package routes

import (
	AdminController "github.com/donghuynh99/ecommerce_api/api/controllers/admin"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	AdminRoutes(router)
	UserRoutes(router)

	return router
}

func AdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	{
		admin.POST("/login", AdminController.Login)
		admin.POST("/register", AdminController.Register)

		password := admin.Group("/password")
		{
			password.POST("/forgot", AdminController.ForgotPassword)
			password.POST("/change", AdminController.ChangePassword)
		}
	}
}

func UserRoutes(router *gin.Engine) {

}
