package models

import (
	"time"
)

type Address struct {
	ID        uint `gorm:"primaryKey"`
	User      User
	UserId    uint   `gorm:"size:255;index"`
	Name      string `gorm:"size:100;not null"`
	IsPrimary bool
	PostCode  string `gorm:"size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
