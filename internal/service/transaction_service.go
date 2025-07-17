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
	UserID   uuid.UUID
	OutletID uuid.UUID
	Items    []struct {
		ProductID   uuid.UUID
		ProductName string
		Category    model.ProductCategory
		Qty         int8
		Price       int32
	}
	Note string
}

func (s *TransactionService) CreateTransaction(ctx context.Context, input CreateTransactionInput) (*model.Transaction, error) {
	// 1. Validate stock availability for each item
	for _, item := range input.Items {
		currentStock, err := s.getCurrentStock(ctx, item.ProductID, input.OutletID)
		if err != nil {
			return nil, fmt.Errorf("failed to check stock for product %s: %w", item.ProductName, err)
		}

		if currentStock < int(item.Qty) {
			return nil, fmt.Errorf("insufficient stock for product '%s': available %d, requested %d",
				item.ProductName, currentStock, item.Qty)
		}
	}

	// 2. Calculate total and prepare details
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

	// 3. Generate invoice code and note
	invoiceCode := generateInvoiceCode()
	note := fmt.Sprintf("INV %s includes: %s. Additional notes: %s",
		invoiceCode,
		strings.Join(itemNames, ", "),
		input.Note,
	)

	// 4. Create transaction object
	transaction := model.Transaction{
		UserID:             input.UserID,
		OutletID:           input.OutletID,
		InvoiceCode:        invoiceCode,
		Total:              total,
		Note:               note,
		TransactionDetails: details, // GORM will auto-create these
	}

	// 5. Save transaction in a database transaction
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the transaction
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 6. Create inventory ledger entries asynchronously after successful transaction
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.createInventoryLedgerEntries(transaction.ID, input.OutletID, input.Items)
	}()

	return &transaction, nil
}

// Generate invoice code from random string
func generateInvoiceCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("INV-%d-%04d", time.Now().Year(), rand.Intn(10000))
}

// getCurrentStock calculates the current on-hand stock for a specific item and outlet
// by summing all quantity_change entries in the inventory_ledgers table.
func (s *TransactionService) getCurrentStock(ctx context.Context, itemID, outletID uuid.UUID) (int, error) {
	var totalStock int64

	err := s.db.WithContext(ctx).
		Model(&model.InventoryLedger{}).
		Where("item_id = ? AND outlet_id = ?", itemID, outletID).
		Select("COALESCE(SUM(quantity_change), 0)").
		Row().
		Scan(&totalStock)

	if err != nil {
		return 0, err
	}

	return int(totalStock), nil
}

// createInventoryLedgerEntries creates negative inventory ledger entries for sold items
// This runs asynchronously after a transaction is successfully created
func (s *TransactionService) createInventoryLedgerEntries(transactionID, outletID uuid.UUID, items []struct {
	ProductID   uuid.UUID
	ProductName string
	Category    model.ProductCategory
	Qty         int8
	Price       int32
}) {
	bgCtx := context.Background()
	fmt.Printf("Starting inventory ledger creation for transaction %s\n", transactionID)

	// Create inventory ledger entries in a database transaction
	err := s.db.WithContext(bgCtx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			ledger := model.InventoryLedger{
				ItemId:         item.ProductID,
				OutletId:       outletID,
				TransactionId:  &transactionID,
				QuantityChange: -int(item.Qty), // Negative for stock-out
			}
			if err := tx.Create(&ledger).Error; err != nil {
				return fmt.Errorf("failed to create inventory ledger for product %s: %w", item.ProductName, err)
			}
		}
		return nil
	})

	if err != nil {
		// Log the error but don't fail the transaction since it's already completed
		fmt.Printf("Failed to create inventory ledger entries for transaction %s: %v\n", transactionID, err)
	} else {
		fmt.Printf("Successfully created inventory ledger entries for transaction %s\n", transactionID)
	}
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

	const isPaidQuery = "is_paid = ?"

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
			tx.Model(&model.Transaction{}).Where(isPaidQuery, true).Select("COALESCE(SUM(total), 0)").Row().Scan(&totalRevenue)
			tx.Model(&model.Transaction{}).Where(isPaidQuery, true).Count(&totalPaidTransactions)
			tx.Model(&model.Transaction{}).Where(isPaidQuery, true).Select("COUNT(DISTINCT user_id)").Row().Scan(&totalUniqueCustomers)

			var totalProductsSold uint64
			tx.Model(&model.TransactionDetail{}).Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
				Where("transactions.is_paid = ?", true).Select("COALESCE(SUM(qty), 0)").Row().Scan(&totalProductsSold)

			type CategoryResult struct {
				Category model.ProductCategory
				Count    int64
			}

			var categoryResults []CategoryResult
			tx.Model(&model.TransactionDetail{}).Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
				Where("transactions.is_paid = ?", true).Select("category, SUM(qty) as count").Group("category").Find(&categoryResults)

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
