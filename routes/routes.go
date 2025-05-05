package routes

import (
	"book_crud/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(app *fiber.App) {
	// Auth routes - unprotected
	SetupAuthRoutes(app)
	
	// Create protected API group
	api := app.Group("/api")
	api.Use(middleware.AuthMiddleware())
	
	// Setup protected routes
	SetupBookRoutes(api)
	SetupStoreRoutes(api)
} 