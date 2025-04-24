package routes

import (
	"book_crud/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	// Auth routes are unprotected, so we use the main app
	auth := app.Group("/api/auth")

	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Get("/verify-email", controllers.VerifyEmail)
} 