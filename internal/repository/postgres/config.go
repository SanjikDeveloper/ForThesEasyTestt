package postgres

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
	// TODO: я писал чтобы ты не использовал дефолты
	ServerPort string `env:"SERVER_PORT" envDefault:":8080"`
}
