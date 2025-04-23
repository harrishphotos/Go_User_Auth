package models

import (
	"time"
)

type Book struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Author      string    `gorm:"not null" json:"author"`
	Description string    `json:"description"`
	StoreID     uint      `json:"store_id"` // Foreign key for Store
	Store       Store     `json:"store,omitempty"` // Belongs to relationship
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
