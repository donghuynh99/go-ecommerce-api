package client

import (
	"log"
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

	user := models.User{
		Addresses: []models.Address{},
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Cart: models.Cart{
			CartItems: cartItems,
		},
		Role:     config.GetConfig().RoleConfig.User,
		Password: password,
	}

	result := controller.db.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("register_fail", nil, nil),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": utils.Translation("register_success", nil, nil),
		"user":    user,
	})

	return
}

func (controller *Controller) ShowProfile(c *gin.Context) {
	user, ok := utils.GetUser(c)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	var avatarlURL config.ImageStruct

	if user.Avatar.ID == 0 {
		avatarlURL = config.ImageStruct{
			Path: config.GetConfig().AppConfig.DefaultAvatarURL,
			Alt:  "default_avatar",
		}
	} else {
		avatarlURL = config.ImageStruct{
			Path: user.Avatar.Path,
			Alt:  user.Avatar.Alt,
		}
	}

	var addresses []config.AddressListStruct

	for _, address := range user.Addresses {
		addresses = append(addresses, config.AddressListStruct{
			ID:        address.ID,
			Name:      address.Name,
			PostCode:  address.PostCode,
			IsPrimary: address.IsPrimary,
		})
	}

	profile := config.ProfileStruct{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    avatarlURL,
		Email:     user.Email,
		Addresses: addresses,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": profile,
	})

	return
}

func (controller *Controller) UpdateProfile(c *gin.Context) {
	var request request.UpdateProfileRequest

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

	user, ok := utils.GetUser(c)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	if request.OldPassword != "" && !utils.ComparePassword(request.OldPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ApiError{
			Param:   "old_password",
			Message: utils.Translation("password_not_correct", nil, nil),
		}})
		return
	}

	user.FirstName = request.FirstName
	user.LastName = request.LastName

	if request.OldPassword != "" && request.Password != "" && request.ConfirmPassword != "" {
		newPassword, err := utils.HashPassword(request.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": utils.Translation("server_error", nil, nil),
			})

			return
		}

		user.Password = newPassword
	}

	err := controller.db.Save(&user).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("update_fail", nil, nil),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": utils.Translation("update_success", nil, nil),
	})

	return
}

func (controller *Controller) UpdateAvatar(c *gin.Context) {
	form, errForm := c.MultipartForm()
	if errForm != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("bad_request", nil, nil),
			"error":   errForm.Error(),
		})

		return
	}

	images := form.File["avatar"]

	if len(images) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("just_one_avatar", nil, nil),
		})

		return
	}

	user, ok := utils.GetUser(c)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	avatar := form.File["avatar"][0]

	destination := config.GetConfig().GeneralConfig.DestinationStoreAvatarUser
	fileName, err := utils.UploadImage(destination, avatar, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("fail_upload", nil, nil),
		})

		return
	}

	if user.Avatar.ID != 0 {
		err := controller.db.Delete(&user.Avatar).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": utils.Translation("fail_upload", nil, nil),
			})

			return
		}
	}

	user.Avatar = models.Image{
		Path: destination + fileName,
		Name: fileName,
		Alt:  "avatar " + user.FirstName + " " + user.LastName,
	}

	err = controller.db.Save(&user).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("fail_upload", nil, nil),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": utils.Translation("update_success", nil, nil),
	})

	return
}
