package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBUser         string
	DBName         string
	DBPort         string
	DBSSLMode      string
	PasetoSecretKey string
	TokenExpiration time.Duration
	Port           string
	Environment    string
	SMTPHost       string
	SMTPPort       string
	SMTPUsername   string
	SMTPPassword   string
	SMTPFrom       string
	APPUrl         string
}

var AppConfig Config

func LoadConfig() error {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return err
	}

	// Parse token expiration duration
	tokenExpiration, err := time.ParseDuration(os.Getenv("TOKEN_EXPIRATION"))
	if err != nil {
		tokenExpiration = 24 * time.Hour // Default to 24 hours
	}

	AppConfig = Config{
		DBHost:          os.Getenv("DB_HOST"),
		DBUser:          os.Getenv("DB_USER"),
		DBName:          os.Getenv("DB_NAME"),
		DBPort:          os.Getenv("DB_PORT"),
		DBSSLMode:       os.Getenv("DB_SSL_MODE"),
		PasetoSecretKey: os.Getenv("PASETO_SECRET_KEY"),
		TokenExpiration: tokenExpiration,
		Port:            os.Getenv("PORT"),
		Environment:     os.Getenv("ENVIRONMENT"),
		SMTPHost:        os.Getenv("SMTP_HOST"),
	    SMTPPort:        os.Getenv("SMTP_PORT"),
	    SMTPUsername:    os.Getenv("SMTP_USERNAME"),
	    SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
	    SMTPFrom:        os.Getenv("SMTP_FROM"),
		APPUrl:          os.Getenv("APP_URL"),

	}

	return nil
} 