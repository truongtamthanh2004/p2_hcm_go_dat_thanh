package main

import (
	"log"
	"mail-service/internal/config"
	"mail-service/internal/kafka"
	"mail-service/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	cfg := config.LoadConfig()
	mailSender := utils.NewMailSender(cfg)
	log.Println("Mail Service started...")
	kafka.StartConsumer(cfg, mailSender)
}
