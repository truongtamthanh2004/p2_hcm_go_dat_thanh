package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"notification-service/internal/model"
	"notification-service/internal/usecase"
)

// Mock repo
type MockNotificationRepo struct {
	mock.Mock
}

func (m *MockNotificationRepo) Create(notification *model.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}
func (m *MockNotificationRepo) GetByUserID(userID uint) ([]model.Notification, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Notification), args.Error(1)
}
func (m *MockNotificationRepo) MarkAsRead(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestSendNotification(t *testing.T) {
	mockRepo := new(MockNotificationRepo)
	uc := usecase.NewNotificationUsecase(mockRepo)

	notif := &model.Notification{
		UserID:  1,
		Type:    "IN_APP",
		Content: "Hello!",
		IsRead:  false,
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.Notification")).Return(nil)

	result, err := uc.SendNotification(1, "IN_APP", "Hello!")

	assert.NoError(t, err)
	assert.Equal(t, notif.UserID, result.UserID)
	assert.Equal(t, notif.Type, result.Type)
	assert.Equal(t, notif.Content, result.Content)
	assert.False(t, result.IsRead)
	mockRepo.AssertExpectations(t)
}

func TestGetUserNotifications(t *testing.T) {
	mockRepo := new(MockNotificationRepo)
	uc := usecase.NewNotificationUsecase(mockRepo)

	expected := []model.Notification{
		{UserID: 1, Content: "Hello"},
		{UserID: 1, Content: "World"},
	}

	mockRepo.On("GetByUserID", uint(1)).Return(expected, nil)

	res, err := uc.GetUserNotifications(1)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "Hello", res[0].Content)
	mockRepo.AssertExpectations(t)
}

func TestMarkAsRead(t *testing.T) {
	mockRepo := new(MockNotificationRepo)
	uc := usecase.NewNotificationUsecase(mockRepo)

	mockRepo.On("MarkAsRead", uint(1)).Return(nil)

	err := uc.MarkAsRead(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSendNotification_Error(t *testing.T) {
	mockRepo := new(MockNotificationRepo)
	uc := usecase.NewNotificationUsecase(mockRepo)

	mockRepo.On("Create", mock.AnythingOfType("*model.Notification")).Return(errors.New("db error"))

	result, err := uc.SendNotification(1, "IN_APP", "Hello")

	assert.Error(t, err)
	assert.Nil(t, result)
}
