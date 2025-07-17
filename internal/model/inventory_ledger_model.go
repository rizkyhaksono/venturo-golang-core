package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryLedger struct {
	ID             uuid.UUID  `gorm:"type:char(36);primary_key" json:"id"`
	ItemId         uuid.UUID  `gorm:"type:char(36);not null" json:"item_id"`
	OutletId       uuid.UUID  `gorm:"type:char(36);not null" json:"outlet_id"`
	TransactionId  *uuid.UUID `gorm:"type:char(36)" json:"transaction_id"` // Nullable for stock-in operations
	QuantityChange int        `gorm:"not null" json:"quantity_change"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Item        Product      `gorm:"foreignKey:ItemId;references:ID" json:"item,omitempty"`
	Outlet      Outlet       `gorm:"foreignKey:OutletId;references:ID" json:"outlet,omitempty"`
	Transaction *Transaction `gorm:"foreignKey:TransactionId;references:ID" json:"transaction,omitempty"`
}

// BeforeCreate is a GORM hook that runs before creating a new inventory ledger.
func (il *InventoryLedger) BeforeCreate(tx *gorm.DB) (err error) {
	il.ID = uuid.New()
	return
}
