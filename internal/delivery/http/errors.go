package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// TODO: ты так и не сделал нормальный враппинг ошибок, должно быть что-то по типу
func errorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

var (
	exampleErrorNotFound       = errors.New("not found")
	exampleErrorInvalidRequest = errors.New("invalid request")
)

// Пример
func (h *TodoHandler) exampleErrorResponse(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, exampleErrorNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, exampleErrorInvalidRequest):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}
