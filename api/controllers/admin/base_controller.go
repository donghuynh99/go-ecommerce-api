package admin

import (
	"errors"
	"strconv"

	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/database"
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
		return errors.New("Not found!")
	}

	return nil
}

func (controller *AdminController) GetDataPagination(c *gin.Context) (int, int, int, error) {
	pageStr := c.DefaultQuery("page", "1")
	page, errPage := strconv.Atoi(pageStr)
	limitStr := c.DefaultQuery("limit", "10")
	limit, errLimit := strconv.Atoi(limitStr)

	if errLimit != nil || errPage != nil {
		return 0, 0, 0, errors.New("Invalid page or limit value")
	}

	offset := (page - 1) * limit

	return page, limit, offset, nil
	// var totalCount int64

	// offset := (page - 1) * limit
	// db := controller.db
	// if relates != nil {
	// 	for _, value := range *relates {
	// 		db = db.Preload(value)
	// 	}
	// }
	// db.Offset(offset).Limit(limit).Find(&items)
	// db.Count(&totalCount)

	// return PaginationStruct{
	// 	Data:       items,
	// 	TotalCount: int(totalCount),
	// 	Page:       page,
	// }, nil
}
