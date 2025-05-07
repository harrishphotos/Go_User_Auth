package utils

import (
	"book_crud/config"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	TokenType string `json:"token_type"` // "access" for access tokens
}

// CreateToken generates a new PASETO token with specified duration
func CreateToken(claims TokenClaims, duration time.Duration) (string, error) {
	v2 := paseto.NewV2()
	
	// Convert hex string to bytes
	key, err := hex.DecodeString(config.AppConfig.PasetoSecretKey)
    if err != nil {
        return "", fmt.Errorf("invalid hex format: %v", err)
    }

    if len(key) != 32 {
        return "", fmt.Errorf("invalid key length: got %d bytes, need 32", len(key))
    }
	
	// Token expires based on provided duration
	expiration := time.Now().Add(duration)
	
	// Create a JSONToken with proper expiration
	jsonToken := paseto.JSONToken{
		Expiration: expiration,
	}
	
	// Encrypt the token with the claims and expiration
	token, err := v2.Encrypt(key, claims, &jsonToken)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %v", err)
	}
	
	return token, nil
}

// CreateAccessToken creates a short-lived access token (15 minutes)
func CreateAccessToken(userID uint, username, role string) (string, error) {
	claims := TokenClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		TokenType: "access",
	}
	
	// Token expiration is set to 1 minute for testing
	// For production, change this to 15*time.Minute
	return CreateToken(claims, 1*time.Minute)
}

// GenerateRefreshToken creates a cryptographically secure random string
// for use as a refresh token (doesn't encode any user information)
func GenerateRefreshToken() (string, error) {
	// Generate 32 bytes of random data
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	
	// Encode as Base64URL (URL-safe Base64)
	return base64.URLEncoding.EncodeToString(b), nil
}

// ValidateToken validates and decodes a PASETO token
func ValidateToken(tokenString string) (*TokenClaims, error) {
	v2 := paseto.NewV2()
	
	// Convert hex string to bytes
	key, err := hex.DecodeString(config.AppConfig.PasetoSecretKey)
    if err != nil {
        return nil, fmt.Errorf("invalid hex format: %v", err)
    }

    if len(key) != 32 {
        return nil, fmt.Errorf("invalid key length: got %d bytes, need 32", len(key))
    }
	
	var claims TokenClaims
	// Create a footer validation rule with proper parsing
	var footer paseto.JSONToken
	
	// Decrypt with footer validation
	err = v2.Decrypt(tokenString, key, &claims, &footer)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	
	// Manually check expiration if needed
	if footer.Expiration.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}
	
	return &claims, nil
}
