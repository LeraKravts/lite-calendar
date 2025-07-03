package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppPort          string `env:"APP_PORT" envDefault:"8080"`
	AppEnv           string `env:"APP_ENV"  envDefault:"local"`
	PostgresHost     string `env:"POSTGRES_HOST"     envDefault:"localhost"`
	PostgresPort     int    `env:"POSTGRES_PORT"     envDefault:"5432"`
	PostgresUser     string `env:"POSTGRES_USER"     envDefault:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" envDefault:"password"`
	PostgresDB       string `env:"POSTGRES_DB"       envDefault:"calendar_db"`
}

func Load() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("❌Failed to load config: v", err)
	}
	return &cfg

}
