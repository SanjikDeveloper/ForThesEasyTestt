package http

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(handler *TodoHandler) *fiber.App {
	app := fiber.New()

	api := app.Group("/todo")
	api.Post("/", handler.createTodo)
	api.Get("/", handler.getAllTodos)
	api.Get("/:id", handler.getTodoById)
	api.Put("/", handler.updateTodo)
	api.Delete("/:id", handler.deleteTodo)

	return app
}
