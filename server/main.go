package main

import (
	"server/config"
	"server/routes"
	"server/ws"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterWS(r *gin.Engine) {
	// Other routes...

	// WebSocket route
	r.GET("/ws", ws.WebSocketHandler)
}

func main() {
	// Connect to the database
	config.ConnectDatabase()
	r := gin.Default()

	// Custom CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5174"}, // your frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // REQUIRED for cookies or Authorization headers
		MaxAge:           12 * time.Hour,
	}))

	routes.RegisterAllRoutes(r)
	RegisterWS(r)

	// Run the server
	r.Run(":8080")
}
