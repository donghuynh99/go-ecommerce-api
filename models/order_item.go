package models

import (
	"time"
)

type OrderItem struct {
	ID        string `gorm:"size:36;not null;auto_increment,uniqueIndex;primary_key"`
	Order     Order
	OrderID   string `gorm:"size:36;index"`
	Product   Product
	ProductID string `gorm:"size:36;index"`
	Qty       int
	CreatedAt time.Time
	UpdatedAt time.Time
}
