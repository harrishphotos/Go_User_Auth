package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"book_crud/config"
	"book_crud/database"
	"book_crud/middleware"
	"book_crud/routes"
)

func main() {
	// Load environment variables
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Initialize database connection
	if err := database.ConnectDB(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Create new Fiber app
	app := fiber.New()

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Setup auth routes (unprotected)
	routes.SetupAuthRoutes(app)
	
	// Create protected API group
	api := app.Group("/api")
	api.Use(middleware.AuthMiddleware())
	
	// Pass the protected api group to route setup functions
	routes.SetupBookRoutes(api)
	routes.SetupStoreRoutes(api)

	// Start server
	log.Fatal(app.Listen(":" + config.AppConfig.Port))
}
