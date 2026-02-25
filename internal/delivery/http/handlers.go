package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"theSone/internal/models"
)

type TodoRepository interface {
	Create(ctx context.Context, t *models.Todo) error
	GetByID(ctx context.Context, id int) (*models.Todo, error)
	Update(ctx context.Context, t *models.Todo) error
	Delete(ctx context.Context, id int) error
}

type TodoHandler struct {
	repo TodoRepository
}

func NewTodoHandler(repo TodoRepository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

func (h *TodoHandler) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return
	}
}

func (h *TodoHandler) validateTodo(t *models.Todo) string {
	if t.TodoList != nil && len(*t.TodoList) > 100 {
		return "Title should be more than 100"
	}
	if t.Description != nil && len(*t.Description) > 500 {
		return "Description should be less than 500"
	}
	return ""
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var t models.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if msg := h.validateTodo(&t); msg != "" {
		ErrorResponse(w, http.StatusBadRequest, msg)
		return
	}

	if err := h.repo.Create(r.Context(), &t); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "creating todo error")
		return
	}

	h.writeJSON(w, http.StatusCreated, t)
}

func (h *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid or missing ID")
		return
	}

	t, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "theres no todo like this")
		return
	}

	h.writeJSON(w, http.StatusOK, t)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid or missing ID")
		return
	}

	var t models.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	t.IdList = id

	if msg := h.validateTodo(&t); msg != "" {
		ErrorResponse(w, http.StatusBadRequest, msg)
		return
	}

	if err := h.repo.Update(r.Context(), &t); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, t)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid or missing ID")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
