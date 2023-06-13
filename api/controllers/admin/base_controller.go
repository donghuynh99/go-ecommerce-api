package admin

import (
	"errors"
	"strconv"

	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/database"
	"github.com/donghuynh99/ecommerce_api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController struct {
	db     *gorm.DB
	config *config.Config
}

func NewController() *AdminController {
	var db *gorm.DB = database.Database
	var config config.Config = *config.GetConfig()

	return &AdminController{db, &config}
}

func (controller *AdminController) CheckExisted(model interface{}, conditions map[string]interface{}, relates *[]string) error {
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

func (controller *AdminController) GetDataPagination(c *gin.Context) (int, int, int, error) {
	pageStr := c.DefaultQuery("page", config.GetConfig().PaginationConfig.Page)
	page, errPage := strconv.Atoi(pageStr)
	limitStr := c.DefaultQuery("limit", config.GetConfig().PaginationConfig.Limit)
	limit, errLimit := strconv.Atoi(limitStr)

	if errLimit != nil || errPage != nil {
		return 0, 0, 0, errors.New(utils.Translation("invalid_pagination", nil, nil))
	}

	offset := (page - 1) * limit

	return page, limit, offset, nil
}
