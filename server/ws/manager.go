package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"server/config"
	"server/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development; restrict in production
		return true
	},
}

var connectedUsers = make(map[uint]*websocket.Conn) // userID -> connection
var mu sync.Mutex

func WebSocketHandler(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}
	userID := uint(claims["user_id"].(float64))

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	connectedUsers[userID] = conn
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(connectedUsers, userID)
		mu.Unlock()
		fmt.Printf("WebSocket connection closed for user %d\n", userID)
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Read error from user %d: %v\n", userID, err)
			break
		}

		var raw map[string]interface{}
		if err := json.Unmarshal(msgBytes, &raw); err != nil {
			fmt.Println("Invalid message format:", err)
			continue
		}

		switch raw["type"] {
		case "initial_message":
			handleInitialMessage(userID, "initial_message", raw, conn)
		case "incoming_message":
			handleDefaultMessage(userID, "incoming_message", raw, conn)
		// case "delete_message":
		// 	handleDeleteMessage(userID, raw, conn)
		case "delete_conversation":
			handleDeleteConversation(userID, raw, conn)
		default:
			fmt.Println("Unknown message type")
		}
	}
}

func handleInitialMessage(senderID uint, messType string, raw map[string]interface{}, conn *websocket.Conn) {
	to := raw["to"].(string)
	body := raw["body"].(string)
	msgType := raw["msgType"].(string)

	receiver, err := FindUserByNicknameOrEmail(to)
	if err != nil {
		conn.WriteJSON(gin.H{"error": "Receiver not found"})
		return
	}

	saveAndSendMessage(senderID, messType, receiver.ID, body, msgType)
}

func handleDefaultMessage(senderID uint, messType string, raw map[string]interface{}, conn *websocket.Conn) {
	toID := uint(raw["to"].(float64))
	body := raw["body"].(string)
	msgType := raw["msgType"].(string)

	saveAndSendMessage(senderID, messType, toID, body, msgType)
}

func saveAndSendMessage(senderID uint, messType string, receiverID uint, body, msgType string) {
	msg := models.UserMessage{
		MessageSender:   senderID,
		MessageReceiver: receiverID,
		MessageBody:     body,
		MessageType:     msgType,
		CreatedAt:       time.Now(),
		IsRead:          false,
	}

	if err := config.DB.Create(&msg).Error; err != nil {
		fmt.Println("DB save error:", err)
		return
	}

	var fullMsg models.UserMessage
	if err := config.DB.Preload("Sender").Preload("Receiver").
		First(&fullMsg, msg.ID).Error; err != nil {
		fmt.Println("DB fetch error:", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	payload := map[string]interface{}{
		"type": messType,
		"data": fullMsg,
	}

	if conn, ok := connectedUsers[receiverID]; ok {
		conn.WriteJSON(payload)
	}
	if conn, ok := connectedUsers[senderID]; ok {
		conn.WriteJSON(payload)
	}
}

func FindUserByNicknameOrEmail(input string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("nickname = ? OR email = ?", input, input).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func handleDeleteConversation(userID uint, raw map[string]interface{}, conn *websocket.Conn) {
	otherUserIDFloat, ok := raw["user_id"].(float64)
	if !ok {
		conn.WriteJSON(gin.H{"error": "Invalid user_id"})
		return
	}
	otherUserID := uint(otherUserIDFloat)
	if err := config.DB.
		Where("(message_sender = ? AND message_receiver = ?) OR (message_sender = ? AND message_receiver = ?)",
			userID, otherUserID, otherUserID, userID).
		Delete(&models.UserMessage{}).Error; err != nil {
		conn.WriteJSON(gin.H{"error": "Failed to delete conversation"})
		return
	}

	conn.WriteJSON(gin.H{
		"type": "delete_conversation",
		"data": gin.H{
			"userID": otherUserID,
		},
	})

	mu.Lock()
	otherConn, isOnline := connectedUsers[otherUserID]
	mu.Unlock()

	if isOnline {
		otherConn.WriteJSON(gin.H{
			"type": "delete_conversation",
			"data": gin.H{
				"userID": userID,
			},
		})
	}
}
