package usecase_test

import (
	"chat-service/internal/constant"
	"chat-service/internal/dto"
	"chat-service/internal/model"
	"chat-service/internal/usecase"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock ChatRepository
type MockChatRepo struct {
	mock.Mock
}

func (m *MockChatRepo) SaveMessage(ctx context.Context, msg *model.ChatMessage) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockChatRepo) GetConversation(ctx context.Context, user1, user2 uint) ([]model.ChatMessage, error) {
	args := m.Called(ctx, user1, user2)
	return args.Get(0).([]model.ChatMessage), args.Error(1)
}

// Mock UserClient
type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) GetUserByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.UserResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

//
// Tests
//

func TestSaveMessage_RepoSuccess_ReturnNil(t *testing.T) {
	mockRepo := new(MockChatRepo)
	mockUser := new(MockUserClient)
	uc := usecase.NewChatUsecase(mockRepo, mockUser)

	msg := &model.ChatMessage{SenderID: 1, ReceiverID: 2, Content: "Hello"}
	mockRepo.On("SaveMessage", mock.Anything, msg).Return(nil)

	err := uc.SaveMessage(context.Background(), msg)
	assert.NoError(t, err)
}

func TestSaveMessage_RepoFails_ReturnError(t *testing.T) {
	mockRepo := new(MockChatRepo)
	mockUser := new(MockUserClient)
	uc := usecase.NewChatUsecase(mockRepo, mockUser)

	msg := &model.ChatMessage{SenderID: 1, ReceiverID: 2, Content: "Hello"}
	mockRepo.On("SaveMessage", mock.Anything, msg).Return(errors.New("db error"))

	err := uc.SaveMessage(context.Background(), msg)
	assert.EqualError(t, err, constant.ErrFailedToSaveMessage)
}

func TestGetConversation_UserNotFound(t *testing.T) {
	mockRepo := new(MockChatRepo)
	mockUser := new(MockUserClient)
	uc := usecase.NewChatUsecase(mockRepo, mockUser)

	mockUser.On("GetUserByID", mock.Anything, uint(2)).
		Return(nil, errors.New(constant.ErrUserNotFound))

	msgs, err := uc.GetConversation(context.Background(), 1, 2)
	assert.Nil(t, msgs)
	assert.EqualError(t, err, constant.ErrUserNotFound)
}

func TestGetConversation_UserClientError(t *testing.T) {
	mockRepo := new(MockChatRepo)
	mockUser := new(MockUserClient)
	uc := usecase.NewChatUsecase(mockRepo, mockUser)

	mockUser.On("GetUserByID", mock.Anything, uint(2)).
		Return(nil, errors.New("some error"))

	msgs, err := uc.GetConversation(context.Background(), 1, 2)
	assert.Nil(t, msgs)
	assert.EqualError(t, err, constant.ErrInternalServer)
}

func TestGetConversation_RepoSuccess_ReturnMessages(t *testing.T) {
	mockRepo := new(MockChatRepo)
	mockUser := new(MockUserClient)
	uc := usecase.NewChatUsecase(mockRepo, mockUser)

	// User tồn tại
	mockUser.On("GetUserByID", mock.Anything, uint(2)).
		Return(&dto.UserResponse{ID: 2, Email: "user2@gmail.com"}, nil)

	expected := []model.ChatMessage{
		{SenderID: 1, ReceiverID: 2, Content: "Hi"},
		{SenderID: 2, ReceiverID: 1, Content: "Hello"},
	}

	mockRepo.On("GetConversation", mock.Anything, uint(1), uint(2)).
		Return(expected, nil)

	msgs, err := uc.GetConversation(context.Background(), 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, expected, msgs)
}
