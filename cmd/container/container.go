package container

import (
	"microservice/cmd/config"
	healthUseCase "microservice/internal/health/application/usecase"
	"microservice/internal/health/presentation/handler"
	"microservice/pkg/logger"

	"go.uber.org/zap"
)

type Container struct {
	HealthHandler *handler.HealthHandler
	Log           *zap.Logger
}

func NewContainer(cfg *config.Configuration, log *zap.Logger) *Container {
	healthUC := healthUseCase.NewHealthUseCase()
	healthHdlr := handler.NewHealthHandler(healthUC)

	logger.Success("Container de dependencias inicializado...")

	return &Container{
		HealthHandler: healthHdlr,
		Log:           log,
	}
}
