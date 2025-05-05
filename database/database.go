package database

import (
	"fmt"

	"book_crud/config"
	"book_crud/models"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
		config.AppConfig.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto Migrate both models
	err = db.AutoMigrate(&models.Book{}, &models.Store{}, &models.User{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	fmt.Println("Connected to database successfully")
	DB = db
	return nil
}