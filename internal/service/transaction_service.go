package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
	"venturo-core/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionService struct {
	db *gorm.DB
	wg *sync.WaitGroup
}

func NewTransactionService(db *gorm.DB, wg *sync.WaitGroup) *TransactionService {
	return &TransactionService{db: db, wg: wg}
}

type CreateTransactionInput struct {
	UserID uuid.UUID
	Items  []struct {
		ProductID   uuid.UUID
		ProductName string
		Category    model.ProductCategory
		Qty         int8
		Price       int32
	}
	Note string
}

func (s *TransactionService) CreateTransaction(ctx context.Context, input CreateTransactionInput) (*model.Transaction, error) {
	// 1. Calculate total and prepare details
	var total int64
	var details []model.TransactionDetail
	var itemNames []string

	for _, item := range input.Items {
		total += int64(item.Qty) * int64(item.Price)
		details = append(details, model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Category:    item.Category,
			Qty:         item.Qty,
			Price:       item.Price,
		})
		itemNames = append(itemNames, item.ProductName)
	}

	// 2. Generate invoice code and note
	invoiceCode := generateInvoiceCode()
	note := fmt.Sprintf("INV %s includes: %s. Additional notes: %s",
		invoiceCode,
		strings.Join(itemNames, ", "),
		input.Note,
	)

	// 3. Create transaction object
	transaction := model.Transaction{
		UserID:             input.UserID,
		InvoiceCode:        invoiceCode,
		Total:              total,
		Note:               note,
		TransactionDetails: details, // GORM will auto-create these
	}

	// 4. Save transaction in the database
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		// The details should be created automatically due to the relationship,
		// but explicit creation is safer if auto-creation is disabled.
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// Generate invoice code from random string
func generateInvoiceCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("INV-%d-%04d", time.Now().Year(), rand.Intn(10000))
}

func (s *TransactionService) MarkAsPaid(ctx context.Context, transactionId uuid.UUID) error {
	var transaction model.Transaction

	if err := s.db.WithContext(ctx).First(&transaction, "id = ?", transactionId).Error; err != nil {
		return errors.New("transaction not found")
	}

	isPaid := true
	transaction.IsPaid = &isPaid
	if err := transaction.Save(s.db); err != nil {
		return err
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		bgCtx := context.Background()

		err := s.db.WithContext(bgCtx).Transaction(func(tx *gorm.DB) error {
			var report model.TransactionReport
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&report, "id = ?", 1).Error; err != nil {
				return nil
			}

			var totalRevenue, totalUniqueCustomers uint64
			var totalPaidTransactions int64
			tx.Model(&model.Transaction{}).Where("is_paid = ?", true).Select("COALESCE(SUM(total), 0)").Row().Scan(&totalRevenue)
			tx.Model(&model.Transaction{}).Where("is_paid = ?", true).Count(&totalPaidTransactions)
			tx.Model(&model.Transaction{}).Where("is_paid = ?", true).Select("COUNT(DISTINCT user_id)").Row().Scan(&totalUniqueCustomers)

			var totalProductsSold uint64
			tx.Model(&model.TransactionDetail{}).Joins("JOIN transactions ON transaction.id = transaction_details.transaction_id").
				Where("transactions.is_paid = ?", true).Select("COALESCE(SUM(qty), 0)").Row().Scan(&totalProductsSold)

			type CategoryResult struct {
				Category model.ProductCategory
				Count    int64
			}

			var categoryResults []CategoryResult
			tx.Model(&model.TransactionDetail{}).Joins("JOIN transaction ON transaction.id = transaction_details.transaction_id").
				Where("transaction.is_paid = ?", true).Select("category, SUM(qty) as count").Group("category").Find(&categoryResults)

			report.TotalRevenue = totalRevenue
			report.TotalPaidTransactions = totalPaidTransactions
			report.TotalProductsSold = totalProductsSold
			report.TotalUniqueCustomers = totalUniqueCustomers
			report.CategorySummary = make(model.CategorySummary)
			for _, res := range categoryResults {
				var categoryName string
				switch res.Category {
				case model.Goods:
					categoryName = "Goods"
				case model.Service:
					categoryName = "Service"
				case model.Subscription:
					categoryName = "Subscription"
				default:
					categoryName = "Other"
				}
				report.CategorySummary[categoryName] = res.Count
			}

			return report.Save(tx)
		})

		if err != nil {
			fmt.Printf("Failed to update transaction report: %v\n", err)
		}
	}()

	return nil
}
