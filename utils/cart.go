package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetCartDataFromCookie(c *gin.Context) map[string]int {
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

	return data
}
