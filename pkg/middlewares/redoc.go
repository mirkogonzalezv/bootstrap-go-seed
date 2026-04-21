package middlewares

import (
	domainErrors "microservice/pkg/errors"
	httpMapper "microservice/pkg/http"
	"microservice/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandlerMiddleware maneja errores globalmente
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		var err error

		switch e := recovered.(type) {
		case *domainErrors.DomainError:
			err = e
		case error:
			err = e
		default:
			err = domainErrors.NewInternalError("PANIC_001", "Unexpected panic occurred")
		}

		// Mapear error a HTTP
		statusCode, errorResponse := httpMapper.MapErrorToHttp(err)

		// Log según severidad
		if statusCode >= 500 {
			logger.L().Error("Server error",
				zap.String("error", err.Error()),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Int("status", statusCode),
			)
		} else {
			logger.L().Warn("Client error",
				zap.String("error", err.Error()),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Int("status", statusCode),
			)
		}

		c.JSON(statusCode, errorResponse)
	})
}
