package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	delivery "theSone/internal/delivery/http"
	"theSone/internal/repository/postgres"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("")

	db, err := sql.Open("", dsn)
	if err != nil {
		log.Fatal("failed to Db connet", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Println("ping db", err)
	}

	repo := postgres.NewTodoRepository(db)
	handler := delivery.NewTodoHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateTodo(w, r)
		case http.MethodGet:
			handler.GetTodoById(w, r)
		case http.MethodPut:
			handler.UpdateTodo(w, r)
		case http.MethodDelete:
			handler.DeleteTodo(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("listening :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
