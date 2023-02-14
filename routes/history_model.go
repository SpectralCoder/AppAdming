package routes

import (
	"appadming/controllers"

	"github.com/gin-gonic/gin"
)

func HistoryRoute(router *gin.Engine) {
	router.POST("/history", controllers.CreateHistory())
	router.GET("/historys/:historyId", controllers.GetAHistory())
	router.PUT("/historys/:historyId", controllers.EditAHistory())
	router.DELETE("/historys/:historyId", controllers.DeleteAHistory())
	router.GET("/historys", controllers.GetAllHistorys())
}
