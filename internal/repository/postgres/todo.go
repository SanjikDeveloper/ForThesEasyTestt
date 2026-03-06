package postgres

import (
	"context"
	"database/sql"
	"theSone/internal/models"
)

const (
	createTodoQuery  = `INSERT INTO todos (todo_list, description, created_at) VALUES ($1, $2, $3) RETURNING id_list`
	getTodoByIDQuery = `SELECT id_list, todo_list, description, created_at FROM todos WHERE id_list = $1`
	updateTodoQuery  = `UPDATE todos 
              SET todo_list = COALESCE($1, todo_list), 
                  description = COALESCE($2, description), 
                  created_at = COALESCE($3, created_at) 
              WHERE id_list = $4`
	deleteTodoQuery = `DELETE FROM todos WHERE id_list = $1`
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(ctx context.Context, t *models.Todo) error {
	// TODO: У тебя при вызове функции каждый раз создается новая строка с одинаковым SQL. Лучше вынеси это в константы
	return r.db.QueryRowContext(ctx, createTodoQuery, t.List, t.Description, t.CreatedAt).Scan(&t.ID)
}

func (r *TodoRepository) GetByID(ctx context.Context, id int) (*models.Todo, error) {
	var t models.Todo
	err := r.db.QueryRowContext(ctx, getTodoByIDQuery, id).Scan(&t.ID, &t.List, &t.Description, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TodoRepository) Update(ctx context.Context, t *models.Todo) error {
	// Мы используем COALESCE($1, todo_list).
	// Это значит: "Возьми новое значение, но если оно пустое ($1 is NULL), оставь старое".
	// TODO: в такой конструкции легко запнуться и обновить данные на пустые значения. Ты перекладываешь ответственность на базу
	// Лучшим и более явным решением будет использовать конструкцию, в которой ты сам проверяешь, если поле не пустое, то ты его обновляешь
	// Напиши такое решение используя библиотеку squirrel
	// Под апдейт ты можешь создать новую модельку уже с указателями и там проверять на наличие переменных
	// Напиши такое решение используя библиотеку squirrel
	// Под апдейт ты можешь создать новую модельку уже с указателями и там проверять на наличие переменных
	res, err := r.db.ExecContext(ctx, updateTodoQuery, t.List, t.Description, t.CreatedAt, t.ID)
	if err != nil {
		return err
	}
	// TODO: почему не обрабатываешь ошибку?
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TodoRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, deleteTodoQuery, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
