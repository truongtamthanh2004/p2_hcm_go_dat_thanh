package main

import (
	"context"
	"fmt"
	"log"
	_ "map-service/docs"
	"map-service/internal/config"
	"map-service/internal/handler"
	"map-service/internal/route"
	"map-service/internal/service"
	"map-service/internal/usecase"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Could not load .env file, using system environment variables")
	}

	ctx := context.Background()
	config.InitRedis(ctx)

	venueService := service.NewVenueService()
	mapUsecase := usecase.NewMapUsecase(venueService)
	mapHandler := handler.NewMapHandler(mapUsecase)

	r := route.SetupRouter(mapHandler)

	port := os.Getenv("MAP_SERVICE_PORT")
	if port == "" {
		port = "8088"
	}

	fmt.Printf("🚀 Map Service running on port %s\n", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
