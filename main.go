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

	router.Run("localhost:6000")
}
