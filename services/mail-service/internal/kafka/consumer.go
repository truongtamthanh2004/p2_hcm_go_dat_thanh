package kafka

import (
	"context"
	"encoding/json"
	"log"
	"mail-service/internal/config"
	"mail-service/internal/constant"
	"mail-service/internal/utils"

	"github.com/segmentio/kafka-go"
)

type MailEvent struct {
	Type  string `json:"type"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func StartConsumer(cfg *config.MailConfig, sender *utils.MailSender) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaBroker},
		Topic:   cfg.KafkaMailTopic,
		GroupID: constant.MailServiceGroup,
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		var event MailEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Println("Invalid event format:", err)
			continue
		}

		switch event.Type {
		case constant.EventTypeVerifyEmail:
			sender.SendVerificationEmail(event.Email, event.Token)
		default:
			log.Println("Unknown mail type:", event.Type)
		}
	}
}
