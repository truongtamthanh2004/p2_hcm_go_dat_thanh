package config

import (
	"log"
	"os"
	"strconv"
)

type MailConfig struct {
	FromEmail      string
	Password       string
	Host           string
	Port           int
	AuthServiceURL string
	KafkaBroker    string
	KafkaMailTopic string
}

func LoadConfig() *MailConfig {
	port, err := strconv.Atoi(os.Getenv("FROM_EMAIL_SMTP_PORT"))
	if err != nil {
		log.Fatal("Invalid SMTP port")
	}
	cfg := &MailConfig{
		FromEmail:      os.Getenv("FROM_EMAIL"),
		Password:       os.Getenv("FROM_EMAIL_PASSWORD"),
		Host:           os.Getenv("FROM_EMAIL_SMTP_HOST"),
		Port:           port,
		AuthServiceURL: os.Getenv("AUTH_SERVICE_URL"),
		KafkaBroker:    os.Getenv("KAFKA_BROKERS"),
		KafkaMailTopic: os.Getenv("KAFKA_TOPIC_VERIFY_EMAIL"),
	}
	if cfg.FromEmail == "" || cfg.Password == "" || cfg.Host == "" || cfg.AuthServiceURL == "" || cfg.KafkaBroker == "" || cfg.KafkaMailTopic == "" {
		log.Fatal("missing required email environment variables")
	}
	return cfg
}
