package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterGroupRoutes(r *gin.Engine) {
	group := r.Group("/group")
	group.Use(middleware.AuthMiddleware())
	group.POST("/create", controllers.CreateGroup)
	group.GET("/groups", controllers.GetUserGroups)
	group.DELETE("/:group_id", controllers.DeleteGroup)

}
