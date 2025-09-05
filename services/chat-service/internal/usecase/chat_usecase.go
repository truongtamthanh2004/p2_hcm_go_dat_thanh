package usecase

import (
	"chat-service/internal/constant"
	"chat-service/internal/model"
	"chat-service/internal/repository"
	"context"
	"errors"
)

type ChatUsecase interface {
	SaveMessage(ctx context.Context, msg *model.ChatMessage) error
	GetConversation(ctx context.Context, user1, user2 uint) ([]model.ChatMessage, error)
}

type chatUsecase struct {
	chatRepo   repository.ChatRepository
	userClient repository.UserClient
}

func NewChatUsecase(chatRepo repository.ChatRepository, userClient repository.UserClient) ChatUsecase {
	return &chatUsecase{chatRepo: chatRepo, userClient: userClient}
}

func (u *chatUsecase) SaveMessage(ctx context.Context, msg *model.ChatMessage) error {
	err := u.chatRepo.SaveMessage(ctx, msg)
	if err != nil {
		return errors.New(constant.ErrFailedToSaveMessage)
	}
	return nil
}

func (u *chatUsecase) GetConversation(ctx context.Context, user1, user2 uint) ([]model.ChatMessage, error) {
	_, err := u.userClient.GetUserByID(ctx, user2)
	if err != nil {
		if err.Error() == constant.ErrUserNotFound {
			return nil, errors.New(constant.ErrUserNotFound)
		}
		return nil, errors.New(constant.ErrInternalServer)
	}
	return u.chatRepo.GetConversation(ctx, user1, user2)
}
