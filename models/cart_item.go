package models

import (
	"time"
)

type CartItem struct {
	CartID    uint `gorm:"primaryKey"`
	ProductID uint `gorm:"primaryKey"`
	Qty       int

	CreatedAt time.Time
}
