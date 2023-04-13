package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "OKOK",
	})
}

func Register() {

}

func ForgotPassword() {

}

func VerifyPassword() {

}

func ChangePassword() {

}
