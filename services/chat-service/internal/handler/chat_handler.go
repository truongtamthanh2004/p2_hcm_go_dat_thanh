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

// SendMessage godoc
// @Summary Connect to chat WebSocket
// @Description Establish a WebSocket connection for real-time messaging
// @Tags Chat
// @Param user_id query int true "User ID"  // nếu bạn muốn query param user_id
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /chat/ws [get]
func (h *ChatHandler) SendMessage(c *gin.Context) {
	if err := h.hub.HandleWebSocket(c.Writer, c.Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
}

// GetConversation godoc
// @Summary Get conversation between two users
// @Description Retrieve all messages between the authenticated user and another user
// @Tags Chat
// @Produce json
// @Param user2 path int true "ID of the other user in conversation"
// @Success 200 {object} map[string]interface{} "List of messages"
// @Failure 400 {object} map[string]string "Invalid user ID or same user"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/conversations/{user2} [get]
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
