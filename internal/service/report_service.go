package service

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportService handles report generation logic.
type ReportService struct {
	db *gorm.DB
}

// NewReportService creates a new report service.
func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

// InventoryReportItem represents the inventory status for a specific item at an outlet.
type InventoryReportItem struct {
	ItemID       uuid.UUID            `json:"item_id"`
	ItemName     string               `json:"item_name"`
	OutletID     uuid.UUID            `json:"outlet_id"`
	OutletName   string               `json:"outlet_name"`
	OnHandQty    int                  `json:"on_hand_quantity"`
	Transactions []TransactionHistory `json:"transaction_history"`
}

// TransactionHistory represents a single transaction that affected the inventory.
type TransactionHistory struct {
	TransactionID   *uuid.UUID `json:"transaction_id"`
	InvoiceCode     *string    `json:"invoice_code"`
	QuantityChange  int        `json:"quantity_change"`
	TransactionType string     `json:"transaction_type"` // "stock-in" or "stock-out"
	CreatedAt       string     `json:"created_at"`
}

// InventoryReportInput represents the filter criteria for the inventory report.
type InventoryReportInput struct {
	ItemID   *uuid.UUID `json:"item_id"`
	OutletID *uuid.UUID `json:"outlet_id"`
}

// GenerateInventoryReport generates a comprehensive inventory report.
func (s *ReportService) GenerateInventoryReport(ctx context.Context, input InventoryReportInput) ([]InventoryReportItem, error) {
	// Build base query for inventory aggregation
	query := s.db.WithContext(ctx).
		Table("inventory_ledgers").
		Select(`
			inventory_ledgers.item_id,
			products.name as item_name,
			inventory_ledgers.outlet_id,
			outlets.name as outlet_name,
			COALESCE(SUM(inventory_ledgers.quantity_change), 0) as on_hand_qty
		`).
		Joins("LEFT JOIN products ON products.id = inventory_ledgers.item_id").
		Joins("LEFT JOIN outlets ON outlets.id = inventory_ledgers.outlet_id").
		Group("inventory_ledgers.item_id, inventory_ledgers.outlet_id, products.name, outlets.name")

	// Apply filters
	if input.ItemID != nil {
		query = query.Where("inventory_ledgers.item_id = ?", *input.ItemID)
	}
	if input.OutletID != nil {
		query = query.Where("inventory_ledgers.outlet_id = ?", *input.OutletID)
	}

	// Execute the aggregation query
	type AggregationResult struct {
		ItemID     uuid.UUID `json:"item_id"`
		ItemName   string    `json:"item_name"`
		OutletID   uuid.UUID `json:"outlet_id"`
		OutletName string    `json:"outlet_name"`
		OnHandQty  int       `json:"on_hand_qty"`
	}

	var aggregationResults []AggregationResult
	if err := query.Find(&aggregationResults).Error; err != nil {
		return nil, err
	}

	// Build the report items with transaction history
	var reportItems []InventoryReportItem
	for _, aggResult := range aggregationResults {
		// Get transaction history for this item/outlet combination
		transactionHistory, err := s.getTransactionHistory(ctx, aggResult.ItemID, aggResult.OutletID)
		if err != nil {
			return nil, err
		}

		reportItem := InventoryReportItem{
			ItemID:       aggResult.ItemID,
			ItemName:     aggResult.ItemName,
			OutletID:     aggResult.OutletID,
			OutletName:   aggResult.OutletName,
			OnHandQty:    aggResult.OnHandQty,
			Transactions: transactionHistory,
		}

		reportItems = append(reportItems, reportItem)
	}

	return reportItems, nil
}

// getTransactionHistory retrieves the transaction history for a specific item/outlet combination.
func (s *ReportService) getTransactionHistory(ctx context.Context, itemID, outletID uuid.UUID) ([]TransactionHistory, error) {
	var history []TransactionHistory

	// Query for all inventory ledger entries with transaction details
	rows, err := s.db.WithContext(ctx).
		Table("inventory_ledgers").
		Select(`
			inventory_ledgers.transaction_id,
			transactions.invoice_code,
			inventory_ledgers.quantity_change,
			inventory_ledgers.created_at
		`).
		Joins("LEFT JOIN transactions ON transactions.id = inventory_ledgers.transaction_id").
		Where("inventory_ledgers.item_id = ? AND inventory_ledgers.outlet_id = ?", itemID, outletID).
		Order("inventory_ledgers.created_at DESC").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transactionID *uuid.UUID
		var invoiceCode *string
		var quantityChange int
		var createdAt string

		if err := rows.Scan(&transactionID, &invoiceCode, &quantityChange, &createdAt); err != nil {
			return nil, err
		}

		// Determine transaction type based on quantity change and transaction ID
		transactionType := "stock-in"
		if quantityChange < 0 {
			transactionType = "stock-out"
		}

		historyItem := TransactionHistory{
			TransactionID:   transactionID,
			InvoiceCode:     invoiceCode,
			QuantityChange:  quantityChange,
			TransactionType: transactionType,
			CreatedAt:       createdAt,
		}

		history = append(history, historyItem)
	}

	return history, nil
}
