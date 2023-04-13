package models

import "time"

type OrderCustomer struct {
	ID        string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	User      User
	UserID    string `gorm:"size:36;index"`
	Order     Order
	OrderID   string `gorm:"size:36;index"`
	Address   Address
	Phone     string `gorm:"size:50;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
