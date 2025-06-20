package routes

import "github.com/gin-gonic/gin"

func RegisterAllRoutes(r *gin.Engine) {
	RegisterUserRoutes(r)
	RegisterGroupRoutes(r)
	RegisterTaskRoutes(r)
	RegisterInviteRoutes(r)
	RegisterUserGroupRoutes(r)
	RegisterUserMessageRoutes(r)
}
