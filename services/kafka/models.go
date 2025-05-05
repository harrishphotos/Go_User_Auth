package kafka

import "time"

// EmailPayload represents the structure of data sent to Kafka for email processing
// It contains all necessary information to generate and send verification emails
type EmailPayload struct {
	// Email is the recipient's email address
	Email string `json:"email"`
	
	// Username is used for personalization in the email
	Username string `json:"username"`
	
	// Token is the verification token to be included in the email link
	Token string `json:"token"`
	
	// ExpiresAt indicates when the verification token expires
	ExpiresAt time.Time `json:"expires_at"`
	
	// MessageType can be used to differentiate between different types of emails
	// (e.g., "verification", "password_reset", "welcome", etc.)
	MessageType string `json:"message_type,omitempty"`
}