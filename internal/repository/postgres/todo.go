package postgres

import (
	"context"
	"database/sql"
	"theSone/internal/models"
	"theSone/pkg/logger"

	"github.com/Masterminds/squirrel"
)

const (
	createTodoQuery  = `INSERT INTO todos (todo_list, description, created_at) VALUES ($1, $2, $3) RETURNING id_list`
	getTodoByIDQuery = `SELECT id_list, todo_list, description, created_at FROM todos WHERE id_list = $1`
	updateTodoQuery  = `UPDATE todos 
              SET todo_list = COALESCE($1, todo_list), 
                  description = COALESCE($2, description), 
                  created_at = COALESCE($3, created_at) 
              WHERE id_list = $4`
	deleteTodoQuery  = `DELETE FROM todos WHERE id_list = $1`
	getAllTodosQuery = `SELECT id_list, todo_list, description, created_at FROM todos ORDER BY created_at DESC`
)

type TodoRepository struct {
	db     *sql.DB
	cfg    *Config
	logger logger.Logger
}

func NewTodoRepository(cfg *Config, log logger.Logger) *TodoRepository {
	return &TodoRepository{cfg: cfg, logger: log}
}

func (r *TodoRepository) Init() error {
	db, err := ConnectDB(r.cfg)
	if err != nil {
		return err
	}
	r.db = db
	return nil
}

func (r *TodoRepository) Run(ctx context.Context) error {
	return nil
}

func (r *TodoRepository) Stop() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

func (r *TodoRepository) Create(ctx context.Context, t *models.Todo) error {
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

func (r *TodoRepository) GetAll(ctx context.Context) ([]*models.Todo, error) {
	rows, err := r.db.QueryContext(ctx, getAllTodosQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		t := new(models.Todo)
		if err := rows.Scan(&t.ID, &t.List, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepository) Update(ctx context.Context, t *models.Todo) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	builder := psql.Update("todos").Where(squirrel.Eq{"id_list": t.ID})

	if t.List != "" {
		builder = builder.Set("todo_list", t.List)
	}
	if t.Description != "" {
		builder = builder.Set("description", t.Description)
	}
	if !t.CreatedAt.IsZero() {
		builder = builder.Set("created_at", t.CreatedAt)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
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
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
