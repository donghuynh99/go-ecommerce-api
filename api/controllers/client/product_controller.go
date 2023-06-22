package client

import (
	"net/http"

	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/gin-gonic/gin"
)

func (controller *Controller) GetPopularProducts(c *gin.Context) {
	var products []models.Product

	err := controller.db.Preload("Images").Find(&products).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	}

	var results []config.ProductJsonStruct

	for _, product := range products {
		thumbnailURL := config.ImageStruct{
			Path: config.GetConfig().AppConfig.DefaultImageURL,
			Alt:  "default_image",
		}

		if len(product.Images) > 0 {
			thumbnailURL = config.ImageStruct{
				Path: product.Images[0].Path,
				Alt:  product.Images[0].Alt,
			}
		}

		results = append(results, config.ProductJsonStruct{
			ID:           product.ID,
			Name:         product.Name,
			Price:        product.Price,
			Description:  product.Description,
			ThumbnailURL: thumbnailURL,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}
