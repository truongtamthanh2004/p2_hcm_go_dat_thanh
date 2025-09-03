package main

import (
	"booking-service/config"
	"booking-service/internal/handler"
	"booking-service/internal/kafka"
	"booking-service/internal/repository"
	"booking-service/internal/router"
	"booking-service/internal/service"
	"booking-service/internal/usecase"
	"fmt"
	"os"
)

func main() {
	config.ConnectDB()

	venueServiceDomain := os.Getenv("VENUE_SERVICE_URL")
	if venueServiceDomain == "" {
		venueServiceDomain = "http://venue-service:8083"
	}
	repo := repository.NewBookingRepository(config.DB)
	venueSvc := service.NewVenueHTTPService(venueServiceDomain)

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}
	topic := os.Getenv("KAFKA_TOPIC_NOTIFICATION_EVENTS")
	if topic == "" {
		topic = "notification-events"
	}
	producer := kafka.NewProducer(
		[]string{brokers},
		topic,
	)
	uc := usecase.NewBookingUsecase(repo, venueSvc, producer)
	h := handler.NewBookingHandler(uc)

	r := router.SetupRouter(h)

	port := os.Getenv("BOOKING_SERVICE_PORT")
	if port == "" {
		port = "8084"
	}

	r.Run(fmt.Sprintf(":%s", port))
}
