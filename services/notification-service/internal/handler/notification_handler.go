package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"notification-service/config"
	"notification-service/internal/usecase"
	"strconv"
)

type NotificationHandler struct {
	usecase usecase.NotificationUsecase
}

func NewNotificationHandler(u usecase.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{u}
}

func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req struct {
		UserID  uint   `json:"user_id"`
		Type    string `json:"type"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	notif, err := h.usecase.SendNotification(req.UserID, req.Type, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Push realtime to WebSocket
	config.SendToUser(req.UserID, notif)

	c.JSON(http.StatusOK, gin.H{"message": "notification.send_success", "data": notif})
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "notification.invalid_user_id"})
		return
	}
	notifications, err := h.usecase.GetUserNotifications(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "notification.get_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "notification.get_success", "data": notifications})
}
