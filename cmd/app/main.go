package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	delivery "theSone/internal/delivery/http"
	"theSone/internal/repository/postgres"
	"theSone/pkg"

	_ "github.com/lib/pq"
)

// TODO: где папка application?
func main() {
	cfg, err := pkg.ReadConfig()
	if err != nil {
		slog.Error("error loading config: %s", err.Error())
		return
	}

	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("error opening db: %s", err.Error())
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("error pinging db: %s", err.Error())
		return
	}

	repo := postgres.NewTodoRepository(db)
	handler := delivery.NewTodoHandler(repo)

	mux := delivery.RegisterRoutes(handler)

	fmt.Printf("Server is running on %s\n", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, mux); err != nil {
		slog.Error("server error: %s", err.Error())
	}
}
