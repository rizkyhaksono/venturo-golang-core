package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionDetail defines the items within a transaction.
type TransactionDetail struct {
	ID            uuid.UUID       `gorm:"type:char(36);primary_key"`
	TransactionID uuid.UUID       `gorm:"type:char(36);not null"`
	ProductID     uuid.UUID       `gorm:"type:char(36);not null"`
	ProductName   string          `gorm:"size:255;not null"`
	Category      ProductCategory `gorm:"not null"`
	Qty           int8            `gorm:"not null"`
	Price         int32           `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// BeforeCreate is a GORM hook.
func (td *TransactionDetail) BeforeCreate(tx *gorm.DB) (err error) {
	td.ID = uuid.New()
	return
}
