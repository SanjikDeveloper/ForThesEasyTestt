package http

import (
	"context"
	"errors"
	"strconv"
	"theSone/internal/application"
	"theSone/internal/models"
	"theSone/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app  *fiber.App
	addr string
}

func NewServer(app *application.Application, addr string) *Server {
	handler := NewTodoHandler(app)
	return &Server{
		app:  RegisterRoutes(handler),
		addr: addr,
	}
}

func (s *Server) Init() error {
	return nil
}

func (s *Server) Run(ctx context.Context) error {
	return s.app.Listen(s.addr)
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}

type TodoHandler struct {
	logger *logger.Logger
	//TODO: где интерфейс? Ты должен был интерфейс добавить сюда
	app *application.Application
}

func NewTodoHandler(app *application.Application) *TodoHandler {
	return &TodoHandler{app: app}
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
	idStr := c.Params("id")
	if idStr == "" {
		idStr = c.Query("id")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return h.errorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid or missing ID"))
	}

	var todo models.Todo
	if err = c.BodyParser(&todo); err != nil {
		return h.errorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}
	todo.ID = id

	if err = h.validateTodo(todo); err != nil {
		return h.errorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if err = h.app.UpdateTodo(c.Context(), &todo); err != nil {
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

func RegisterRoutes(handler *TodoHandler) *fiber.App {
	app := fiber.New()

	api := app.Group("/todo")
	api.Post("/", handler.createTodo)
	api.Get("/", handler.getAllTodos)
	api.Get("/:id", handler.getTodoById)
	api.Put("/:id", handler.updateTodo)
	api.Delete("/:id", handler.deleteTodo)

	return app
}
