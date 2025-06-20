package controllers

import (
	"net/http"
	"server/config"
	"server/models"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func AcceptInvite(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to accept this invite"})
		return
	}

	// Check if the user is already in the group
	var existingUG models.UserGroup
	if err := config.DB.
		Where("user_id = ? AND group_id = ?", userID, invite.GroupID).
		First(&existingUG).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "You are already a member of this group"})
		return
	}

	// Create the UserGroup entry
	userGroup := models.UserGroup{
		UserID:  userID,
		GroupID: invite.GroupID,
		IsAdmin: false,
	}

	if err := config.DB.Create(&userGroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join group"})
		return
	}

	// Delete the invite after accepting
	if err := config.DB.Delete(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "You have joined the group"})
}

func AdminRemoveUserFromGroup(c *gin.Context) {
	var input struct {
		GroupID uint `json:"group_id" binding:"required"`
		UserID  uint `json:"user_id" binding:"required"` // target user
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	adminID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if requester is admin in group
	var adminUG models.UserGroup
	if err := config.DB.
		Where("user_id = ? AND group_id = ?", adminID, input.GroupID).
		First(&adminUG).Error; err != nil || !adminUG.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can remove users"})
		return
	}

	// Don't allow admin to remove themselves
	if adminID == input.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admins cannot remove themselves using this route"})
		return
	}

	// Remove target user from group
	if err := config.DB.
		Where("user_id = ? AND group_id = ?", input.UserID, input.GroupID).
		Delete(&models.UserGroup{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from group"})
}

func LeaveGroup(c *gin.Context) {
	var input struct {
		GroupID uint `json:"group_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if the user is actually in the group
	var userGroup models.UserGroup
	if err := config.DB.
		Where("user_id = ? AND group_id = ?", userID, input.GroupID).
		First(&userGroup).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "You are not a member of this group"})
		return
	}

	// Delete the UserGroup record
	if err := config.DB.
		Delete(&userGroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave the group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "You have left the group successfully"})
}

func PromoteUserToAdmin(c *gin.Context) {
	var input struct {
		GroupID uint `json:"group_id" binding:"required"`
		UserID  uint `json:"user_id" binding:"required"`
	}

	// Bind the incoming request JSON to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Get the current user's ID from the context (decoded from JWT)
	adminID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if the current user is an admin in the specified group
	var adminUserGroup models.UserGroup
	if err := config.DB.Where("user_id = ? AND group_id = ?", adminID, input.GroupID).First(&adminUserGroup).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this group"})
		return
	}

	if !adminUserGroup.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not an admin in this group"})
		return
	}

	// Check if the target user is already a member of the group
	var targetUserGroup models.UserGroup
	if err := config.DB.Where("user_id = ? AND group_id = ?", input.UserID, input.GroupID).First(&targetUserGroup).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User is not a member of this group"})
		return
	}

	// Update the target user's IsAdmin status to true
	targetUserGroup.IsAdmin = true

	// Save the updated UserGroup record
	if err := config.DB.Save(&targetUserGroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote user to admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully promoted to admin"})
}
