package controllers

import (
	"fmt"
	"net/http"
	"server/config"
	"server/models"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EditMessage updates the content of a message by its ID
func EditMessage(c *gin.Context) {
	idParam := c.Param("id")
	messageID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var input struct {
		MessageBody string `json:"message_body" binding:"required"`
		MessageType string `json:"message_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var message models.UserMessage
	if err := config.DB.First(&message, messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	message.MessageBody = input.MessageBody
	message.MessageType = input.MessageType

	if err := config.DB.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully"})
}

// DeleteMessage deletes a specific message by ID
func DeleteMessage(c *gin.Context) {
	idParam := c.Param("id")
	messageID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	if err := config.DB.Delete(&models.UserMessage{}, messageID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

// GetUserMessages fetches all messages sent or received by the logged-in user
func GetUserMessages(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var messages []models.UserMessage
	if err := config.DB.
		Where("message_sender = ? OR message_receiver = ?", userID, userID).
		Order("created_at ASC").
		Preload("Sender").
		Preload("Receiver").
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	// Group messages by the conversation partner
	grouped := make(map[string][]models.UserMessage)
	for _, msg := range messages {
		var otherID uint
		if msg.MessageSender == userID {
			otherID = msg.MessageReceiver
		} else {
			otherID = msg.MessageSender
		}
		key := fmt.Sprintf("%d", otherID)
		grouped[key] = append(grouped[key], msg)
	}

	c.JSON(http.StatusOK, gin.H{"messages": grouped})
}

func GetUserConversations(c *gin.Context) {
	// Step 1: Get current user ID from context
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	// Step 2: Use raw SQL to find distinct user IDs from conversations (sender/receiver)
	var otherUserIDs []uint
	if err := config.DB.Raw(`
		SELECT DISTINCT message_receiver AS id FROM user_messages WHERE message_sender = ?
		UNION
		SELECT DISTINCT message_sender AS id FROM user_messages WHERE message_receiver = ?
	`, userID, userID).Scan(&otherUserIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversation IDs"})
		return
	}

	// Step 3: Filter out self ID if included
	var filteredIDs []uint
	for _, id := range otherUserIDs {
		if id != userID {
			filteredIDs = append(filteredIDs, id)
		}
	}

	// Step 4: Fetch user details for all unique conversation partners
	var users []models.User
	if err := config.DB.Where("id IN ?", filteredIDs).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
		return
	}

	// ✅ Step 5: Return just the array
	c.JSON(http.StatusOK, users)
}

// DeleteAllMessages deletes all messages in the database
func DeleteAllMessages(c *gin.Context) {
	if err := config.DB.Where("1 = 1").Delete(&models.UserMessage{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete all messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All messages deleted"})
}
