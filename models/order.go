package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID               string `gorm:"size:36;not null;auto_increment;uniqueIndex;primary_key"`
	UserID           string `gorm:"size:36;index"`
	User             User
	OrderItems       []OrderItem
	Code             string `gorm:"size:50;index"`
	Status           int
	OrderDate        time.Time
	Note             string `gorm:"type:text"`
	ApprovedBy       string `gorm:"size:36"`
	ApprovedAt       time.Time
	CompletedAt      time.Time
	CancelledBy      string `gorm:"size:36"`
	CancelledAt      time.Time
	CancellationNote string `gorm:"size:255"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}
