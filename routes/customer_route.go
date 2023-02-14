package routes

import (
	"appadming/controllers"

	"github.com/gin-gonic/gin"
)

func CustomerRoute(router *gin.Engine) {
	router.POST("/customer", controllers.CreateCustomer())
	router.GET("/customers/:customerId", controllers.GetACustomer())
	router.PUT("/customers/:customerId", controllers.EditACustomer())
	router.DELETE("/customers/:customerId", controllers.DeleteACustomer())
	router.GET("/customers", controllers.GetAllCustomers())

}
