package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort     string `env:"HTTP_PORT" env-default:"8080"`
	JWTSecret    string `env:"JWT_SECRET" env-required:"true"`
	PostgresUser string `env:"POSTGRES_USER" env-required:"true"`
	PostgresPass string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresDB   string `env:"POSTGRES_DB" env-required:"true"`
	PostgresHost string `env:"POSTGRES_HOST" env-required:"true"`
	PostgresPort string `env:"POSTGRES_PORT" env-required:"true"`
}

func New() (*Config, error) {
	cfg := &Config{}

	_ = godotenv.Load()

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
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
