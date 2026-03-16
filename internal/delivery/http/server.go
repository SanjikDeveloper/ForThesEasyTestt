package http

import (
	"context"
	"theSone/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app  *fiber.App
	addr string
}

func NewServer(app TodoService, addr string, log logger.Logger) *Server {
	handler := NewTodoHandler(app, log)
	return &Server{
		app:  RegisterRoutes(handler),
		addr: addr,
	}
}

func (s *Server) Init() error {
	return nil
}

func (s *Server) Run(ctx context.Context) error {
	return s.app.Listen(s.addr)
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}
