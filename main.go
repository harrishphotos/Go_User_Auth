package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"book_crud/config"
	"book_crud/database"
	"book_crud/middleware"
	"book_crud/routes"
	"book_crud/services/kafka"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load config
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Config loading failed: %v", err)
	}

	// Initialize db
	if err := database.ConnectDB(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Initialize Kafka consumer for email processing
	initKafkaConsumer()

	app := fiber.New()

	// Configure middleware
	middleware.SetupMiddleware(app)

	// Setup routes
	routes.SetupRoutes(app)

	// Get port from config, with fallback options
	port := config.AppConfig.Port
	if port == "" {
		port = "3000"
	}
	
	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

// Add this function to initialize and start the Kafka consumer
func initKafkaConsumer() {
	// Create a context that can be cancelled for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	
	// Set up signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	
	// Create and start the consumer
	consumer := kafka.NewConsumer()
	go consumer.Start(ctx)
	
	// Handle shutdown signal
	go func() {
		<-signalChan
		log.Println("Shutting down application...")
		cancel() // Cancel the context to signal the consumer to stop
		// Allow some time for cleanup
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()
}
