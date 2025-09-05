package handler

import (
	"chat-service/internal/constant"
	"chat-service/internal/usecase"
	"chat-service/internal/websocket"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatUsecase usecase.ChatUsecase
	hub         *websocket.Hub
}

func NewChatHandler(chatUsecase usecase.ChatUsecase, hub *websocket.Hub) *ChatHandler {
	return &ChatHandler{chatUsecase, hub}
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	if err := h.hub.HandleWebSocket(c.Writer, c.Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
}

func (h *ChatHandler) GetConversation(c *gin.Context) {
	user1Val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized})
		return
	}
	user1, ok := user1Val.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrInvalidUserID})
		return
	}

	user2, err := strconv.Atoi(c.Param("user2"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidUserID})
		return
	}
	if user1 == uint(user2) {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrSameUserConversation})
		return
	}

	messages, err := h.chatUsecase.GetConversation(c.Request.Context(), uint(user1), uint(user2))
	if err != nil {
		switch err.Error() {
		case constant.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": constant.ErrUserNotFound})
		case constant.ErrFailedToGetConversation:
			c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrFailedToGetConversation})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrInternalServer})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": constant.SuccessGetConversation,
		"data":    messages,
	})
}
