package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"theSone/internal/application"
	"theSone/internal/models"
	"theSone/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
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

func (h *TodoHandler) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return
	}
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

// TODO: все ручки тоже переделай на fiber...
func (h *TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validateTodo(todo); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.app.CreateTodo(r.Context(), &todo); err != nil {
		errorResponse(w, http.StatusInternalServerError, "error creating todo")
		return
	}

	h.writeJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) getAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.app.GetAllTodo(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "error fetching todos")
		return
	}

	h.writeJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) getTodoById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid or missing ID")
		return
	}

	todo, err := h.app.GetTodoByID(r.Context(), id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "todo not found")
		return
	}

	h.writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid or missing ID")
		return
	}

	var todo models.Todo
	// TODO: я же говорил везде перепроверить на то, что ты переиспользуешь переменные

	if err = json.NewDecoder(r.Body).Decode(&todo); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	todo.ID = id

	// тут тоже просто err =
	if err = h.validateTodo(todo); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.app.UpdateTodo(r.Context(), &todo); err != nil {
		errorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid or missing ID")
		return
	}

	// тут тоже err =
	// ты сам должен все это проверить
	if err := h.app.DeleteTodo(r.Context(), id); err != nil {
		errorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", http.StatusText(http.StatusOK))
}

func RegisterRoutes(handler *TodoHandler) *fiber.App {
	app := fiber.New()
	// TODO: ручки во множественном числе не надо называть
	app.Post("/todos", adaptor.HTTPHandlerFunc(handler.createTodo))
	app.Get("/todos", func(c *fiber.Ctx) error {
		// TODO: просто сделай отдельную ручку на /todo/:id
		if c.Query("id") != "" {
			return adaptor.HTTPHandlerFunc(handler.getTodoById)(c)
		}
		return adaptor.HTTPHandlerFunc(handler.getAllTodos)(c)
	})
	app.Put("/todos", adaptor.HTTPHandlerFunc(handler.updateTodo))
	app.Delete("/todos", adaptor.HTTPHandlerFunc(handler.deleteTodo))

	return app
}
