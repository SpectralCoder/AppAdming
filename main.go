package main

import (
	"appadming/configs"
	"appadming/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//run database
	configs.ConnectDB()

	//routes
	routes.HistoryRoute(router)
	routes.SellsRoute(router)
	routes.CustomerRoute(router)
	routes.ProductRoute(router)

	err := router.Run(":6000")
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
