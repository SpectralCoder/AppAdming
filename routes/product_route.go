package routes

import (
	"appadming/controllers"

	"github.com/gin-gonic/gin"
)

func ProductRoute(router *gin.Engine) {
	router.POST("/product", controllers.CreateProduct())
	router.GET("/products/:productId", controllers.GetAProduct())
	router.PUT("/products/:productId", controllers.EditAProduct())
	router.DELETE("/products/:productId", controllers.DeleteAProduct())
	router.GET("/products", controllers.GetAllProducts())
}
