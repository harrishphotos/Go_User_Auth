package routes

import (
	"github.com/gofiber/fiber/v2"
	"book_crud/controllers"
)

func SetupStoreRoutes(api fiber.Router) {
	stores := api.Group("/stores")

	// Store CRUD routes
	stores.Get("/", controllers.GetStores)
	stores.Get("/:id", controllers.GetStore)
	stores.Post("/", controllers.CreateStore)
	stores.Put("/:id", controllers.UpdateStore)
	stores.Delete("/:id", controllers.DeleteStore)

	// Store-Book relationship routes
	stores.Get("/:id/books", controllers.GetStoreBooks)
} 