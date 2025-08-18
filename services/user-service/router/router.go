package router

import (
	"user-service/db"
	"user-service/internal/handler"
	"user-service/internal/repository"
	"user-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	userRepo := repository.NewUserRepository(db.DB)
	userUC := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUC)
	// User routes
	api := r.Group("/api/v1/users")
	api.POST("/", userHandler.CreateUser)
}
