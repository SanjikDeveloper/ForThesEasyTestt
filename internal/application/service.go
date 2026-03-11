package application

import (
	"context"
	"theSone/internal/models"
	"theSone/pkg/logger"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type todoRepository interface {
	Create(ctx context.Context, t *models.Todo) error
	GetByID(ctx context.Context, id int) (*models.Todo, error)
	GetAll(ctx context.Context) ([]*models.Todo, error)
	Update(ctx context.Context, t *models.Todo) error
	Delete(ctx context.Context, id int) error
}

type Application struct {
	repo   todoRepository
	logger Logger
}

func NewApplication(repo todoRepository, logger *logger.Logger) *Application {
	return &Application{repo: repo, logger: logger}
}

func (a *Application) Init() error {
	return nil
}

func (a *Application) Run(ctx context.Context) error {
	return nil
}

func (a *Application) Stop() error {
	return nil
}

func (a *Application) CreateTodo(ctx context.Context, todo *models.Todo) error {
	return a.repo.Create(ctx, todo)
}

func (a *Application) GetTodoByID(ctx context.Context, id int) (*models.Todo, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *Application) GetAllTodo(ctx context.Context) ([]*models.Todo, error) {
	return a.repo.GetAll(ctx)
}

func (a *Application) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	return a.repo.Update(ctx, todo)
}

func (a *Application) DeleteTodo(ctx context.Context, id int) error {
	return a.repo.Delete(ctx, id)
}
