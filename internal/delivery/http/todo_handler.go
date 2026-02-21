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
	json.NewEncoder(w).Encode(v)
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var t models.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &t); err != nil {
		http.Error(w, "creating todo error", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, http.StatusCreated, t)
}

func (h *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}

	t, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, http.StatusOK, t)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}

	var t models.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	t.IdList = id

	if err := h.repo.Update(r.Context(), &t); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, http.StatusOK, t)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
