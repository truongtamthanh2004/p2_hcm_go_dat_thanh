package router

import (
	"auth-service/internal/handler"
	"auth-service/internal/kafka"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, dbConn *gorm.DB, kafkaProducer kafka.Producer) {
	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		log.Fatal("missing env: APP_BASE_URL")
	}

	// Init dependencies
	authRepo := repository.NewAuthRepository(dbConn)
	userClient := repository.NewUserClient(baseURL)

	authUC := usecase.NewAuthUsecase(authRepo, userClient, kafkaProducer)
	authHandler := handler.NewAuthHandler(*authUC)

	// Routes
	api := r.Group("")
	api.POST("/sign-up", authHandler.SignUp)
	api.GET("/verify-account", authHandler.VerifyAccount)
	api.POST("/login", authHandler.Login)
	api.POST("/refresh-token", authHandler.RefreshToken)
	api.POST("/reset-password", authHandler.ResetPassword)

	//user-service
	api.PUT("/users", authHandler.UpdateAuthUser)
}
