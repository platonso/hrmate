package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type PostgresConfig struct {
	Host         string `env:"POSTGRES_HOST" env-required:"true"`
	Port         string `env:"POSTGRES_PORT" env-required:"true"`
	Username     string `env:"POSTGRES_USER" env-required:"true"`
	Password     string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Database     string `env:"POSTGRES_DB" env-required:"true"`
	MigrationDir string `env:"MIGRATION_DIR" env-default:"./migrations"`
}

type Config struct {
	PostgresConfig
	HTTPPort      string `env:"HTTP_PORT" env-default:"8080"`
	JWTSecret     string `env:"JWT_SECRET" env-required:"true"`
	AdminEmail    string `env:"ADMIN_EMAIL" env-required:"true"`
	AdminPassword string `env:"ADMIN_PASSWORD" env-required:"true"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func NewDB() (*PostgresConfig, error) {
	var cfg PostgresConfig

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *PostgresConfig) GetConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}
