package repository

import (
	"context"
	"database/sql"
)

type PostgresConf struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type repositoryDB struct {
	db  *sql.DB
	cfg *PostgresConf
}

func NewRepository(cfg *PostgresConf) *repositoryDB {
	return &repositoryDB{
		cfg: *PostgresConf{},
	}
}

type TodoList struct {
}

type TodoDO interface {
	Create(ctx context.Context, todoDO *TodoDO) error
	GetTodoById(ctx context.Context, todoDO *TodoDO) error
	Update(ctx context.Context, todoDO *TodoDO) error
	Delete(ctx context.Context, todoDO *TodoDO) error
}
