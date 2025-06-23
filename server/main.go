package main

import (
	"fmt"
	"os"
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
	fmt.Println("🚀 Server starting")
	// Connect to the database
	config.ConnectDatabase()
	r := gin.Default()

	// Custom CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // ⚠️ Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // REQUIRED for cookies or Authorization headers
		MaxAge:           12 * time.Hour,
	}))

	routes.RegisterAllRoutes(r)
	RegisterWS(r)

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "80" // fallback default
	}
	r.Run(":" + port)
}
