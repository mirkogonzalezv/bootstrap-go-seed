package routes

import (
	"microservice/cmd/container"
	"microservice/pkg/logger"

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
	ConfigureModule(router, r.container)
	logger.Success("Ruta API Registradas correctamente")
}
