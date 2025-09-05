package main

import (
	"fmt"
	"notification-service/config"
	"notification-service/internal/handler"
	"notification-service/internal/kafka"
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

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}
	topic := os.Getenv("KAFKA_TOPIC_NOTIFICATION_EVENTS")
	if topic == "" {
		topic = "notification-events"
	}
	group := os.Getenv("KAFKA_CONSUMER_GROUP_NOTIFICATION_SERVICE")
	if group == "" {
		group = "notification-service"
	}

	kafka.StartBookingConsumer(
		[]string{brokers},
		topic,
		group,
		uc,
	)

	r := route.SetupRouter(h)

	port := os.Getenv("NOTIFICATION_SERVICE_PORT")
	if port == "" {
		port = "8087"
	}

	r.Run(fmt.Sprintf(":%s", port))
}
