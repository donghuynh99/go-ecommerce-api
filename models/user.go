package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint      `gorm:"primaryKey"`
	Addresses     []Address `json:"addresses"`
	Cart          Cart
	FirstName     string     `gorm:"size:100;not null"`
	LastName      string     `gorm:"size:100;not null"`
	Role          string     `gorm:"size:100, not null"`
	Avatar        Image      `gorm:"polymorphic:Imageable;"`
	Email         string     `gorm:"size:100;not null;uniqueIndex"`
	Password      string     `gorm:"size:255;not null"`
	RememberToken *string    `gorm:"size:255"`
	Token         *string    `gorm:"size:255"`
	ExpiredAt     *time.Time `gorm:"default:null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
