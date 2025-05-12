package controllers

import (
	"book_crud/database"
	"book_crud/models"
	"book_crud/services/kafka"
	"book_crud/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"email": req.Email,
		}).Error("Password hashing failed")

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong. Please try again later.",
		})
	}

	token, err := utils.GenerateVerificationToken()
	if err != nil {
		fmt.Printf("Failed to generate verification token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong. Please try again later.",
		})
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	tx := database.DB.Begin()
	defer tx.Rollback()

	user := models.User{
		Email:             req.Email,
		Username:          req.Username,
		Password:          hashedPassword,
		Role:              "user",
		IsVerified:        false,
		VerificationToken: token,
		TokenExpiresAt:    expiresAt,
	}

	if err := tx.Create(&user).Error; err != nil {
		fmt.Printf("User creation failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create user",
		})
	}

	producer := kafka.NewProducer()
	defer producer.Close()

	if err := producer.SendVerificationEmail(user.Email, user.Username, token, expiresAt); err != nil {
		fmt.Printf("Failed to enqueue verification email: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send verification email. Please try again later.",
		})
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Printf("Transaction commit failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong. Please try again later.",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully. Please check your email for verification.",
	})
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if !user.IsVerified {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":              "Please verify your email before logging in",
			"needsVerification": true,
		})
	}

	match, err := utils.ComparePasswords(req.Password, user.Password)
	if err != nil || !match {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Generate access token
	accessToken, err := utils.CreateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"userId": user.ID,
		}).Error("Failed to create access token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create access token",
		})
	}
	
	// Generate refresh token (now a random string)
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"userId": user.ID,
		}).Error("Failed to generate refresh token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create refresh token",
		})
	}
	
	// Store refresh token in the database
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days
	user.RefreshToken = refreshToken
	user.RefreshExpiresAt = refreshExpiresAt
	
	if err := database.DB.Save(&user).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"userId": user.ID,
		}).Error("Failed to save refresh token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not save refresh token",
		})
	}
	
	// Set refresh token as HttpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  refreshExpiresAt,
		HTTPOnly: true,
		Secure:   true, // Set to false for development to work with HTTP
		SameSite: "lax", // Set to lax to allow redirects
		Path:     "/api/auth", // Restrict to auth endpoints
	})

	// Return only the access token in the response body
	return c.JSON(fiber.Map{
		"access_token": accessToken,
		"expires_in": 60, // 15 minutes in seconds
		"user": fiber.Map{
			"id": user.ID,
			"username": user.Username,
			"email": user.Email,
			"role": user.Role,
		},
	})
}

// RefreshToken generates a new access token using a refresh token
func RefreshToken(c *fiber.Ctx) error {
	// Get refresh token from cookie
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing refresh token",
		})
	}
	
	// Find user by refresh token
	var user models.User
	if err := database.DB.Where("refresh_token = ?", refreshToken).First(&user).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Invalid refresh token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}
	
	// Check if refresh token is expired
	if time.Now().After(user.RefreshExpiresAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Refresh token expired",
		})
	}
	
	// Generate new access token
	accessToken, err := utils.CreateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"userId": user.ID,
		}).Error("Failed to create new access token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create access token",
		})
	}
	
	// Optional: Implement refresh token rotation for better security
	// This creates a new refresh token with each use
	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"userId": user.ID,
		}).Error("Failed to generate new refresh token")
		// Continue with the existing refresh token
	} else {
		// Update the refresh token in the database
		user.RefreshToken = newRefreshToken
		user.RefreshExpiresAt = time.Now().Add(7 * 24 * time.Hour)
		
		if err := database.DB.Save(&user).Error; err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"userId": user.ID,
			}).Error("Failed to save new refresh token")
			// Continue with the existing refresh token
		} else {
			// Set the new refresh token in the cookie
			c.Cookie(&fiber.Cookie{
				Name:     "refresh_token",
				Value:    newRefreshToken,
				Expires:  user.RefreshExpiresAt,
				HTTPOnly: true,
				Secure:   false, // Set to false for development
				SameSite: "lax", // Changed to lax for better compatibility
				Path:     "/api/auth",
			})
		}
	}
	
	// Return the new access token
	return c.JSON(fiber.Map{
		"access_token": accessToken,
		"expires_in": 900, // 15 minutes in seconds
	})
}

