package routes

import (
	"appadming/controllers"

	"github.com/gin-gonic/gin"
)

func SellsRoute(router *gin.Engine) {
	router.POST("/sell", controllers.CreateSell())
	router.GET("/sells/:sellsId", controllers.GetASell())
	router.PUT("/sells/:sellsId", controllers.EditASell())
	router.DELETE("/sells/:sellsId", controllers.DeleteASell())
	router.GET("/sells", controllers.GetAllSells())
}
