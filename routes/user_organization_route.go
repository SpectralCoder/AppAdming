package routes

import (
	controller "appadming/controllers"

	"github.com/gin-gonic/gin"
)

// UserOrganizationRoutes function
func UserOrganizationRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/org/:org_id", controller.CreateUserOrganization())
	incomingRoutes.GET("/user/orgs", controller.GetAllUserOrganizations())
	incomingRoutes.PUT("/approve/role/:id", controller.ApproveRole())
	incomingRoutes.GET("/user/orgs/:org_id", controller.GetUserOfOrganization())
	incomingRoutes.GET("/org/me/:org_id", controller.GetMyOrganization())
}