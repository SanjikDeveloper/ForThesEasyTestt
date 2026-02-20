package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	handlers "theSone/internal/deleviry/http"

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

	handlers.SetDB(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateTodo(w, r)
		case http.MethodGet:
			handlers.GetTodoById(w, r)
		case http.MethodPut:
			handlers.UpdateTodo(w, r)
		case http.MethodDelete:
			handlers.DeleteTodo(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("listening :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
