package models

import "time"

type GroupMessage struct {
	ID          uint   `gorm:"primaryKey"`
	Text        string `gorm:"not null"`
	UserGroupID uint   `gorm:"not null"`
	CreatedAt   time.Time
	UserGroup   UserGroup `gorm:"foreignKey:UserGroupID"`
}
