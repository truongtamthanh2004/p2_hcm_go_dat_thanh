package main

import (
	"chat-service/internal/db"
	"chat-service/router"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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
