package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"notification-service/internal/handler"
	"notification-service/internal/model"
	"testing"
)

// Mock Usecase
type MockNotificationUsecase struct {
	mock.Mock
}

func (m *MockNotificationUsecase) SendNotification(userID uint, notifType, content string) (*model.Notification, error) {
	args := m.Called(userID, notifType, content)
	return args.Get(0).(*model.Notification), args.Error(1)
}

func (m *MockNotificationUsecase) GetUserNotifications(userID uint) ([]model.Notification, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Notification), args.Error(1)
}

func (m *MockNotificationUsecase) MarkAsRead(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestSendNotificationHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(MockNotificationUsecase)
	h := handler.NewNotificationHandler(mockUC)

	// mock return
	expected := &model.Notification{UserID: 1, Type: "IN_APP", Content: "Hello"}
	mockUC.On("SendNotification", uint(1), "IN_APP", "Hello").Return(expected, nil)

	// request body
	body := `{"user_id":1,"type":"IN_APP","content":"Hello"}`
	req, _ := http.NewRequest("POST", "/api/v1/notifications", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/api/v1/notifications", h.SendNotification)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.Notification
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expected.Content, resp.Content)

	mockUC.AssertExpectations(t)
}

func TestGetNotificationsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(MockNotificationUsecase)
	h := handler.NewNotificationHandler(mockUC)

	// mock return
	expected := []model.Notification{
		{UserID: 1, Type: "IN_APP", Content: "Hi"},
		{UserID: 1, Type: "EMAIL", Content: "Welcome"},
	}
	mockUC.On("GetUserNotifications", uint(1)).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/api/v1/notifications/1", nil)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/api/v1/notifications/:userId", h.GetNotifications)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []model.Notification
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 2)
	assert.Equal(t, "Hi", resp[0].Content)

	mockUC.AssertExpectations(t)
}
