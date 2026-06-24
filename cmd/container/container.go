package container

import (
	"microservice/cmd/config"
	healthUseCase "microservice/internal/health/application/usecase"
	healthInfra "microservice/internal/health/infrastructure"
	"microservice/internal/health/presentation/handler"
	"microservice/pkg/database"
	"microservice/pkg/logger"
	"microservice/pkg/pubsub"

	"go.uber.org/zap"
)

type Container struct {
	DB            database.Database
	PubSub        pubsub.PubSub
	HealthHandler *handler.HealthHandler
	Log           *zap.Logger
}

func NewContainer(cfg *config.Configuration, db database.Database, ps pubsub.PubSub, log *zap.Logger) (*Container, error) {

	repo := healthInfra.NewPostgresHealthRepository(db)
	healthUC := healthUseCase.NewHealthUseCase(repo)
	healthHdlr := handler.NewHealthHandler(healthUC)

	logger.Success("Container de dependencias inicializado...")

	return &Container{
		DB:            db,
		PubSub:        ps,
		HealthHandler: healthHdlr,
		Log:           log,
	}, nil
}
