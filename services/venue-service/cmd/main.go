package main

import (
	"fmt"
	"os"
	"venue-service/config"
	"venue-service/internal/handler"
	"venue-service/internal/repository"
	"venue-service/internal/route"
	"venue-service/internal/usecase"
)

func main() {
	config.ConnectDB()

	venueRepository := repository.NewVenueRepository(config.DB)
	venueUsecase := usecase.NewVenueUsecase(venueRepository)
	venueHandler := handler.NewVenueHandler(venueUsecase)

	spaceRepository := repository.NewSpaceRepository(config.DB)
	spaceUsecase := usecase.NewSpaceUsecase(spaceRepository, venueRepository)
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
