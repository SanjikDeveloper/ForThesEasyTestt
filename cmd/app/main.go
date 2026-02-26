package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	delivery "theSone/internal/delivery/http"
	"theSone/internal/repository/postgres"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:1234@localhost:5432/todo_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
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
			delivery.ErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
