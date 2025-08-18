package main

import (
	"auth-service/internal/db"
	"auth-service/internal/kafka"
	"auth-service/router"
	"log"
	"os"
	"strings"

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
