package main

import (
	"log"
	"user-service/db"
	"user-service/router"

	_ "user-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title       User Service API
// @version     1.0
// @description This is the User Service for the coworking booking system.
// @termsOfService http://swagger.io/terms/

// @contact.name  API Support
// @contact.url   http://www.example.com/support
// @contact.email support@example.com

// @host        localhost:8082
// @BasePath    /api/v1
// @schemes     http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.InitDB()
	db.AutoMigrate()
	r := gin.Default()
	router.SetupRouter(r)
	err = r.Run(":8082")
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
