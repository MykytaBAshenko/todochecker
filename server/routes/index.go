package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAllRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	RegisterUserRoutes(r)
	RegisterGroupRoutes(r)
	RegisterTaskRoutes(r)
	RegisterInviteRoutes(r)
	RegisterUserGroupRoutes(r)
	RegisterUserMessageRoutes(r)
}
