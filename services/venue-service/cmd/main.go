package main

import (
	"fmt"
	"log"
	"os"
	"venue-service/config"
	_ "venue-service/docs"
	"venue-service/internal/handler"
	"venue-service/internal/repository"
	"venue-service/internal/route"
	"venue-service/internal/usecase"
)

// @title Venue Service API
// @version 1.0
// @description This is the venue service API for the coworking booking system.

// @host localhost:8083
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config.ConnectDB()

	baseURL := os.Getenv("BOOKING_SERVICE_URL")
	if baseURL == "" {
		log.Fatal("missing env: BOOKING_SERVICE_URL")
	}
	bookingClient := repository.NewBookingClient(baseURL)
	venueRepository := repository.NewVenueRepository(config.DB)
	venueUsecase := usecase.NewVenueUsecase(venueRepository)
	venueHandler := handler.NewVenueHandler(venueUsecase)

	spaceRepository := repository.NewSpaceRepository(config.DB)
	spaceUsecase := usecase.NewSpaceUsecase(spaceRepository, venueRepository, bookingClient)
	spaceHandler := handler.NewSpaceHandler(spaceUsecase)

	amenityRepository := repository.NewAmenityRepository(config.DB)
	amenityUsecase := usecase.NewAmenityUsecase(amenityRepository)
	amenityHandler := handler.NewAmenityHandler(amenityUsecase)
	r := route.SetupRouter(venueHandler, spaceHandler, amenityHandler)

	port := os.Getenv("VENUE_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	r.Run(fmt.Sprintf(":%s", port))
}
