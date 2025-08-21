package main

import (
	"fmt"
	"notification-service/config"
	"notification-service/internal/handler"
	"notification-service/internal/repository"
	"notification-service/internal/route"
	"notification-service/internal/usecase"
	"os"
)

func main() {
	config.ConnectDB()

	repo := repository.NewNotificationRepository(config.DB)
	uc := usecase.NewNotificationUsecase(repo)
	h := handler.NewNotificationHandler(uc)

	r := route.SetupRouter(h)

	port := os.Getenv("NOTIFICATION_SERVICE_PORT")
	if port == "" {
		port = "8087"
	}

	r.Run(fmt.Sprintf(":%s", port))
}
