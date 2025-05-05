package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupMiddleware configures middleware for the application
func SetupMiddleware(app *fiber.App) {
	// Enable CORS
	app.Use(cors.New())
	
	// Enable request logging
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
} 