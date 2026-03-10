package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(cfg *Config) (*Storage, error) {
	db, err := ConnectDB(cfg)
	if err != nil {
		return nil, err
	}

	return &Storage{DB: db}, nil
}

func ConnectDB(cfg *Config) (*sql.DB, error) {
	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging db: %w", err)
	}

	return db, nil
}
