package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	Name        string    `gorm:"size:255;not null"`
	Price       int32
	Stock       int16
	ImageURL    string `gorm:"size:255"`
	ImageStatus string `gorm:"size:20;not null;default:'default'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return nil
}

func (p *Product) Save(db *gorm.DB) (err error) {
	return db.WithContext(context.Background()).Save(p).Error
}
