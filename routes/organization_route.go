package routes

import (
	"appadming/controllers"

	"github.com/gin-gonic/gin"
)

func OrganizationRoute(router *gin.Engine) {
	router.POST("/organization", controllers.CreateOrganization())
	router.GET("/organizations/:organizationId", controllers.GetAOrganization())
	router.PUT("/organizations/:organizationId", controllers.EditAOrganization())
	router.DELETE("/organizations/:organizationId", controllers.DeleteAOrganization())
	router.GET("/organizations", controllers.GetAllOrganizations())
}
