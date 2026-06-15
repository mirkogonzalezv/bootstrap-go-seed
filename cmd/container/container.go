package container

import (
	"context"
	"fmt"
	"microservice/cmd/config"
	healthUseCase "microservice/internal/health/application/usecase"
	healthInfra "microservice/internal/health/infrastructure"
	"microservice/internal/health/presentation/handler"
	"microservice/pkg/database"
	"microservice/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Container struct {
	DB            *pgxpool.Pool
	HealthHandler *handler.HealthHandler
	Log           *zap.Logger
}

func NewContainer(cfg *config.Configuration, log *zap.Logger) (*Container, error) {

	// Configuración de la DB
	dbCfg := database.PostgresConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}

	pool, err := database.NewPool(context.Background(), dbCfg)
	if err != nil {
		return nil, fmt.Errorf("container: %w", err)
	}

	repo := healthInfra.NewPostgresHealthRepository(pool)
	healthUC := healthUseCase.NewHealthUseCase(repo)
	healthHdlr := handler.NewHealthHandler(healthUC)

	logger.Success("Container de dependencias inicializado...")

	return &Container{
		DB:            pool,
		HealthHandler: healthHdlr,
		Log:           log,
	}, nil
}
