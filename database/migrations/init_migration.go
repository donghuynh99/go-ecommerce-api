package migrations

import (
	"log"

	"github.com/donghuynh99/ecommerce_api/database"
	"github.com/donghuynh99/ecommerce_api/models"
)

func RunMigration() {
	for _, model := range models.RegisteredModel() {
		err := database.Database.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}
}
