package controllers

import (
	"net/http"
	"strconv"
	"time"

	"server/config"
	"server/models"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func CreateGroup(c *gin.Context) {
	// Parse group input
	var input struct {
		Name  string `json:"name" binding:"required"`
		Image string `json:"image"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Get user ID from JWT
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Create group
	group := models.Group{
		Name:      input.Name,
		Image:     input.Image,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := config.DB.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	// Create user-group link with admin rights
	userGroup := models.UserGroup{
		IsAdmin:   true,
		UserID:    userID,
		GroupID:   group.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := config.DB.Create(&userGroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link user to group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"group":      group,
		"user_group": userGroup,
	})
}

func GetUserGroups(c *gin.Context) {
	// Step 1: Get the current user's ID from the JWT context
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Step 2: Retrieve all UserGroup records for the current user
	var userGroups []models.UserGroup
	if err := config.DB.
		Preload("Group"). // Preload the associated group data
		Where("user_id = ?", userID).
		Find(&userGroups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user groups"})
		return
	}

	// Step 3: Extract all the groups from the UserGroup records
	var groups []models.Group
	for _, userGroup := range userGroups {
		groups = append(groups, userGroup.Group)
	}

	// Step 4: Return the list of groups the user is part of
	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

func DeleteGroup(c *gin.Context) {
	// Step 1: Get the current user's ID from the JWT context
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Step 2: Get the group ID from the URL parameters
	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Step 3: Check if the user is an admin of the group
	var userGroup models.UserGroup
	if err := config.DB.Where("user_id = ? AND group_id = ?", userID, groupID).First(&userGroup).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this group"})
		return
	}

	// Step 4: Ensure the user is an admin of the group
	if !userGroup.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not an admin of this group"})
		return
	}

	// Step 5: Delete all UserGroup associations for this group (optional, for clean-up)
	if err := config.DB.Where("group_id = ?", groupID).Delete(&models.UserGroup{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove users from group"})
		return
	}

	// Step 6: Delete the group
	if err := config.DB.Where("id = ?", groupID).Delete(&models.Group{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}

	// Step 7: Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}
