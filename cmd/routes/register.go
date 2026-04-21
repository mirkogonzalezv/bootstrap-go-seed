package routes

import (
	"microservice/cmd/container"

	health "microservice/internal/modules/health/presentation/router"
)

func ConfigureModule(c *container.Container) {
	health.ConfigHealthVersion(c.HealthHandler)
	// Registramos las rutas de forma dinámica con su versión desde aquí
}
