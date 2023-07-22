package main

import (
	"appadming/configs"
	middleware "appadming/middlewares"
	"appadming/routes"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	//run database
	configs.ConnectDB()

	router.Use(gin.Logger())
	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5432"} // Replace with the allowed origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	routes.AuthRoutes(router)
	router.Use(middleware.Authentication())

	routes.UserRoutes(router)

	// API-2
	router.GET("/api-1", func(c *gin.Context) {

		c.JSON(200, gin.H{"success": "Access granted for api-1"})

	})

	// API-1
	router.GET("/api-2", func(c *gin.Context) {
		email, exists := c.Get("email")
		if !exists {
			fmt.Println("not found")
		}
		fmt.Println(email)
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	routes.HistoryRoute(router)
	routes.SellsRoute(router)
	routes.CustomerRoute(router)
	routes.ProductRoute(router)

	err := router.Run("0.0.0.0:9000")
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
