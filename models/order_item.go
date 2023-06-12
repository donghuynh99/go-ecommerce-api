package models

import (
	"time"
)

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	Order     Order
	OrderID   uint `gorm:"size:255;index"`
	Product   Product
	ProductID uint `gorm:"size:255;index"`
	Qty       int
	CreatedAt time.Time
	UpdatedAt time.Time
}
