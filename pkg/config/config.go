package config

import (
	"log/slog"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

func ReadConfig() (*cmd.Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found, using environment variables", "error", err.Error())
	}
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
