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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	baseURL := os.Getenv("USER_SERVICE_URL")
	if baseURL == "" {
		log.Fatal("missing env: USER_SERVICE_URL")
	}
	chatRepo := repository.NewChatRepository(db)
	userClient := repository.NewUserClient(baseURL)
	chatUC := usecase.NewChatUsecase(chatRepo, userClient)
	hub := ws.NewHub(chatUC)
	go hub.Run()
	chatHandler := handler.NewChatHandler(chatUC, hub)
	chatApi := r.Group("api/v1/chat")
	chatApi.GET("/ws", chatHandler.SendMessage)
	chatApi.GET("/conversations/:user2", middleware.RequireAuth(), chatHandler.GetConversation)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
