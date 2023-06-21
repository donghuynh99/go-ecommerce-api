package models

import (
	"time"
)

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint `gorm:"size:255;index"`
	ProductID uint `gorm:"size:255;index"`
	Product   Product
	Qty       int
	CreatedAt time.Time
	UpdatedAt time.Time
}
