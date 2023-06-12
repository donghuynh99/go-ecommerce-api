package client

import (
	"net/http"

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

	c.JSON(http.StatusOK, gin.H{
		"data": products,
	})
}
