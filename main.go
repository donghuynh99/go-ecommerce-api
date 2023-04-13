package main

import (
	"flag"
	"strconv"

	"github.com/donghuynh99/ecommerce_api/api/routes"
	"github.com/donghuynh99/ecommerce_api/cmd"
	"github.com/donghuynh99/ecommerce_api/config"
)

func main() {
	config := config.Load()
	router := routes.SetupRoutes()
	flag.Parse()
	arg := flag.Arg(0)

	if arg != "" {
		cmd.InitCommand()
	} else {
		router.Run(":" + strconv.Itoa(config.AppConfig.AppPort))
	}

	router.Run()
}
