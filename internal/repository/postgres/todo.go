package postgres

import (
	"context"
	"database/sql"
	"theSone/internal/models"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(ctx context.Context, t *models.Todo) error {
	query := `INSERT INTO todos (todo_list, description, created_at) VALUES ($1, $2, $3) RETURNING id_list`
	return r.db.QueryRowContext(ctx, query, t.TodoList, t.Description, t.CreatedAt).Scan(&t.IdList)
}

func (r *TodoRepository) GetByID(ctx context.Context, id int) (*models.Todo, error) {
	var t models.Todo
	query := `SELECT id_list, todo_list, description, created_at FROM todos WHERE id_list = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&t.IdList, &t.TodoList, &t.Description, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TodoRepository) Update(ctx context.Context, t *models.Todo) error {
	// Мы используем COALESCE($1, todo_list).
	// Это значит: "Возьми новое значение, но если оно пустое ($1 is NULL), оставь старое".
	query := `UPDATE todos 
              SET todo_list = COALESCE($1, todo_list), 
                  description = COALESCE($2, description), 
                  created_at = COALESCE($3, created_at) 
              WHERE id_list = $4`
	res, err := r.db.ExecContext(ctx, query, t.TodoList, t.Description, t.CreatedAt, t.IdList)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TodoRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM todos WHERE id_list = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
