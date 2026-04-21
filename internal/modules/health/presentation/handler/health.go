package handler

import (
	errResponse "microservice/internal/modules/health/application/dtos/error"
	usecases "microservice/internal/modules/health/application/use_cases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	uc      *usecases.HealthUseCase
	Version int
}

// Constructor del Controller
func NewHealthController(uc *usecases.HealthUseCase) *HealthHandler {
	return &HealthHandler{
		uc:      uc,
		Version: 1, // Asignamos la versión del módulo
	}
}

// Acceder a la versión siempre debe ir!
func (h *HealthHandler) GetVersion() int {
	return h.Version
}

// @Summary Health Check
// @Description Returns the health status of the service
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} responses.HealthResponse
// @Failure 503 {object} responses.ErrorResponse
// @Router /health [get]
func (h *HealthHandler) GetHealth(c *gin.Context) {
	res, err := h.uc.Execute()

	if err != nil {
		errorResponse := errResponse.NewErrorResponse("Service unhealthy")
		c.JSON(http.StatusServiceUnavailable, errorResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}
