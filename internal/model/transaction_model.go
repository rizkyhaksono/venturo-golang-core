package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	UserID      uuid.UUID `gorm:"type:char(36);not null"`
	InvoiceCode string    `gorm:"size:20;not null;unique"`
	Total       int64     `gorm:"not null"`
	IsPaid      *bool     `gorm:"not null;default:false" json:"is_paid"`
	Note        string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relationships
	User               User                `gorm:"foreignKey:UserID"`
	TransactionDetails []TransactionDetail `gorm:"foreignKey:TransactionID"`
}

// BeforeCreate is a GORM hook.
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	// Set default IsPaid to false if it's nil
	if t.IsPaid == nil {
		b := false
		t.IsPaid = &b
	}
	return
}

// Save creates or updates a record.
func (t *Transaction) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(t).Error
}
