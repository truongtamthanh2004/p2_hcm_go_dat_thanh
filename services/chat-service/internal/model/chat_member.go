package model

import (
	"gorm.io/gorm"
	"time"
)

type ChatMember struct {
	gorm.Model
	RoomID   uint      `gorm:"not null;index"`
	UserID   uint      `gorm:"not null;index"`
	JoinedAt time.Time `gorm:"not null"`
}
