package model

import (
	"gorm.io/gorm"
	"time"
)

type ChatMessage struct {
	gorm.Model
	RoomID  uint      `gorm:"not null;index"`
	UserID  uint      `gorm:"not null;index"`
	Message string    `gorm:"type:text;not null"`
	SentAt  time.Time `gorm:"not null"`
}
