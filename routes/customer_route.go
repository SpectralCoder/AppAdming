package routes

import (
	"appadming/controllers"

	"github.com/gin-gonic/gin"
)

func CustomerRoute(router *gin.Engine) {
	router.POST("/customer", controllers.CreateCustomer())
	router.GET("/customer/:customerId", controllers.GetACustomer())
	router.PUT("/customer/:customerId", controllers.EditACustomer())
	router.DELETE("/customer/:customerId", controllers.DeleteACustomer())
	router.GET("/customers", controllers.GetAllCustomers())

}
