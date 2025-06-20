package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(r *gin.Engine) {
	t := r.Group("/tasks")
	{
		t.Use(middleware.AuthMiddleware())
		t.POST("/create", controllers.CreateTask)
		t.GET("/", controllers.GetAllTasks)
		t.GET("/:id", controllers.GetTask)
		t.PUT("/:id", controllers.UpdateTask)
		t.DELETE("/:id", controllers.DeleteTask)
	}
}
