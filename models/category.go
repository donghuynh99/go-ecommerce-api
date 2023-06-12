package models

import (
	"time"

	"github.com/donghuynh99/ecommerce_api/database"
)

type Category struct {
	ID        uint   `gorm:"primaryKey"`
	ParentID  *uint  `gorm:"size:100"`
	Name      string `gorm:"size:100;not null;unique"`
	Slug      string `gorm:"size:100;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Products []Product  `gorm:"many2many:product_categories;" json:"products"`
	Children []Category `gorm:"foreignkey:ParentID"`
}

func (c *Category) GetChildren() ([]Category, error) {
	var allDescendants []Category

	GetAllDescendants("parent_id", c.ID, &allDescendants)

	return allDescendants, nil
}

func (c *Category) GetParents() []Category {
	var allParents []Category

	if c.ParentID != nil {
		GetAllDescendants("id", c.ParentID, &allParents)
	}

	allParents = append(allParents, *c)

	return allParents
}

func GetAllDescendants(key string, value interface{}, categories *[]Category) {
	db := database.Database
	db.Where(key, value).Find(categories)

	for _, child := range *categories {
		var childCategories []Category

		if key == "id" {
			GetAllDescendants(key, child.ParentID, &childCategories)
		} else {
			GetAllDescendants(key, child.ID, &childCategories)
		}

		*categories = append(*categories, childCategories...)
	}
}
