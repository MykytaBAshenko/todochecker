package models

import "time"

type UserMessage struct {
	ID              uint   `gorm:"primaryKey"`
	MessageSender   uint   `gorm:"not null"`                  // Foreign key to User
	MessageReceiver uint   `gorm:"not null"`                  // Foreign key to User
	MessageBody     string `gorm:"type:text;not null"`        // Supports long text
	MessageType     string `gorm:"type:varchar(10);not null"` // "string" or "file"
	IsRead          bool   `gorm:"default:false"`
	CreatedAt       time.Time

	// Associations
	Sender   User `gorm:"foreignKey:MessageSender;constraint:OnDelete:CASCADE"`
	Receiver User `gorm:"foreignKey:MessageReceiver;constraint:OnDelete:CASCADE"`
}
