package routes

import (
	"github.com/gofiber/fiber/v2"
	"book_crud/controllers"
)

func SetupBookRoutes(api fiber.Router) {
	books := api.Group("/books")

	// Book CRUD routes
	books.Get("/", controllers.GetBooks)
	books.Get("/:id", controllers.GetBook)
	books.Post("/", controllers.CreateBook)
	books.Put("/:id", controllers.UpdateBook)
	books.Delete("/:id", controllers.DeleteBook)
} 