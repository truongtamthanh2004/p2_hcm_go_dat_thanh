package router

import (
	"log"
	"os"
	"user-service/db"
	"user-service/internal/handler"
	"user-service/internal/middleware"
	"user-service/internal/repository"
	"user-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		log.Fatal("missing env: APP_BASE_URL")
	}
	authClient := repository.NewAuthClient(baseURL)
	userRepo := repository.NewUserRepository(db.DB)
	userUC := usecase.NewUserUsecase(userRepo, authClient)
	userHandler := handler.NewUserHandler(userUC)
	// User routes
	api := r.Group("")
	//admin
	api.GET("/", middleware.RequireAuth("admin"), userHandler.GetUserList)
	api.GET("/:id", middleware.RequireAuth("admin"), userHandler.GetUserByID)
	api.PUT("/:id", middleware.RequireAuth("admin"), userHandler.UpdateUser)
	//user
	api.GET("/profile", middleware.RequireAuth("user"), userHandler.GetUserProfile)
	api.PUT("/profile", middleware.RequireAuth("user"), userHandler.UpdateUserProfile)

	//auth-service
	api.POST("/", userHandler.CreateUser)
}
