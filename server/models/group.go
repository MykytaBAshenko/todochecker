package models

import "time"

type Group struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
