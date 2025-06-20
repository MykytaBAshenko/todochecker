package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserGroupRoutes(r *gin.Engine) {
	usergroup := r.Group("/usergroup")
	usergroup.Use(middleware.AuthMiddleware())
	usergroup.POST("/invites/:id/accept", controllers.AcceptInvite)
	usergroup.DELETE("/group/leave", controllers.LeaveGroup)
	usergroup.DELETE("/group/remove-user", controllers.AdminRemoveUserFromGroup)
	usergroup.PUT("/group/promote", controllers.PromoteUserToAdmin)
}
