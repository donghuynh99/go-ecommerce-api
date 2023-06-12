package client

import (
	"net/http"

	"github.com/donghuynh99/ecommerce_api/api/request"
	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/donghuynh99/ecommerce_api/utils"
	"github.com/gin-gonic/gin"
)

func (controller *Controller) Register(c *gin.Context) {
	var request request.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("bad_request", nil, nil),
			"error":   err.Error(),
		})

		return
	}

	errors := utils.ValidateStruct(request)

	if errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors})
		return
	}

	password, _ := utils.HashPassword(request.Password)

	// cartData := utils.GetCartDataFromCookie(c)

	// var cartItems []models.CartItem

	// for productID, quantity := range cartData {
	// 	var product models.Product
	// 	err := controller.CheckExisted(&product, map[string]interface{}{
	// 		"id": productID,
	// 	}, nil)
	// 	if err != nil {
	// 		log.Println(err)

	// 		continue
	// 	}

	// 	cartItems = append(cartItems, models.CartItem{
	// 		ProductID: product.ID,
	// 		Qty:       quantity,
	// 	})
	// }

	user := models.User{
		Addresses: []models.Address{},
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		// Cart: models.Cart{
		// 	CartItems: cartItems,
		// },
		Role:     config.GetConfig().RoleConfig.User,
		Password: password,
	}

	result := controller.db.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Register user fail!",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Register successful!",
		"user":    user,
	})

	return
}
