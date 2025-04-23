package utils

import (
	"book_crud/config"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// CreateToken generates a new PASETO token
func CreateToken(claims TokenClaims) (string, error) {
	v2 := paseto.NewV2()
	
	// Convert hex string to bytes
	key, err := hex.DecodeString(config.AppConfig.PasetoSecretKey)
	if err != nil || len(key) != 32 {
		return "", fmt.Errorf("invalid secret key: must be 32 bytes")
	}
	
	// Token expires based on config
	expiration := time.Now().Add(config.AppConfig.TokenExpiration)
	
	token, err := v2.Encrypt(key, claims, expiration)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %v", err)
	}
	
	return token, nil
}

// ValidateToken validates and decodes a PASETO token
func ValidateToken(tokenString string) (*TokenClaims, error) {
	v2 := paseto.NewV2()
	
	// Convert hex string to bytes
	key, err := hex.DecodeString(config.AppConfig.PasetoSecretKey)
	if err != nil || len(key) != 32 {
		return nil, fmt.Errorf("invalid secret key: must be 32 bytes")
	}
	
	var claims TokenClaims
	err = v2.Decrypt(tokenString, key, &claims, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	
	return &claims, nil
}
