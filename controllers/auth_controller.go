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

	token, err := utils.CreateToken(utils.TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		fmt.Printf("Verification token is missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "verification failed",
		})
	}

	tx := database.DB.Begin()
	defer tx.Rollback()

	var user models.User
	result := tx.Where("verification_token = ? AND token_expires_at > ?", token, time.Now()).First(&user)

	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or Expired link",
		})
	}

	user.IsVerified = true
	user.VerificationToken = ""

	if err := tx.Save(&user).Error; err != nil {
		fmt.Printf("error while saving user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong. Please try again later.",
		})
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Printf("Transaction commit failed: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong. Please try again later.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Email verified successfully. You can now log in.",
	})
}
