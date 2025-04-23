package controllers

import (
	"book_crud/database"
	"book_crud/models"

	"github.com/gofiber/fiber/v2"
)

func GetBooks(c *fiber.Ctx) error{

	var books []models.Book
	database.DB.Find(&books)
	return c.JSON(books)
}

func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book

	if err := database.DB.Find(&book, id).Error; err!=nil{

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":"book not found",
		})
	}

	return c.JSON(book)
}

func CreateBook(c *fiber.Ctx) error {
	var book models.Book

	if err := c.BodyParser(&book); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{

			"error":"cant prase json",
		})


	}

	if err :=  database.DB.Save(&book).Error; err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"failed to create a book",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(book)

}

func UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book
	if err := database.DB.First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := database.DB.Save(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update book",
		})
	}

	return c.JSON(book)
}

func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book
	if err := database.DB.First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	if err := database.DB.Delete(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete book",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
} 