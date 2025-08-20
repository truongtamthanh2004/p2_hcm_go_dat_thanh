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
	Email string            `json:"email"`
	Type  string            `json:"type"` // "VERIFY_EMAIL"
	Data  map[string]string `json:"data,omitempty"`
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
			token := event.Data["token"]
			if token == "" {
				log.Println("Missing token in verify email event")
				continue
			}
			sender.SendVerificationEmail(event.Email, token)
		case constant.EventTypeResetPassword:
			newPassword := event.Data["newPassword"]
			if newPassword == "" {
				log.Println("Missing new password in reset password email event")
				continue
			}
			sender.SendResetPassword(event.Email, newPassword)
		default:
			log.Println("Unknown mail type:", event.Type)
		}
	}
}
