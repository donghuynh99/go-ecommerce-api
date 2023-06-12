package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint            `gorm:"primaryKey"`
	Categories  []Category      `gorm:"many2many:product_categories;" json:"categories"`
	Images      []Image         `gorm:"polymorphic:Imageable;" json:"images"`
	Name        string          `gorm:"size:100;not null;"`
	Slug        string          `gorm:"size:100;not null;"`
	Price       decimal.Decimal `gorm:"type:decimal(16,2);not null;"`
	Description string          `gorm:"type:text;not null;"`
	Status      int             `gorm:"default:0;not null;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
