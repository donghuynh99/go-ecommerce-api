package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint       `gorm:"primaryKey"`
	UserID    uint       `gorm:"size:255;index"`
	CartItems []CartItem `json:"cart_items"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeleteAt  gorm.DeletedAt
}
