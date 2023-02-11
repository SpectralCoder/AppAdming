package routes

import (
	"appadming/controllers"

	"github.com/gin-gonic/gin"
)

func ProductRoute(router *gin.Engine) {
	router.POST("/product", controllers.CreateProduct())
	router.GET("/product/:productId", controllers.GetAProduct())
	router.PUT("/product/:productId", controllers.EditAProduct())
	router.DELETE("/product/:productId", controllers.DeleteAProduct())
	router.GET("/products", controllers.GetAllProducts())
}
