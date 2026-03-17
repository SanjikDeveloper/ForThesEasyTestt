package http

import (
	"context"
	"errors"
	"strconv"
	"theSone/internal/models"
	"theSone/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type TodoService interface {
	CreateTodo(ctx context.Context, todo *models.Todo) error
	GetTodoByID(ctx context.Context, id int) (*models.Todo, error)
	GetAllTodo(ctx context.Context) ([]*models.Todo, error)
	UpdateTodo(ctx context.Context, todo *models.Todo) error
	DeleteTodo(ctx context.Context, id int) error
}

type TodoHandler struct {
	logger logger.Logger
	app    TodoService
}

func NewTodoHandler(app TodoService, log logger.Logger) *TodoHandler {
	return &TodoHandler{app: app, logger: log}
}

func (h *TodoHandler) validateTodo(todo models.Todo) error {
	if len(todo.List) > 100 {
		return errors.New("title should be less than 100 characters")
	}
	if len(todo.Description) > 500 {
		return errors.New("description should be less than 500 characters")
	}
	return nil
}

func (h *TodoHandler) createTodo(c *fiber.Ctx) error {
	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		return h.errorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid request body"))
	}

	if err := h.validateTodo(todo); err != nil {
		return h.errorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if err := h.app.CreateTodo(c.Context(), &todo); err != nil {
		return h.errorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(todo)
}

func (h *TodoHandler) getAllTodos(c *fiber.Ctx) error {
	todos, err := h.app.GetAllTodo(c.Context())
	if err != nil {
		return h.errorResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(todos)
}

func (h *TodoHandler) getTodoById(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		idStr = c.Query("id")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return h.errorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid or missing ID"))
	}

	todo, err := h.app.GetTodoByID(c.Context(), id)
	if err != nil {
		return h.errorResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(todo)
}

func (h *TodoHandler) updateTodo(c *fiber.Ctx) error {
	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		return h.errorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if todo.ID == 0 {
		return h.errorResponse(c, fiber.NewError(fiber.StatusBadRequest, "missing ID in request body"))
	}

	if err := h.validateTodo(todo); err != nil {
		return h.errorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if err := h.app.UpdateTodo(c.Context(), &todo); err != nil {
		return h.errorResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(todo)
}

func (h *TodoHandler) deleteTodo(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		idStr = c.Query("id")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return h.errorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid or missing ID"))
	}

	if err := h.app.DeleteTodo(c.Context(), id); err != nil {
		return h.errorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
