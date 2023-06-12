package seeders

import (
	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/database"
	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/donghuynh99/ecommerce_api/utils"
)

func GenerateAdmin(email string, password string) error {
	db := database.Database

	hashPassword, err := utils.HashPassword(password)

	if err != nil {
		return err
	}

	admin := models.User{
		Addresses: []models.Address{},
		FirstName: "Admin",
		LastName:  "Site",
		Role:      config.GetConfig().RoleConfig.Admin,
		Email:     email,
		Password:  hashPassword,
	}

	result := db.Create(&admin)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
