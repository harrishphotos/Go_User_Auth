package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupMiddleware configures middleware for the application
func SetupMiddleware(app *fiber.App) {
	// Enable CORS with proper configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173", 
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Enable request logging
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
}