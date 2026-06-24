package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
)

type Configuration struct {
	AppEnv          string `env:"ENVIRONMENT" envDefault:"dev" validate:"oneof= dev qa prod"`
	Port            int    `env:"PORT" envDefault:"3200"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`
	ShutDownTimeOut int    `env:"SHUTDOWN_TIMEOUT" envDefault:"30"`
	// Se definen las variables de entorno necesarias para el proyecto
	// DB, Redis, API_PATH, etc...
	DBHost          string `env:"DB_HOST" envDefault:"localhost"`
	DBPort          int    `env:"DB_PORT" envDefault:"5432"`
	DBUser          string `env:"DB_USER" envDefault:"postgres"`
	DBPassword      string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName          string `env:"DB_NAME" envDefault:"microservice"`
	DBSSLMode       string `env:"DB_SSL_MODE" envDefault:"disable"`
	PubSubProjectID string `env:"PUBSUB_PROJECT_ID" envDefault:""`
}

func LoadVars() (*Configuration, error) {
	cfg := &Configuration{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	v := validator.New()

	if err := v.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
