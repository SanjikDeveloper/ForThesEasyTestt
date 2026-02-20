package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"theSone/internal/models"
)

var DB *sql.DB

func SetDB(db *sql.DB) {
	DB = db
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var t models.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO todos (todo_list, description, created_at) VALUES ($1, $2, $3) RETURNING id_list`
	err := DB.QueryRow(query, t.TodoList, t.Description, t.CreatedAt).Scan(&t.IdList)
	if err != nil {
		http.Error(w, "creating todo error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

func GetTodoById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}

	var t models.Todo
	query := `SELECT id_list, todo_list, description, created_at FROM todos WHERE id_list = $1`
	err = DB.QueryRow(query, id).Scan(&t.IdList, &t.TodoList, &t.Description, &t.CreatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, t)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}

	var t models.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `UPDATE todos SET todo_list = $1, description = $2, created_at = $3 WHERE id_list = $4`
	res, err := DB.Exec(query, t.TodoList, t.Description, t.CreatedAt, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	t.IdList = id
	writeJSON(w, http.StatusOK, t)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM todos WHERE id_list = $1`
	res, err := DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
