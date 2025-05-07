package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"book_crud/config"
	"book_crud/services"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
    reader       *kafka.Reader
    emailService *services.EmailService
}

func NewConsumer() *Consumer {
    brokers := strings.Split(config.AppConfig.KAFKABrokers, ",")
    topic := config.AppConfig.KAFKATopicsendverification
    groupID := config.AppConfig.KAFKAGroupid

    return &Consumer{
        reader: kafka.NewReader(kafka.ReaderConfig{
            Brokers: brokers,
            Topic:   topic,
            GroupID: groupID,
        }),
        emailService: services.NewEmailService(),
    }
}

func (c *Consumer) Start(ctx context.Context) {
    // Use a separate goroutine to handle graceful shutdown
    go func() {
        <-ctx.Done()
        log.Println("Shutting down Kafka consumer...")
        c.Close() // This doesn't return a value, so don't assign it
    }()
    
    log.Printf("Starting Kafka consumer for topic: %s", c.reader.Config().Topic)
    
    // Keep processing messages until context is cancelled
    for {
        select {
        case <-ctx.Done():
            return
        default:
            msg, err := c.reader.FetchMessage(ctx)
            if err != nil {
                if ctx.Err() != nil {
                    // Context was cancelled, exit gracefully
                    return
                }
                log.Printf("Kafka fetch error: %v", err)
                time.Sleep(time.Second) // Add a small delay before retrying
                continue
            }
            
            var payload EmailPayload
            if err := json.Unmarshal(msg.Value, &payload); err != nil {
                log.Printf("Invalid payload: %v", err)
                // Commit message to avoid reprocessing invalid messages
                if err := c.reader.CommitMessages(ctx, msg); err != nil {
                    log.Printf("Failed to commit invalid message: %v", err)
                }
                continue
            }
            
            // Process the message
            if err := c.processEmailMessage(payload); err != nil {
                log.Printf("Email processing error: %v - will retry later", err)
                // Don't commit the message so it will be reprocessed
                continue
            }
            
            // Commit the successfully processed message
            if err := c.reader.CommitMessages(ctx, msg); err != nil {
                log.Printf("Commit failed: %v", err)
            }
        }
    }
}

// processEmailMessage handles different types of email messages
func (c *Consumer) processEmailMessage(payload EmailPayload) error {
    // Default to verification if type is empty (for backward compatibility)
    if payload.MessageType == "" {
        payload.MessageType = "verification"
    }
    
    switch payload.MessageType {
    case "verification":
        return c.emailService.SendVerificationEmail(payload.Email, payload.Username, payload.Token)
    case "password_reset":
        return c.emailService.SendPasswordResetEmail(payload.Email, payload.Username, payload.Token)
    default:
        return fmt.Errorf("unknown message type: %s", payload.MessageType)
    }
}

// Close closes the Kafka reader
func (c *Consumer) Close() {
    if c.reader != nil {
        if err := c.reader.Close(); err != nil {
            log.Printf("Failed to close Kafka reader: %v", err)
        } else {
            log.Println("Kafka reader closed successfully")
        }
    }
}
