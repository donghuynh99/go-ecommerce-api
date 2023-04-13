package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	ID           string     `gorm:"size:100;not null;auto_increment,uniqueIndex;primary_key"`
	Categories   []Category `gorm:"many2many:product_categories"`
	ProductImage []ProductImage
	Name         string          `gorm:"size:100;not null"`
	Slug         string          `gorm:"size:100;not null"`
	Price        decimal.Decimal `gorm:"type:decimal(16,2)"`
	Description  string          `gorm:"type:text"`
	Status       int             `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
