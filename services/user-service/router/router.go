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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(r *gin.Engine) {
	baseURL := os.Getenv("AUTH_SERVICE_URL")
	if baseURL == "" {
		log.Fatal("missing env: AUTH_SERVICE_URL")
	}
	authClient := repository.NewAuthClient(baseURL)
	userRepo := repository.NewUserRepository(db.DB)
	userUC := usecase.NewUserUsecase(userRepo, authClient)
	userHandler := handler.NewUserHandler(userUC)
	// User routes
	api := r.Group("api/v1/users")
	//admin
	api.GET("/", middleware.RequireAuth("admin" , "moderator"), userHandler.GetUserList)
	api.GET("/:id", userHandler.GetUserByID)
	api.PUT("/:id", middleware.RequireAuth("admin"), userHandler.UpdateUser)
	//user
	api.GET("/profile", middleware.RequireAuth("user"), userHandler.GetUserProfile)
	api.PUT("/profile", middleware.RequireAuth("user"), userHandler.UpdateUserProfile)

	//auth-service
	api.POST("/", userHandler.CreateUser)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
