package models

import "time"

type OrderCustomer struct {
	ID        uint `gorm:"primaryKey"`
	User      User
	UserID    uint `gorm:"size:255;index"`
	Order     Order
	OrderID   uint `gorm:"size:255;index"`
	Address   Address
	AddressID uint   `gorm:"size:255;index"`
	Phone     string `gorm:"size:50;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
