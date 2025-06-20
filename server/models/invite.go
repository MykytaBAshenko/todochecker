package models

import "time"

type Invite struct {
	ID        uint `gorm:"primaryKey"`
	SentByID  uint `gorm:"not null"` // references UserGroup.ID
	SentToID  uint `gorm:"not null"` // references User.ID
	GroupID   uint `gorm:"not null"` // references Group.ID
	CreatedAt time.Time
	UpdatedAt time.Time

	SentBy UserGroup `gorm:"foreignKey:SentByID"`
	SentTo User      `gorm:"foreignKey:SentToID"`
	Group  Group     `gorm:"foreignKey:GroupID"`
}
