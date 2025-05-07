package routes

import (
	"book_crud/controllers"
	"book_crud/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	// Auth routes are unprotected, so we use the main app
	auth := app.Group("/api/auth")

	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Get("/verify-email", controllers.VerifyEmail)
	auth.Post("/refresh", controllers.RefreshToken)
	auth.Post("/forgot-password", controllers.ForgotPassword)
	auth.Post("/reset-password", controllers.ResetPassword)
	
	// Protected route - requires authentication because this logout clears out the refresh token 
	auth.Post("/logout", middleware.AuthMiddleware(), controllers.Logout)
} 