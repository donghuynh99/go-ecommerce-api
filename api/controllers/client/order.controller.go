package client

import (
	"net/http"
	"time"

	"github.com/donghuynh99/ecommerce_api/api/request"
	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/donghuynh99/ecommerce_api/utils"
	"github.com/gin-gonic/gin"
)

func (controller *Controller) ListMyOrder(c *gin.Context) {
	user, ok := utils.GetUser(c)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	var orders []models.Order
	err := controller.db.Preload("Address").Preload("OrderItems.Product.Images").Where("user_id", user.ID).Find(&orders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("server_error", nil, nil),
		})

		return
	}

	var result []config.OrderJsonStruct

	for _, order := range orders {
		result = append(result, order.FormatOrder())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})

	return
}

func (controller *Controller) ShowOrder(c *gin.Context) {
	var request request.OrderCancelRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var order models.Order
	user, ok := utils.GetUser(c)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	err := controller.CheckExisted(&order, map[string]interface{}{
		"id":      c.Param("id"),
		"user_id": user.ID,
	}, &[]string{"Address", "OrderItems.Product.Images"})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": order.FormatOrder(),
	})

	return
}

func (controller *Controller) HandleOrder(c *gin.Context) {
	var request request.OrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, ok := utils.GetUser(c)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	var primaryAddress models.Address
	for _, address := range user.Addresses {
		if address.IsPrimary {
			primaryAddress = address
			break
		}
	}

	if primaryAddress.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("warning_address_order", nil, nil),
		})

		return
	}

	var cartItems []models.CartItem

	err := controller.db.Find(&cartItems, "cart_id", user.Cart.ID).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("server_error", nil, nil),
		})

		return
	}

	var orderItems []models.OrderItem

	for _, cartItem := range cartItems {
		orderItems = append(orderItems, models.OrderItem{
			ProductID: cartItem.ProductID,
			Qty:       cartItem.Qty,
		})
	}

	if len(orderItems) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("add_more_items", nil, nil),
		})

		return
	}

	order := models.Order{
		UserID:     user.ID,
		AddressID:  primaryAddress.ID,
		OrderItems: orderItems,
		Note:       request.Note,
	}

	err = controller.db.Where("cart_id", user.Cart.ID).Delete(&models.CartItem{}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("create_fail", nil, nil),
		})

		return
	}

	err = controller.db.Create(&order).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("create_fail", nil, nil),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": utils.Translation("create_success", nil, nil),
		"data":    order,
	})

	return
}

func (controller *Controller) CancelOrder(c *gin.Context) {
	var request request.OrderCancelRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var order models.Order
	user, ok := utils.GetUser(c)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	err := controller.CheckExisted(&order, map[string]interface{}{
		"id":      c.Param("id"),
		"user_id": user.ID,
	}, nil)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	if order.GetStatus() != config.GetConfig().StatusOrderConfig.Pending {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("condition_cancel", nil, nil),
		})

		return
	}

	currentTime := time.Now()
	order.CancelledBy = &user.ID
	order.CancelledAt = &currentTime
	order.CancellationNote = &request.Note

	err = controller.db.Save(&order).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("cancel_order_fail", nil, nil),
		})

		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": utils.Translation("cancel_order_success", nil, nil),
	})

	return
}
