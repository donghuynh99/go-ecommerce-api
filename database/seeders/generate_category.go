package seeders

import (
	"github.com/donghuynh99/ecommerce_api/database"
	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/gosimple/slug"
)

func GenerateBaseCategory() error {
	db := database.Database

	tx := db.Begin()

	categories := []models.Category{
		{
			Name: "Phone",
			Slug: slug.Make("Phone"),
		},
		{
			Name: "Laptop",
			Slug: slug.Make("Laptop"),
		},
		{
			Name: "Accessory",
			Slug: slug.Make("Accessory"),
		},
		{
			Name: "Personal computer",
			Slug: slug.Make("Personal computer"),
		},
	}

	err := db.Create(&categories)

	if err.Error != nil {
		tx.Rollback()

		return err.Error
	}

	var phone models.Category
	db.First(&phone, "name = ?", "Phone")

	var laptop models.Category
	db.First(&laptop, "name = ?", "Laptop")

	var accessory models.Category
	db.First(&accessory, "name = ?", "Accessory")

	childCategories := []models.Category{
		{
			ParentID: &phone.ID,
			Name:     "Iphone",
			Slug:     slug.Make("Iphone"),
		},
		{
			ParentID: &phone.ID,
			Name:     "Samsung",
			Slug:     slug.Make("Samsung"),
		},
		{
			ParentID: &phone.ID,
			Name:     "Xiaomi",
			Slug:     slug.Make("Xiaomi"),
		},
		{
			ParentID: &phone.ID,
			Name:     "Oppo",
			Slug:     slug.Make("Oppo"),
		},
		{
			ParentID: &laptop.ID,
			Name:     "Mac",
			Slug:     slug.Make("Mac"),
		},
		{
			ParentID: &laptop.ID,
			Name:     "Dell",
			Slug:     slug.Make("Dell"),
		},
		{
			ParentID: &laptop.ID,
			Name:     "Asus",
			Slug:     slug.Make("Asus"),
		},
		{
			ParentID: &laptop.ID,
			Name:     "Lenovo",
			Slug:     slug.Make("Lenovo"),
		},
		{
			ParentID: &accessory.ID,
			Name:     "Phone accessory",
			Slug:     slug.Make("Phone accessory"),
		},
		{
			ParentID: &accessory.ID,
			Name:     "Laptop accessory",
			Slug:     slug.Make("Laptop accessory"),
		},
		{
			ParentID: &accessory.ID,
			Name:     "Camera",
			Slug:     slug.Make("Camera"),
		},
		{
			ParentID: &accessory.ID,
			Name:     "Network accessory",
			Slug:     slug.Make("Network accessory"),
		},
	}

	err = db.Create(&childCategories)

	if err.Error != nil {
		tx.Rollback()

		return err.Error
	}

	tx.Commit()

	return nil
}
