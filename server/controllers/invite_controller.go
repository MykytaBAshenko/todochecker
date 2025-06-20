package controllers

import (
	"net/http"

	"server/config"
	"server/models"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func InviteUserToGroup(c *gin.Context) {
	var input struct {
		GroupID  uint   `json:"group_id" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	requesterID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 🧠 Check if the user is an admin of the group
	var requesterUG models.UserGroup
	if err := config.DB.
		Where("user_id = ? AND group_id = ?", requesterID, input.GroupID).
		First(&requesterUG).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this group"})
		return
	}

	if !requesterUG.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only group admins can invite users"})
		return
	}

	// 🔎 Find the user to invite
	var userToInvite models.User
	if err := config.DB.Where("nickname = ?", input.Nickname).First(&userToInvite).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 🛑 Prevent duplicate invites
	var existingInvite models.Invite
	if err := config.DB.
		Where("sent_to_id = ? AND group_id = ?", userToInvite.ID, input.GroupID).
		First(&existingInvite).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already invited to this group"})
		return
	}

	// ✅ Create the invite
	invite := models.Invite{
		SentByID: requesterUG.ID,
		SentToID: userToInvite.ID,
		GroupID:  input.GroupID,
	}

	if err := config.DB.Create(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User invited successfully"})
}

func GetInvites(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var invites []models.Invite
	if err := config.DB.
		Preload("SentBy.User").
		Preload("Group").
		Where("sent_to_id = ?", userID).
		Find(&invites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invites"})
		return
	}

	c.JSON(http.StatusOK, invites)
}

func GetSentInvites(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Step 1: Get all UserGroup IDs where user is a member (or optionally admin)
	var userGroupIDs []uint
	if err := config.DB.
		Model(&models.UserGroup{}).
		Where("user_id = ?", userID).
		Pluck("id", &userGroupIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user groups"})
		return
	}

	// Step 2: Find all invites where SentByID is in those UserGroup IDs
	var invites []models.Invite
	if err := config.DB.
		Preload("SentTo").
		Preload("Group").
		Where("sent_by_id IN ?", userGroupIDs).
		Find(&invites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sent invites"})
		return
	}

	c.JSON(http.StatusOK, invites)
}

func DeleteReceivedInvite(c *gin.Context) {
	inviteID := c.Param("id")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var invite models.Invite
	if err := config.DB.First(&invite, inviteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invite not found"})
		return
	}

	if invite.SentToID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this invite"})
		return
	}

	if err := config.DB.Delete(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invite deleted successfully"})
}

func DeleteSentInvite(c *gin.Context) {
	inviteID := c.Param("id")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get all UserGroup IDs for current user
	var userGroupIDs []uint
	if err := config.DB.
		Model(&models.UserGroup{}).
		Where("user_id = ?", userID).
		Pluck("id", &userGroupIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user groups"})
		return
	}

	// Check if invite exists and was sent by user's group
	var invite models.Invite
	if err := config.DB.
		Where("id = ? AND sent_by_id IN ?", inviteID, userGroupIDs).
		First(&invite).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this invite"})
		return
	}

	if err := config.DB.Delete(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sent invite deleted successfully"})
}
