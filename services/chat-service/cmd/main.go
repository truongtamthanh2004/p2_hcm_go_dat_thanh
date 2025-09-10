package main

import (
	_ "chat-service/docs"
	"chat-service/internal/db"
	"chat-service/router"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Chat Service API
// @version 1.0
// @description This is the chat service API for the coworking booking system.


// @host localhost:8086
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.InitDB()
	db.AutoMigrate()
	r := gin.Default()
	router.SetupRouter(r, db.DB)
	err = r.Run(":8086")
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
