package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
)

type Configuration struct {
	AppEnv          string `env:"ENVIRONMENT" envDefault:"dev" validate:"oneof= dev qa prod"`
	Port            int    `env:"PORT" envDefault:"4200"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`
	ShutDownTimeOut int    `env:"SHUTDOWN_TIMEOUT" envDefault:"30"`
	// Se definen las variables de entorno necesarias para el proyecto
	// DB, Redis, API_PATH, etc...
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
