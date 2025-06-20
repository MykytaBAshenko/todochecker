package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserMessageRoutes(r *gin.Engine) {
	protected := r.Group("/usermessages")
	protected.Use(middleware.AuthMiddleware()) // Protect endpoints below
	{
		protected.GET("/all", controllers.GetUserMessages)
		protected.GET("/allusers", controllers.GetUserConversations)
	}
}
