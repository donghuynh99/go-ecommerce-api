package utils

import (
	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateUUID() string {
	uuidObj := uuid.New()

	// Convert UUID to string
	uuidStr := uuidObj.String()

	return uuidStr
}

func GetUser(c *gin.Context) (*models.User, bool) {
	user, ok := c.Get("user")

	if !ok {
		return &models.User{}, false
	}

	userModel, ok := user.(*models.User)
	if !ok {
		return &models.User{}, false
	}

	return userModel, true
}
