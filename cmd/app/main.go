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

// TODO: почему gitignore пустой?
// TODO: где папка application?
// TODO: я просил тебя не использовать ИИ для написания чего-либо. Перепиши taskfile своими руками хотя бы для команды go run
func main() {
	// TODO: переделай на godotenv
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:1234@localhost:5432/todo_db?sslmode=disable"
	}

	// TODO: сделай лучше подключение к бд отдельно именно в бд, логику грейсфуллшатдауна тоже там же
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		// TODO: замени фатал на лог и return
		log.Fatal(err)
	}
	defer db.Close()

	// TODO: можно переиспользовать err
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewTodoRepository(db)
	handler := delivery.NewTodoHandler(repo)

	// TODO: вынеси роутер отдельно в папку http, так логически будет правильнее
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
	// TODO: порт выведи в конфиг и не используй фаталы вообще, лучше обработай ошибку и выйди с функции через return
	log.Fatal(http.ListenAndServe(":8080", mux))
}
