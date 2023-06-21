package models

import (
	"time"

	"github.com/donghuynh99/ecommerce_api/config"
	"gorm.io/gorm"
)

type Order struct {
	ID               uint `gorm:"primaryKey"`
	UserID           uint `gorm:"size:255;index"`
	User             User
	AddressID        uint `gorm:"size:255;index"`
	Address          Address
	OrderItems       []OrderItem `json:"order_items"`
	Note             string      `gorm:"type:text"`
	ApprovedBy       *uint       `gorm:"size:255;index"`
	ApprovedAt       *time.Time
	CompletedAt      *time.Time
	CancelledBy      *uint `gorm:"size:255;index"`
	CancelledAt      *time.Time
	CancellationNote *string `gorm:"size:255"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

func (order *Order) GetStatus() string {
	if order.ApprovedAt == nil {
		return config.GetConfig().StatusOrderConfig.Pending
	} else if order.CompletedAt != nil {
		return config.GetConfig().StatusOrderConfig.Completed
	} else if order.CancelledAt != nil {
		return config.GetConfig().StatusOrderConfig.Canceled
	} else {
		return config.GetConfig().StatusOrderConfig.Approved
	}
}

func (order *Order) FormatOrder() config.OrderJsonStruct {
	var orderItemsJSON []config.OrderItemJsonStruct

	for _, orderItem := range order.OrderItems {
		thumbnailURL := config.ThumbnailURLStruct{
			Path: config.GetConfig().AppConfig.DefaultImageURL,
			Alt:  "default_image",
		}

		if len(orderItem.Product.Images) > 0 {
			thumbnailURL = config.ThumbnailURLStruct{
				Path: orderItem.Product.Images[0].Path,
				Alt:  orderItem.Product.Images[0].Alt,
			}
		}

		orderItemsJSON = append(orderItemsJSON, config.OrderItemJsonStruct{
			Name:         orderItem.Product.Name,
			ThumbnailURL: thumbnailURL,
			Price:        orderItem.Product.Price,
			Quantity:     orderItem.Qty,
		})
	}

	orderJson := config.OrderJsonStruct{
		ID: order.ID,
		Address: config.AddressJsonStruct{
			ID:       order.AddressID,
			Name:     order.Address.Name,
			PostCode: order.Address.PostCode,
		},
		OrderItems:       orderItemsJSON,
		Note:             order.Note,
		ApprovedBy:       order.ApprovedBy,
		ApprovedAt:       order.ApprovedAt,
		CompletedAt:      order.CompletedAt,
		CancelledBy:      order.CancelledBy,
		CancelledAt:      order.CancelledAt,
		CancellationNote: order.CancellationNote,
	}

	return orderJson
}
