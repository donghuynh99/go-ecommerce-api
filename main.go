package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/donghuynh99/ecommerce_api/api/routes"
	"github.com/donghuynh99/ecommerce_api/cmd"
	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/database"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

var (
	developmentMode bool
)

func init() {
	flag.BoolVar(&developmentMode, "dev", false, "run in development mode")
	flag.Parse()
}

func main() {
	if developmentMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	config := config.Load()
	database.Connect()
	router := routes.SetupRoutes()
	flag.Parse()
	arg := flag.Arg(0)

	if arg != "" {
		cmd.InitCommand()
	} else {
		router.Run(":" + strconv.Itoa(config.AppConfig.AppPort))
	}
}
