package kafka

import (
	"auth-service/internal/constant"
	"auth-service/internal/dto"
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishVerificationEvent(ctx context.Context, email string, token string) error
	Close() error
}

type producer struct {
	writer *kafka.Writer
}

func New(brokers []string, topic string) Producer {
	return &producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.Hash{},
		},
	}
}

func (p *producer) PublishVerificationEvent(ctx context.Context, email string, token string) error {
	event := dto.VerifyEmailEvent{
		Email: email,
		Token: token,
		Type:  constant.EventTypeVerifyEmail,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return errors.New(constant.ErrMarshalRequest)
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(email),
		Value: payload,
	}); err != nil {
		return errors.New(constant.ErrPublishEvent)
	}

	return nil
}

func (p *producer) Close() error {
	return p.writer.Close()
}
