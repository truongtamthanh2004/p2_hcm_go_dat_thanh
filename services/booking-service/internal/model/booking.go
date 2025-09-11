package model

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID     uint      `gorm:"not null;index"`
	SpaceID    uint      `gorm:"not null;index"`
	StartTime  time.Time `gorm:"not null"`
	EndTime    time.Time `gorm:"not null"`
	Status     string    `gorm:"type:varchar(50);default:'PENDING'"` // PENDING, CONFIRMED, CANCELED
	TotalPrice float64   `gorm:"type:decimal(10,2);not null"`
}
