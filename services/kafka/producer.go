package kafka

import (
	"book_crud/config"
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
    writer *kafka.Writer
}

func NewProducer() *Producer {
    brokers := strings.Split(config.AppConfig.KAFKABrokers, ",")
    topic := config.AppConfig.KAFKATopicsendverification

    return &Producer{
        writer: &kafka.Writer{
            Addr:     kafka.TCP(brokers...),
            Topic:    topic,
            Balancer: &kafka.LeastBytes{},
            RequiredAcks: kafka.RequireOne,
            WriteTimeout: 2 * time.Second,
        },
    }
}

func (p *Producer) SendVerificationEmail(email, username, token string, expiresAt time.Time) error {
    payload := EmailPayload{
        Email: email, 
        Username: username, 
        Token: token, 
        ExpiresAt: expiresAt,
        MessageType: "verification",
    }

    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    msg := kafka.Message{
        Key:   []byte(email),
        Value: data,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() {
    if err := p.writer.Close(); err != nil {
        log.Printf("Failed to close Kafka writer: %v", err)
    }
}