// Logout invalidates the refresh token
func Logout(c *fiber.Ctx) error {
	// Get the user ID from the authenticated user (this is a protected route)
	user := c.Locals("user").(*utils.TokenClaims)
	
	// Clear the refresh token in the database
	if err := database.DB.Model(&models.User{}).Where("id = ?", user.UserID).
		Updates(map[string]interface{}{
			"refresh_token": "",
			"refresh_expires_at": time.Now(),
		}).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"userId": user.UserID,
		}).Error("Failed to clear refresh token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not complete logout",
		})
	}
	
	// Clear the refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expired
		HTTPOnly: true,
		Secure:   false, // Set to false for development
		SameSite: "lax",
		Path:     "/api/auth",
	})
	
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		fmt.Printf("Verification token is missing")
		return c.Redirect("http://localhost:5173/verify-email?status=failed")
	}

	tx := database.DB.Begin()
	defer tx.Rollback()

	var user models.User
	result := tx.Where("verification_token = ? AND token_expires_at > ?", token, time.Now()).First(&user)

	if result.Error != nil || result.RowsAffected == 0 {
		return c.Redirect("http://localhost:5173/verify-email?status=failed")
	}

	user.IsVerified = true
	user.VerificationToken = ""
	user.TokenExpiresAt= time.Time{}

	if err := tx.Save(&user).Error; err != nil {
		fmt.Printf("error while saving user: %v", err)
		return c.Redirect("http://localhost:5173/verify-email?status=failed")
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Printf("Transaction commit failed: %v\n", err)
		return c.Redirect("http://localhost:5173/verify-email?status=failed")
	}

	// Redirect to the frontend verification page with success status
	return c.Redirect("http://localhost:5173/verify-email?status=success")
}

// ForgotPassword sends a password reset email to the user
func ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Don't reveal that the email doesn't exist (for security)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "If your email exists in our system, you will receive a password reset link shortly.",
		})
	}

	// Generate reset token
	token, err := utils.GenerateVerificationToken()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"email": req.Email,
		}).Error("Failed to generate reset token")
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong. Please try again later.",
		})
	}

	// Set token expiry (1 hour)
	expiresAt := time.Now().Add(1 * time.Hour)

	// Update user with reset token
	tx := database.DB.Begin()
	defer tx.Rollback()

	user.ResetPasswordToken = token
	user.ResetTokenExpiresAt = expiresAt

	if err := tx.Save(&user).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"userId": user.ID,
		}).Error("Failed to save reset token")
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong. Please try again later.",
		})
	}

	// Send reset password email via Kafka
	producer := kafka.NewProducer()
	defer producer.Close()

	if err := producer.SendPasswordResetEmail(user.Email, user.Username, token, expiresAt); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"userId": user.ID,
		}).Error("Failed to send reset password email")
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send reset password email. Please try again later.",
		})
	}

	if err := tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"userId": user.ID,
		}).Error("Transaction commit failed")
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong. Please try again later.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "If your email exists in our system, you will receive a password reset link shortly.",
	})
}

// ResetPassword resets the user's password using the reset token
func ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if req.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Reset token is required",
		})
	}

	if req.NewPassword == "" || len(req.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters",
		})
	}

	tx := database.DB.Begin()
	defer tx.Rollback()

	// Find user by reset token and check if token is valid
	var user models.User
	result := tx.Where("reset_password_token = ? AND reset_token_expires_at > ?", req.Token, time.Now()).First(&user)
	
	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or expired reset token",
		})
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"userId": user.ID,
		}).Error("Password hashing failed")
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong. Please try again later.",
		})
	}

	// Update user password and clear reset token
	user.Password = hashedPassword
	user.ResetPasswordToken = ""
	user.ResetTokenExpiresAt = time.Time{}

	if err := tx.Save(&user).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"userId": user.ID,
		}).Error("Failed to update password")
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong. Please try again later.",
		})
	}

	if err := tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"userId": user.ID,
		}).Error("Transaction commit failed")
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong. Please try again later.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password has been reset successfully. You can now log in with your new password.",
	})
}