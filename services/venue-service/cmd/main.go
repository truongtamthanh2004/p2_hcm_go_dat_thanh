package main

import (
	"fmt"
	"os"
	"p2_hcm_go_dat_thanh/services/venue-service/config"
	"p2_hcm_go_dat_thanh/services/venue-service/internal/handler"
	"p2_hcm_go_dat_thanh/services/venue-service/internal/repository"
	"p2_hcm_go_dat_thanh/services/venue-service/internal/route"
	"p2_hcm_go_dat_thanh/services/venue-service/internal/usecase"
)

func main() {
	config.ConnectDB()

	repo := repository.NewVenueRepository(config.DB)
	uc := usecase.NewVenueUsecase(repo)
	h := handler.NewVenueHandler(uc)

	r := route.SetupRouter(h)

	port := os.Getenv("VENUE_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	r.Run(fmt.Sprintf(":%s", port))
}
