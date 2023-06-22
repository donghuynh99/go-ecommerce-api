package routes

import (
	"time"

	AdminController "github.com/donghuynh99/ecommerce_api/api/controllers/admin"
	ClientController "github.com/donghuynh99/ecommerce_api/api/controllers/client"
	"github.com/donghuynh99/ecommerce_api/api/middleware"
	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// CORs setup
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Static("/assets", "./assets")

	router.Use(middleware.CheckLocalization())

	apiRouter := router.Group("/api")

	apiRouter.GET("/change-locale", ClientController.ChangeLocale)

	AdminRoutes(apiRouter)
	UserRoutes(apiRouter)

	return router
}

func AdminRoutes(router *gin.RouterGroup) {
	adminController := AdminController.NewController()

	adminAuth := router.Group("/admin")
	{
		adminAuth.POST("/login", adminController.Login)

		adminAuth.POST("/logout", middleware.HandleAuthenticate(config.GetConfig().RoleConfig.Admin), adminController.Logout)
	}

	admin := router.Group("/admin", middleware.HandleAuthenticate(config.GetConfig().RoleConfig.Admin))
	{
		product := admin.Group("/products")
		{
			product.GET("/", adminController.ListProduct)
			product.GET("/:id", adminController.ShowProduct)
			product.POST("/", adminController.StoreProduct)
			product.PUT("/:id", adminController.UpdateProduct)
			product.PATCH("/:id/status", adminController.UpdateStatusProduct)
			product.DELETE("/:id", adminController.DeleteProduct)
		}

		category := admin.Group("/categories")
		{
			category.GET("/", adminController.ListCategories)
			category.POST("/", adminController.StoreCategories)
			category.PUT("/:id", adminController.UpdateCategory)
			category.DELETE("/:id", adminController.DeleteCategory)
			category.DELETE("/:id/force-delete", adminController.ForceDeleteCategory)
		}
	}
}

func UserRoutes(router *gin.RouterGroup) {
	clientController := ClientController.NewController()
	adminController := AdminController.NewController()

	router.POST("/login", adminController.Login)
	router.POST("/logout", middleware.HandleAuthenticate(config.GetConfig().RoleConfig.User), adminController.Logout)
	router.POST("/register", clientController.Register)

	// cart
	router.GET("/cart", clientController.GetCart)
	router.PUT("/cart", clientController.UpdateCart)

	// product
	router.GET("/products/most-popular", clientController.GetPopularProducts)

	// authentication route
	router.Use(middleware.HandleAuthenticate(config.GetConfig().RoleConfig.User))

	router.GET("/profile", clientController.ShowProfile)
	router.PUT("/profile", clientController.UpdateProfile)
	router.PUT("/profile/avatar", clientController.UpdateAvatar)

	address := router.Group("/addresses")
	{
		address.GET("/", clientController.ListAddress)
		address.POST("/", clientController.AddAddress)
		address.PUT("/:id", clientController.UpdateAddress)
		address.DELETE("/:id", clientController.RemoveAddress)
	}

	order := router.Group("/orders")
	{
		order.GET("/", clientController.ListMyOrder)
		order.POST("/", clientController.HandleOrder)
		order.GET("/:id", clientController.ShowOrder)
		order.PATCH("/:id/cancel-order", clientController.CancelOrder)
	}
}
