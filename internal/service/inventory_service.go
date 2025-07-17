package service

import (
	"context"
	"venturo-core/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InventoryService handles inventory management logic.
type InventoryService struct {
	db *gorm.DB
}

// NewInventoryService creates a new inventory service.
func NewInventoryService(db *gorm.DB) *InventoryService {
	return &InventoryService{db: db}
}

// StockInInput represents the data needed for stock-in operation.
type StockInInput struct {
	ItemID   uuid.UUID `json:"item_id" validate:"required"`
	OutletID uuid.UUID `json:"outlet_id" validate:"required"`
	Quantity int       `json:"quantity" validate:"required,min=1"`
}

// StockIn creates a new record in inventory_ledgers with positive quantity_change.
func (s *InventoryService) StockIn(ctx context.Context, input StockInInput) (*model.InventoryLedger, error) {
	// Create inventory ledger entry
	ledger := model.InventoryLedger{
		ItemId:         input.ItemID,
		OutletId:       input.OutletID,
		TransactionId:  nil,            // No transaction for stock-in operations
		QuantityChange: input.Quantity, // Positive for stock-in
	}

	// Save to database
	if err := s.db.WithContext(ctx).Create(&ledger).Error; err != nil {
		return nil, err
	}

	return &ledger, nil
}
