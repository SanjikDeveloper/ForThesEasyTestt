package pkg

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// TODO: сделай на каждый слой свой конфиг, а в мейнике чтобы этот конфиг собирался из всех и дефолт не юзай
type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
	ServerPort string `env:"SERVER_PORT" envDefault:":8080"`
}

func ReadConfig() (*Config, error) {
	// TODO: и что делать при обработке ошибки?
	if err := godotenv.Load(); err != nil {
	}
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
