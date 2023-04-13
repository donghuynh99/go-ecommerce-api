package models

import "time"

type Category struct {
	ID        string `gorm:"size:100;not null;auto_increment,uniqueIndex;primary_key"`
	ParentID  string `gorm:"size:100"`
	Name      string `gorm:"size:100;not null"`
	Slug      string `gorm:"size:100;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Products []Product `gorm:"many2many:product_categories;"`
}
