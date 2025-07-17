package model

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Outlet struct {
	ID   uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Name string    `gorm:"size:255;not null" json:"name"`
}

func (o *Outlet) BeforeCReate() (err error) {
	o.ID = uuid.New()
	return
}

func (o *Outlet) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(0).Error
}
