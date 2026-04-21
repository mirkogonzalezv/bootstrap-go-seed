package container

import (
	"microservice/cmd/config"
	usecases "microservice/internal/modules/health/application/use_cases"
	"microservice/internal/modules/health/presentation/handler"
	"microservice/pkg/logger"
)

type Container struct {
	HealthHandler *handler.HealthHandler
}

func NewContainer(cfg *config.Configuration) *Container {
	healthUseCase := usecases.NewHealthUseCase()
	healthHdlr := handler.NewHealthController(healthUseCase)

	logger.Success("Container de dependencias inicializado...")

	return &Container{
		HealthHandler: healthHdlr,
	}
}
