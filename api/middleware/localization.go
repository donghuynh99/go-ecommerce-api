package middleware

import (
	"net/http"

	"github.com/donghuynh99/ecommerce_api/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func CheckLocalization() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale, err := c.Request.Cookie("locale")

		if err != nil || !slices.Contains([]string{"vi", "en"}, locale.Value) {
			locale = &http.Cookie{
				Name:   "locale",
				Value:  "en",
				MaxAge: 86400,
				Path:   "/",
			}

			http.SetCookie(c.Writer, locale)
		}

		utils.InitLocalizer(locale.Value)

		c.Next()
	}
}
