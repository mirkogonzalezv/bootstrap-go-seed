package middleware

import (
	"microservice/pkg/apperrors"
	"microservice/pkg/httputil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandlerMiddleware maneja errores globalmente
func ErrorHandlerMiddleware(log *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		var err error

		switch e := recovered.(type) {
		case *apperrors.DomainError:
			err = e
		case error:
			err = e
		default:
			err = apperrors.NewInternalError("PANIC_001", "Unexpected panic occurred")
		}

		// Mapear error a HTTP
		statusCode, errorResponse := httputil.MapErrorToHttp(err)

		// Log según severidad
		if statusCode >= 500 {
			log.Error("Server error",
				zap.String("error", err.Error()),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Int("status", statusCode),
			)
		} else {
			log.Warn("Client error",
				zap.String("error", err.Error()),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Int("status", statusCode),
			)
		}

		c.JSON(statusCode, errorResponse)
	})
}
