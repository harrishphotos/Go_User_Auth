package models

import (
	"time"
)

type Store struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Location    string    `gorm:"not null" json:"location"`
	Books       []Book    `json:"books,omitempty"` // One-to-many relationship with Book
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
