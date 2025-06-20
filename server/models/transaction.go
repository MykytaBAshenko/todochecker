package models

import "time"

type Transaction struct {
	ID        uint          `gorm:"primaryKey"`
	GroupID   uint          `gorm:"not null"`
	Status    BillingStatus `gorm:"not null"`
	Amount    float64       `gorm:"not null"`
	BilledAt  time.Time     `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Group     Group `gorm:"foreignKey:GroupID"`
}
