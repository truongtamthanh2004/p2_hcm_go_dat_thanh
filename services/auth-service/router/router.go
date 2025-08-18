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
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		log.Fatal("missing env: USER_SERVICE_URL")
	}

	// Init dependencies
	authRepo := repository.NewAuthRepository(dbConn)
	userClient := repository.NewUserClient(userServiceURL)

	authUC := usecase.NewAuthUsecase(authRepo, userClient, kafkaProducer)
	authHandler := handler.NewAuthHandler(*authUC)

	// Routes
	api := r.Group("/api/v1/auth")
	api.POST("/sign-up", authHandler.SignUp)
	api.GET("/verify-account", authHandler.VerifyAccount)
	api.POST("/login", authHandler.Login)
	api.POST("/refresh-token", authHandler.RefreshToken)
}
