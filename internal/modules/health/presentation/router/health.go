package router

import (
	"microservice/internal/modules/health/presentation/handler"
	"microservice/pkg/version"

	"github.com/gin-gonic/gin"
)

func ConfigHealthVersion(healthHandler *handler.HealthHandler) {
	wrapperFunc := func(rg *gin.RouterGroup, hdl any) {
		healthHandler := hdl.(*handler.HealthHandler)
		RegisterHealthRoutes(rg, healthHandler)
	}

	version.ConfigControllerVersion("health", healthHandler, wrapperFunc)
}

func RegisterHealthRoutes(rg *gin.RouterGroup, healthHandler *handler.HealthHandler) {
	rg.GET("", healthHandler.GetHealth)
}
