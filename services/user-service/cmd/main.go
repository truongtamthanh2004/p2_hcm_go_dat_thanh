package main

import (
	"log"
	"user-service/db"
	"user-service/router"

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
	router.SetupRouter(r)
	err = r.Run(":8082")
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
