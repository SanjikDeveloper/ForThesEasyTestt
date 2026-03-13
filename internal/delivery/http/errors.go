package http

import (
	"database/sql"
	"errors"
	"log"
	"theSone/internal/models"

	"github.com/gofiber/fiber/v2"
)

func (h *TodoHandler) errorResponse(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(fiber.Map{})
	}

	var status int
	var message string

	switch {
	case errors.Is(err, models.ErrNotFound), errors.Is(err, sql.ErrNoRows):
		status = fiber.StatusNotFound
		message = "Todo not found"
	case errors.Is(err, models.ErrInvalidInput):
		status = fiber.StatusBadRequest
		message = "Invalid input data"
	case errors.Is(err, models.ErrUnauthorized):
		status = fiber.StatusUnauthorized
		message = "Unauthorized"
	default:
		status = fiber.StatusInternalServerError
		message = "Internal server error"

		if h.app != nil {
			logInternalError(err)
		}
	}
	// TODO: пользователям не надо видеть что за ошибка на беке. Ее главное залогировать и отдать 500 статус
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

// TODO: сделай данную функуцию методом handler и у нее используй общий логгер для сервиса. Везде где нужно что-то логировать
// прокидывай общий логгер и его используй
func logInternalError(err error) {

	log.Printf("Internal error: %+v", err)
}
