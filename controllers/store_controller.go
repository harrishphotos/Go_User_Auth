package controllers

import (
	"github.com/gofiber/fiber/v2"
	"book_crud/database"
	"book_crud/models"
)

// Get all stores with their books
func GetStores(c *fiber.Ctx) error {
	var stores []models.Store
	database.DB.Preload("Books").Find(&stores)
	return c.JSON(stores)
}

// Get a specific store with its books
func GetStore(c *fiber.Ctx) error {
	id := c.Params("id")
	var store models.Store
	if err := database.DB.Preload("Books").First(&store, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Store not found",
		})
	}
	return c.JSON(store)
}

// Create a new store
func CreateStore(c *fiber.Ctx) error {
	var store models.Store
	if err := c.BodyParser(&store); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := database.DB.Create(&store).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create store",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(store)
}

// Update a store
func UpdateStore(c *fiber.Ctx) error {
	id := c.Params("id")
	var store models.Store
	if err := database.DB.First(&store, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Store not found",
		})
	}

	if err := c.BodyParser(&store); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := database.DB.Save(&store).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update store",
		})
	}

	return c.JSON(store)
}

// Delete a store
func DeleteStore(c *fiber.Ctx) error {
	id := c.Params("id")
	var store models.Store
	if err := database.DB.First(&store, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Store not found",
		})
	}

	// Delete associated books first
	if err := database.DB.Where("store_id = ?", id).Delete(&models.Book{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete associated books",
		})
	}

	if err := database.DB.Delete(&store).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete store",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Get all books in a store
func GetStoreBooks(c *fiber.Ctx) error {
	storeID := c.Params("id")
	var books []models.Book
	if err := database.DB.Where("store_id = ?", storeID).Find(&books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch books",
		})
	}
	return c.JSON(books)
} 