package main

import (
	"fmt"
	"log/slog"
	"theSone/internal/application"
	delivery "theSone/internal/delivery/http"
	"theSone/internal/repository/postgres"
	"theSone/pkg"
	logger "theSone/pkg/logger"
)

func main() {
	cfg, err := pkg.ReadConfig()
	if err != nil {
		slog.Error("error loading config", "error", err.Error())
		return
	}
	log := logger.NewLogger(&cfg.Logger)
	// TODO: я же говорил вынести подключение на уровень репозитория
	db, err := postgres.ConnectDB(&cfg.Repo)
	if err != nil {
		slog.Error("error connecting to db", "error", err.Error())
		return
	}
	defer db.Close()

	repo := postgres.NewTodoRepository(db, log)
	app := application.NewApplication(repo, log)
	handler := delivery.NewTodoHandler(app)

	server := delivery.NewServer(handler)

	fmt.Printf("Server is running on %s\n", cfg.App.ServerPort)
	if err := server.Start(cfg.App.ServerPort); err != nil {
		slog.Error("server error", "error", err.Error())
	}
}
