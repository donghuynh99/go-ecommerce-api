package client

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/donghuynh99/ecommerce_api/api/middleware"
	"github.com/donghuynh99/ecommerce_api/api/request"
	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/gin-gonic/gin"
)

func (controller *Controller) UpdateCart(c *gin.Context) {
	var request request.CartRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := middleware.CheckAuthenticate(c, config.GetConfig().RoleConfig.User)

	if user == nil {
		err := controller.UpdateCartForSession(request.Data, c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})

			return
		}
	} else {
		err, statusCode := controller.UpdateCartForUser(request.Data, user)

		if err != nil {
			c.JSON(statusCode, gin.H{
				"message": err.Error(),
			})

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Update successful",
	})

	return
}

func (controller *Controller) GetCart(c *gin.Context) {
	user := middleware.CheckAuthenticate(c, config.GetConfig().RoleConfig.User)
	var cartData []config.CartInformation

	if user == nil {
		cartCookie, err := c.Request.Cookie("cart")

		if err != nil {
			data, _ := json.Marshal(struct{}{})
			cartCookie = &http.Cookie{
				Name:   "cart",
				Value:  string(data),
				MaxAge: 604800, // Set the cookie to expire after 1 week
				Path:   "/",
			}

			http.SetCookie(c.Writer, cartCookie)
		}

		var data map[string]int

		valueCart := cartCookie.Value
		valueCart = strings.ReplaceAll(valueCart, ":", `":`)
		valueCart = strings.ReplaceAll(valueCart, ",", `,"`)
		valueCart = strings.ReplaceAll(valueCart, "{", `{"`)

		err = json.Unmarshal([]byte(valueCart), &data)

		if err != nil {
			log.Println(err)
		}

		for productID, quantity := range data {
			var product models.Product
			err := controller.CheckExisted(&product, map[string]interface{}{
				"id": productID,
			}, &[]string{"Images"})
			if err != nil {
				log.Println(errors.New("ID of product not found!"))

				continue
			}

			var thumbnail string
			if len(product.Images) > 0 {
				thumbnail = product.Images[0].Path
			} else {
				thumbnail = config.GetConfig().AppConfig.DefaultImageURL
			}

			cartData = append(cartData, config.CartInformation{
				ProductName: product.Name,
				Thumbnail:   thumbnail,
				Quantity:    quantity,
			})
		}

	} else {
		var cartItems []models.CartItem
		controller.db.Find(&cartItems, "cart_id", user.Cart.ID)

		for _, cartItem := range cartItems {
			var product models.Product
			err := controller.CheckExisted(&product, map[string]interface{}{
				"id": cartItem.ProductID,
			}, &[]string{"Images"})
			if err != nil {
				log.Println(errors.New("ID of product not found!"))

				continue
			}

			var thumbnail string
			if len(product.Images) > 0 {
				thumbnail = product.Images[0].Path
			} else {
				thumbnail = config.GetConfig().AppConfig.DefaultImageURL
			}

			cartData = append(cartData, config.CartInformation{
				ProductName: product.Name,
				Thumbnail:   thumbnail,
				Quantity:    cartItem.Qty,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cartData,
	})

	return
}

func (controller *Controller) UpdateCartForSession(data map[string]int, c *gin.Context) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.New("Failed to marshal JSON")
	}

	cookie := &http.Cookie{
		Name:   "cart",
		Value:  string(jsonData),
		MaxAge: 604800, // Set the cookie to expire after 1 week
		Path:   "/",
	}
	http.SetCookie(c.Writer, cookie)

	return nil
}

func (controller *Controller) UpdateCartForUser(data map[string]int, user *models.User) (error, int) {
	cart, err := controller.GetCartOfUser(user)

	if err != nil {
		return err, http.StatusInternalServerError
	}

	// load product item
	var cartItems []models.CartItem
	for productID, quantity := range data {
		var product models.Product
		err := controller.CheckExisted(&product, map[string]interface{}{
			"id": productID,
		}, nil)
		if err != nil {
			return errors.New("ID of product not found!"), http.StatusBadRequest
		}

		cartItems = append(cartItems, models.CartItem{
			ProductID: product.ID,
			CartID:    cart.ID,
			Qty:       quantity,
		})
	}

	controller.db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})

	controller.db.Model(&cart).Association("CartItems").Clear()
	errChange := controller.db.Model(&cart).Association("CartItems").Append(&cartItems)

	if errChange != nil {
		return errChange, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (controller *Controller) GetCartOfUser(user *models.User) (models.Cart, error) {
	var cart models.Cart
	err := controller.CheckExisted(&cart, map[string]interface{}{
		"user_id": user.ID,
	}, &[]string{"CartItems"})

	if err == nil {
		return cart, nil
	}

	newCart := models.Cart{
		UserID:    user.ID,
		CartItems: []models.CartItem{},
	}

	err = controller.db.Create(&newCart).Error

	if err != nil {
		return models.Cart{}, errors.New("Create fail")
	}

	return newCart, nil
}
