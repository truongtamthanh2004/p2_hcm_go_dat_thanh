package router

import (
	"chat-service/internal/handler"
	"chat-service/internal/middleware"
	"chat-service/internal/repository"
	"chat-service/internal/usecase"
	ws "chat-service/internal/websocket"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		log.Fatal("missing env: APP_BASE_URL")
	}
	chatRepo := repository.NewChatRepository(db)
	userClient := repository.NewUserClient(baseURL)
	chatUC := usecase.NewChatUsecase(chatRepo, userClient)
	hub := ws.NewHub(chatUC)
	go hub.Run()
	chatHandler := handler.NewChatHandler(chatUC, hub)
	chatApi := r.Group("")
	chatApi.GET("/ws", chatHandler.SendMessage)
	chatApi.GET("/conversations/:user2", middleware.RequireAuth(), chatHandler.GetConversation)
}
