package pkg

import (
	"log/slog"
	"theSone/internal/repository/postgres"
	"theSone/pkg/logger"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort string `env:"SERVER_PORT" envDefault:":8080"`
}

type Config struct {
	Repo   postgres.Config `envPrefix:"REPO_"`
	Logger logger.Config   `envPrefix:"LOGGER_"`
	App    AppConfig       `envPrefix:"APP_"`
}

func ReadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found, using environment variables", "error", err.Error())
	}
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
