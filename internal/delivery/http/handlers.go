package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"theSone/internal/models"
)

// TODO: лучше сделать не экспортируемым данный интерфейс, ты его используешь только в нынешнем пакете
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
	// TODO: почему ошибку не логируешь?
	if err != nil {
		return
	}
}

// TODO: почему функция валидации возвращает string? Не лучше ли возвращать ошибку? Она логически подходит под валидацию
// Если не прошел валидацию = ошибка
func (h *TodoHandler) validateTodo(t *models.Todo) string {
	if t.TodoList != nil && len(*t.TodoList) > 100 {
		// TODO: текста и ошибки, все пиши в нижнем регистре
		return "Title should be more than 100"
	}
	if t.Description != nil && len(*t.Description) > 500 {
		return "Description should be less than <500"
	}
	return ""
}

// TODO: почему методы экспортируемые?
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// TODO: называй переменные нормальными именами, если это один объект, то можно назвать "todo", если это слайс
	// То можешь во множественном числе назвать
	var t models.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if msg := h.validateTodo(&t); msg != "" {
		ErrorResponse(w, http.StatusBadRequest, msg)
		return
	}
	// TODO: если нет явной цели передавать модель по указателю, то лучше передавать без него, это может создать дополнительные сложноти
	// Это может привести к тому, что где-то из слоев выше кто-то может модифицировать структуру, а ты можешь об этом и не узнать
	// Если будешь работать в рамках данного пакета
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
	// TODO:
	t, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "theres no todo like this")
		return
	}

	h.writeJSON(w, http.StatusOK, t)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// TODO: в put методе ты можешь просто принять все тело запроса, не надо отдельно с URL доставать ID и отдельно
	// С тела запроса брать все остальные данные
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid or missing ID")
		return
	}
	// TODO: лучше будет сделать отдельную модельку под реквест, условно updateTodoRequest в которой будут переменные
	// которые ты принимаешь на стороне сервера, а потом уже из реквеста делать модельку для передачи данных выше
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
	// TODO: статус поменять надо на 200. Статус 204 говорит о том, что к тебе обратились за контентом, а у тебя его нет
	w.WriteHeader(http.StatusNoContent)
}
