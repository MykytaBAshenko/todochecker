package models

import "time"

type Task struct {
	ID         uint   `gorm:"primaryKey"`
	TaskBody   string `gorm:"not null"`
	AssignedTo uint   `gorm:"not null"`
	GroupID    uint   `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserGroup  UserGroup `gorm:"foreignKey:AssignedTo"`
	Group      Group     `gorm:"foreignKey:GroupID"`
}
