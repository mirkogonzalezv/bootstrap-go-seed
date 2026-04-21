package routes

import (
	"microservice/cmd/container"
	"microservice/pkg/logger"
	"microservice/pkg/version"

	"github.com/gin-gonic/gin"
)

type APIRouter struct {
	container *container.Container
}

func NewAPIRouter(cont *container.Container) *APIRouter {
	return &APIRouter{
		container: cont,
	}
}

func (r *APIRouter) RegisterRoutes(router *gin.Engine) {
	ConfigureModule(r.container)

	// Base Path: Ej -> "api"
	version.BuildRoutes(router, "api")

	logger.Success("Ruta API Registradas correctamente")
}
