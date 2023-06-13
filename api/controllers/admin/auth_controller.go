package admin

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/donghuynh99/ecommerce_api/api/request"
	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/donghuynh99/ecommerce_api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (controller *AdminController) Login(c *gin.Context) {
	var request request.LoginRequest

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

	user, err := AuthenticateUser(controller.db, request.Email, request.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("credential_wrong", nil, nil),
		})

		return
	}

	token, expiredAt, err := utils.GenerateToken(int(user.ID), user.Role)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("bad_request", nil, nil),
		})

		return
	}

	if user.Role == config.GetConfig().RoleConfig.User && user.Cart.ID == 0 {
		cartData := utils.GetCartDataFromCookie(c)

		var cartItems []models.CartItem

		for productID, quantity := range cartData {
			var product models.Product
			err := controller.CheckExisted(&product, map[string]interface{}{
				"id": productID,
			}, nil)
			if err != nil {
				log.Println(err)

				continue
			}

			cartItems = append(cartItems, models.CartItem{
				ProductID: product.ID,
				Qty:       quantity,
			})
		}
		user.Cart = models.Cart{
			CartItems: cartItems,
		}
	}

	user.Token = &token
	user.ExpiredAt = &expiredAt

	if user.Role == config.GetConfig().RoleConfig.User {
		controller.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user)
	} else {
		controller.db.Save(&user)
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expired_at": expiredAt.Format(time.DateTime),
	})

	return
}

func AuthenticateUser(db *gorm.DB, email string, password string) (*models.User, error) {
	var user models.User
	result := db.Preload(clause.Associations).First(&user, "email = ?", email)

	if result.Error != nil || !utils.ComparePassword(password, user.Password) {
		return nil, errors.New("Fail")
	}

	return &user, nil
}

func (controller *AdminController) Logout(c *gin.Context) {
	user, ok := c.Get("user")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	userModel, ok := user.(models.User)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	currentTime := time.Now()
	userModel.ExpiredAt = &currentTime
	controller.db.Save(&userModel)

	c.JSON(http.StatusOK, gin.H{
		"message": utils.Translation("login_success", nil, nil),
	})

	return
}
