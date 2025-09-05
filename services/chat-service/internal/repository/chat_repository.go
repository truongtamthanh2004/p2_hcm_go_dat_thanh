package repository

import (
	"chat-service/internal/model"
	"context"

	"gorm.io/gorm"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, msg *model.ChatMessage) error
	GetConversation(ctx context.Context, user1, user2 uint) ([]model.ChatMessage, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) SaveMessage(ctx context.Context, msg *model.ChatMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *chatRepository) GetConversation(ctx context.Context, user1, user2 uint) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage
	err := r.db.WithContext(ctx).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", user1, user2, user2, user1).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}
