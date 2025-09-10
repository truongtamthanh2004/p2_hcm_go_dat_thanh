package main

import (
	_ "auth-service/docs"
	"auth-service/internal/db"
	"auth-service/internal/kafka"
	"auth-service/router"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Auth Service API
// @version 1.0
// @description This is the authentication service for the coworking booking system.

// @host localhost:8081
// @BasePath /api/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.InitDB()
	db.AutoMigrate()
	r := gin.Default()

	kafkaBrokers := os.Getenv("KAFKA_BROKERS") // format: "broker1:9092,broker2:9092"
	if kafkaBrokers == "" {
		log.Fatal("missing env: KAFKA_BROKERS")
	}
	brokerList := strings.Split(kafkaBrokers, ",")

	kafkaTopic := os.Getenv("KAFKA_TOPIC_VERIFY_EMAIL")
	if kafkaTopic == "" {
		log.Fatal("missing env: KAFKA_TOPIC_VERIFY_EMAIL")
	}
	producer := kafka.New(brokerList, kafkaTopic)
	defer producer.Close()

	router.SetupRouter(r, db.DB, producer)
	err = r.Run(":8081")
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
