package config

import (
	"log/slog"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

func ReadConfig(cfg any) error {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found, using environment variables", "error", err.Error())
	}
	if err := env.Parse(cfg); err != nil {
		return err
	}
	return nil
}
