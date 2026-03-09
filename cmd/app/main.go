package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"theSone/internal/application"
	delivery "theSone/internal/delivery/http"
	"theSone/internal/repository/postgres"
	"theSone/pkg"
	logger "theSone/pkg/logger"

	// TODO: поменя на pgx/v5
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := pkg.ReadConfig()
	if err != nil {
		slog.Error("error loading config", "error", err.Error())
		return
	}
	log := logger.NewLogger(&cfg.Logger)
	//TODO: вынеси в репозиторий
	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Repo.DBHost, cfg.Repo.DBPort, cfg.Repo.DBUser, cfg.Repo.DBPassword, cfg.Repo.DBName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("error opening db", "error", err.Error())
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("error pinging db", "error", err.Error())
		return
	}

	repo := postgres.NewTodoRepository(db)
	app := application.NewApplication(repo, log)
	handler := delivery.NewTodoHandler(app)

	mux := delivery.RegisterRoutes(handler)

	fmt.Printf("Server is running on %s\n", cfg.App.ServerPort)
	if err := http.ListenAndServe(cfg.App.ServerPort, mux); err != nil {
		slog.Error("server error", "error", err.Error())
	}
}
