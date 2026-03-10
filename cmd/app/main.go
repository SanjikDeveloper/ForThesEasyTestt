package main

import (
	"context"
	"log/slog"
	"theSone/internal/application"
	delivery "theSone/internal/delivery/http"
	"theSone/internal/repository/postgres"
	"theSone/pkg/config"
	logger "theSone/pkg/logger"
	service "theSone/pkg/services"
)

type Config struct {
	Repo   postgres.Config    `envPrefix:"REPO_"`
	Logger logger.Config      `envPrefix:"LOGGER_"`
	App    delivery.AppConfig `envPrefix:"APP_"`
}

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		slog.Error("error loading config", "error", err.Error())
		return
	}
	log := logger.NewLogger(&cfg.Logger)

	repos := postgres.NewTodoRepository(&cfg.Repo, log)
	app := application.NewApplication(repos, log)
	server := delivery.NewServer(app, cfg.App.ServerPort)

	manager := service.NewManager(log)
	manager.AddService(repos, app, server)

	log.Info("Starting todoservice")
	if err := manager.Run(context.Background()); err != nil {
		log.Error("failed to start", "error", err.Error())
	}
}
