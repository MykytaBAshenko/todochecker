package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Print raw JWT parts
		parts := strings.Split(tokenStr, ".")
		if len(parts) == 3 {
			headerBytes, _ := base64.RawURLEncoding.DecodeString(parts[0])
			payloadBytes, _ := base64.RawURLEncoding.DecodeString(parts[1])
			signature := parts[2]

			fmt.Println("🔐 JWT Header:")
			fmt.Println(string(headerBytes))
			fmt.Println("📦 JWT Payload:")
			fmt.Println(string(payloadBytes))
			fmt.Println("✍️ JWT Signature:")
			fmt.Println(signature)
		} else {
			fmt.Println("⚠️ JWT token does not have 3 parts!")
		}

		// Parse and verify token
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// Extract and inject user ID
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := uint(claims["user_id"].(float64))
			c.Set("user_id", userID)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
		}
	}
}
