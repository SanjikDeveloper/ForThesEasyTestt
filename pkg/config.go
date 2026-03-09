package pkg

import (
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
	// TODO: и что делать при обработке ошибки?
	if err := godotenv.Load(); err != nil {
		// Log error if needed or ignore if .env is optional
	}
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
