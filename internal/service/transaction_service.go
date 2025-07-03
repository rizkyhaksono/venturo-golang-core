package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"venturo-core/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{db: db}
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
