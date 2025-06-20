package controllers

import (
	"net/http"
	"time"

	"server/config"
	"server/models"
	"server/utils"

	"github.com/gin-gonic/gin"
)

// CreateTask creates a new task
func CreateTask(c *gin.Context) {
	var input struct {
		TaskBody string `json:"task_body" binding:"required"`
		GroupID  uint   `json:"group_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Get user ID from JWT context
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Find the UserGroup ID for the current user in the specified group
	var userGroup models.UserGroup
	if err := config.DB.Where("user_id = ? AND group_id = ?", userID, input.GroupID).First(&userGroup).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not part of the group"})
		return
	}

	// Create the task using the UserGroup ID
	task := models.Task{
		TaskBody:   input.TaskBody,
		GroupID:    input.GroupID,
		AssignedTo: userGroup.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetAllTasks returns all tasks
func GetAllTasks(c *gin.Context) {
	var tasks []models.Task
	if err := config.DB.Preload("UserGroup.User").Preload("Group").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTask returns a specific task by ID
func GetTask(c *gin.Context) {
	var task models.Task
	if err := config.DB.Preload("UserGroup.User").Preload("Group").
		First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// UpdateTask edits a task by ID
func UpdateTask(c *gin.Context) {
	var task models.Task
	if err := config.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input struct {
		TaskBody   string `json:"task_body"`
		AssignedTo uint   `json:"assigned_to"`
		GroupID    uint   `json:"group_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.TaskBody != "" {
		task.TaskBody = input.TaskBody
	}
	if input.AssignedTo != 0 {
		task.AssignedTo = input.AssignedTo
	}
	if input.GroupID != 0 {
		task.GroupID = input.GroupID
	}
	task.UpdatedAt = time.Now()

	if err := config.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask deletes a task by ID
func DeleteTask(c *gin.Context) {
	if err := config.DB.Delete(&models.Task{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
