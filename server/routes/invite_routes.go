package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterInviteRoutes(r *gin.Engine) {
	invite := r.Group("/invite")
	invite.Use(middleware.AuthMiddleware())
	invite.POST("/create", controllers.InviteUserToGroup)
	invite.GET("/invites", controllers.GetInvites)
	invite.GET("/invites/sent", controllers.GetSentInvites)
	invite.DELETE("/invites/received/:id", controllers.DeleteReceivedInvite)
	invite.DELETE("/invites/sent/:id", controllers.DeleteSentInvite)
}
