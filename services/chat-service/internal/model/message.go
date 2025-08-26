package model

import (
	"gorm.io/gorm"
)

type ChatMessage struct {
	gorm.Model
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
}
