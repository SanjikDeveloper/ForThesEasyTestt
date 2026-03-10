package http

import (
	"errors"
	"log"
	"net/http"
)

// TODO: посмотри снова пример, ты тут принимаешь статус и отдаешь этот статус, выходит масло масленное, даже непонятно для чеео этот метод нужен
// Если ты просто повторяешь тот статус, который передал. Посмотри пример ниже, у тебя есть кастомные ошибки приложения, ты кастомную ошибку приложения
// Разворачиваешь и дальше решаешь, что отправить
func errorResponse(w http.ResponseWriter, status int, message string) {
	switch status {
	case http.StatusOK:
		w.WriteHeader(http.StatusOK)
	case http.StatusBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	case http.StatusUnauthorized:
		w.WriteHeader(http.StatusUnauthorized)
	case http.StatusForbidden:
		w.WriteHeader(http.StatusForbidden)
	case http.StatusNotFound:
		w.WriteHeader(http.StatusNotFound)
	case http.StatusInternalServerError:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusInternalServerError)
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
