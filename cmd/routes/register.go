package routes

import (
	"microservice/cmd/container"

	health "microservice/internal/health/presentation/router"

	"github.com/gin-gonic/gin"
)

func ConfigureModule(router *gin.Engine, c *container.Container) {
	api := router.Group("api")
	health.RegisterRoutes(api, c.HealthHandler)
	// Registramos las rutas de forma dinámica con su versión desde aquí
}
