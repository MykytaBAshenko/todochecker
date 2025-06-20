package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	r.POST("/auth/signup", controllers.Signup)
	r.POST("/auth/login", controllers.Login)
	r.GET("/auth/validate-token", controllers.ValidateToken)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware()) // Protect endpoints below
	{
		protected.GET("/user", controllers.GetUser)
	}
}
