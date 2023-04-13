package models

import "time"

type Address struct {
	ID        string `gorm:"size:36;not null;auto_increment,uniqueIndex;primary_key"`
	User      User
	UserId    string `gorm:"size:36;index"`
	Name      string `gorm:"size:100;not null"`
	IsPrimary bool
	PostCode  string `gorm:"size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
