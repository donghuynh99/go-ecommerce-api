package models

import (
	"log"
	"os"

	"gorm.io/gorm"
)

type Image struct {
	ID            uint   `gorm:"primaryKey"`
	ImageableID   uint   `gorm:"not null"`
	ImageableType string `gorm:"not null"`
	Path          string `gorm:"type:text;not null"`
	Name          string `gorm:"size:255;not null"`
	Alt           string `gorm:"type:text"`
}

func (i *Image) BeforeDelete(tx *gorm.DB) (err error) {
	error := os.Remove(i.Path)
	if error != nil {
		log.Println(error, i.Path)
	}

	return nil
}
