package http

import (
	"database/sql"
	"errors"
	"log"
	"theSone/internal/models"

	"github.com/gofiber/fiber/v2"
)

// TODO: в принципе мессадж нужен только если фронтенд просит что-то отправлять, чтобы самим что-то отрисовать от твоего сообщения
// можно убрать, а так не принципиально. Статус коды уже почти всю инфу говорят
func (h *TodoHandler) errorResponse(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(fiber.Map{
			"error": fiberErr.Message,
		})
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

	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

func logInternalError(err error) {

	log.Printf("Internal error: %+v", err)
}
