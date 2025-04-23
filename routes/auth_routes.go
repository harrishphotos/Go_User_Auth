package routes

import (
	"github.com/gofiber/fiber/v2"
	"book_crud/controllers"
)

func SetupAuthRoutes(app *fiber.App) {
	// Auth routes are unprotected, so we use the main app
	auth := app.Group("/api/auth")

	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
} 