package router

import (
	"microservice/internal/health/presentation/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.HealthHandler) {
	v1 := api.Group("/v1/health")

	v1.GET("", h.Health)
}
