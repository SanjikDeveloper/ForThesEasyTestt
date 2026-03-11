package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// TODO: посмотри снова пример, ты тут принимаешь статус и отдаешь этот статус, выходит масло масленное, даже непонятно для чеео этот метод нужен
// Если ты просто повторяешь тот статус, который передал. Посмотри пример ниже, у тебя есть кастомные ошибки приложения, ты кастомную ошибку приложения
// Разворачиваешь и дальше решаешь, что отправить
func errorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
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
