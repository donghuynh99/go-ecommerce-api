package client

import (
	"errors"
	"net/http"

	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/database"
	"github.com/donghuynh99/ecommerce_api/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type Controller struct {
	db     *gorm.DB
	config *config.Config
}

func NewController() *Controller {
	var db *gorm.DB = database.Database
	var config config.Config = *config.GetConfig()

	return &Controller{db, &config}
}

func (controller *Controller) CheckExisted(model interface{}, conditions map[string]interface{}, relates *[]string) error {
	db := controller.db
	for key, value := range conditions {
		db = db.Where(key, value)
	}
	if relates != nil {
		for _, value := range *relates {
			db = db.Preload(value)
		}
	}
	checkExisted := db.Take(model)

	if checkExisted.Error != nil {
		return errors.New(utils.Translation("not_found", nil, nil))
	}

	return nil
}

func ChangeLocale(c *gin.Context) {
	locale := c.Query("locale")

	if locale == "" || !slices.Contains([]string{"vi", "en"}, locale) {
		locale = "en"
	}

	cookie := &http.Cookie{
		Name:   "locale",
		Value:  locale,
		MaxAge: 86400, // Set the cookie to expire after 24 hours
		Path:   "/",
	}
	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{
		"message": utils.Translation("change_language_success", nil, nil),
	})

	return
}
