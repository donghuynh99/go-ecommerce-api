package client

import (
	"net/http"

	"github.com/donghuynh99/ecommerce_api/api/request"
	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/donghuynh99/ecommerce_api/utils"
	"github.com/gin-gonic/gin"
)

func (controller *Controller) ListAddress(c *gin.Context) {
	var addresses []config.AddressListStruct
	user, ok := utils.GetUser(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": utils.Translation("unauthenticate", nil, nil),
		})

		return
	}

	err := controller.db.Model(&models.Address{}).Where("user_id", user.ID).Scan(&addresses).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("server_error", nil, nil),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": addresses,
	})

	return
}

func (controller *Controller) AddAddress(c *gin.Context) {
	var request request.AddressRequest

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

	if len(user.Addresses) > 0 {
		if request.IsPrimary == true {
			controller.db.Model(&models.Address{}).Where("user_id", user.ID).Update("is_primary", false)
		}
	} else {
		request.IsPrimary = true
	}

	address := models.Address{
		UserId:    user.ID,
		Name:      request.Name,
		IsPrimary: request.IsPrimary,
		PostCode:  request.PostCode,
	}

	errStore := controller.db.Create(&address).Error

	if errStore != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("create_fail", nil, nil),
			"error":   errStore.Error,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": utils.Translation("create_success", nil, nil),
		"data":    address,
	})

	return
}

func (controller *Controller) RemoveAddress(c *gin.Context) {
	var address models.Address
	user, ok := utils.GetUser(c)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	err := controller.CheckExisted(&address, map[string]interface{}{
		"id":      c.Param("id"),
		"user_id": user.ID,
	}, nil)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	if len(user.Addresses) > 1 && address.IsPrimary {
		var nextAddress models.Address
		err := controller.CheckExisted(&nextAddress, map[string]interface{}{
			"user_id": user.ID,
			"id != ?": address.ID,
		}, nil)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": utils.Translation("server_error", nil, nil),
			})

			return
		}

		err = controller.db.Model(&nextAddress).Update("is_primary", true).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": utils.Translation("server_error", nil, nil),
			})

			return
		}
	}

	err = controller.db.Delete(&address).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("delete_fail", nil, nil),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": utils.Translation("delete_success", nil, nil),
	})

	return
}

func (controller *Controller) UpdateAddress(c *gin.Context) {
	var request request.AddressRequest

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

	var address models.Address
	user, ok := utils.GetUser(c)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	err := controller.CheckExisted(&address, map[string]interface{}{
		"id":      c.Param("id"),
		"user_id": user.ID,
	}, nil)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	if len(user.Addresses) > 1 && request.IsPrimary == true {
		controller.db.Model(&models.Address{}).Where("user_id", user.ID).Update("is_primary", false)
	}

	if !request.IsPrimary && address.IsPrimary {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("must_have_primary_address", nil, nil),
		})

		return
	}

	address.Name = request.Name
	address.PostCode = request.PostCode
	address.IsPrimary = request.IsPrimary

	err = controller.db.Save(&address).Error

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
