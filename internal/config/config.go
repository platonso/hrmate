package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPPort       string `env:"HTTP_PORT" env-default:"true"`
	JWTSecret      string `env:"JWT_SECRET" env-required:"true"`
	PostgresUser   string `env:"POSTGRES_USER" env-required:"true"`
	PostgresPass   string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresDB     string `env:"POSTGRES_DB" env-required:"true"`
	PostgresHost   string `env:"POSTGRES_HOST" env-required:"true"`
	PostgresPort   string `env:"POSTGRES_PORT" env-required:"true"`
	MigrationsPath string `env:"GOOSE_MIGRATION_DIR" env-default:"./migrations"`
}

func New() (*Config, error) {
	var cfg Config

	//if err := godotenv.Load(); err != nil {
	//	log.Println(".env file not found, using system env")
	//}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) GetConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		c.PostgresUser,
		c.PostgresPass,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
	)
}
