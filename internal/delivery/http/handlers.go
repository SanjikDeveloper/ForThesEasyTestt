package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"theSone/internal/application"
	"theSone/internal/models"
)

type TodoHandler struct {
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
	todos, err := h.app.GetAll(r.Context())
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
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	todo.ID = id

	if err := h.validateTodo(todo); err != nil {
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
	if err := h.app.DeleteTodo(r.Context(), id); err != nil {
		errorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "todo deleted successfully"}`)
}

func RegisterRoutes(handler *TodoHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.createTodo(w, r)
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				handler.getTodoById(w, r)
			} else {
				handler.getAllTodos(w, r)
			}
		case http.MethodPut:
			handler.updateTodo(w, r)
		case http.MethodDelete:
			handler.deleteTodo(w, r)
		default:
			errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	return mux
}
