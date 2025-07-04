package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type CategorySummary map[string]int64

func (cs CategorySummary) Value() (driver.Value, error) {
	return json.Marshal(cs)
}

func (cs *CategorySummary) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan CategorySummary: value is not a byte slice")
	}
	return json.Unmarshal(b, &cs)
}

type TransactionReport struct {
	ID                    uint8 `gorm:"primary_key"`
	TotalRevenue          uint64
	TotalPaidTransactions int64
	TotalProductsSold     uint64
	TotalUniqueCustomers  uint64
	CategorySummary       CategorySummary `gorm:"type:json"`
	UpdatedAt             time.Time
}

func (tr *TransactionReport) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(tr).Error
}
