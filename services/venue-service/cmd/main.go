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

	repo := repository.NewVenueRepository(config.DB)
	uc := usecase.NewVenueUsecase(repo)
	h := handler.NewVenueHandler(uc)

	spaceRepository := repository.NewSpaceRepository(config.DB)
	spaceUsecase := usecase.NewSpaceUsecase(spaceRepository)
	spaceHandler := handler.NewSpaceHandler(spaceUsecase)

	r := route.SetupRouter(h, spaceHandler)

	port := os.Getenv("VENUE_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	r.Run(fmt.Sprintf(":%s", port))
}
