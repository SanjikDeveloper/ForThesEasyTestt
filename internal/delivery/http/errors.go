package http

import (
	"database/sql"
	"errors"
	"theSone/internal/models"

	"github.com/gofiber/fiber/v2"
)

func (h *TodoHandler) errorResponse(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		h.logger.Error("Fiber error", "code", fiberErr.Code, "err", err)
		return c.SendStatus(fiberErr.Code)
	}

	var status int
	var message string

	switch {
	case errors.Is(err, models.ErrNotFound), errors.Is(err, sql.ErrNoRows):
		status = fiber.StatusNotFound
	case errors.Is(err, models.ErrInvalidInput):
		status = fiber.StatusBadRequest
		message = "Invalid input data"
	case errors.Is(err, models.ErrUnauthorized):
		status = fiber.StatusUnauthorized
	default:
		status = fiber.StatusInternalServerError
		h.logger.Error("Internal server error")
	}

	if message == "" {
		return c.SendStatus(status)
	}

	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}
