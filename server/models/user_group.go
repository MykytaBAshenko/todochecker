package models

import "time"

type UserGroup struct {
	ID        uint `gorm:"primaryKey"`
	IsAdmin   bool `gorm:"default:false"`
	UserID    uint `gorm:"not null"`
	GroupID   uint `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User  `gorm:"foreignKey:UserID"`
	Group     Group `gorm:"foreignKey:GroupID"`
}
