package middleware

import (
	"net/http"
	"strings"

	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/donghuynh99/ecommerce_api/utils"
	"github.com/gin-gonic/gin"
)

func CheckAuthenticate(c *gin.Context, role string) *models.User {
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		return nil
	}

	splitToken := strings.Split(tokenString, "Bearer ")

	if len(splitToken) != 2 {
		return nil
	}

	token := splitToken[1]

	user, err := utils.ValidateToken(token, role)
	if err != nil {
		return nil
	}

	return &user
}

func HandleAuthenticate(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := CheckAuthenticate(c, role)

		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthenticate",
			})

			return
		}

		c.Set("user", user)

		c.Next()
	}
}
