package usecases

import (
	healthRes "microservice/internal/modules/health/application/dtos/response"
)

type HealthUseCase struct{}

// Constructor del Caso de Uso
func NewHealthUseCase() *HealthUseCase {
	return &HealthUseCase{}
}

func (h *HealthUseCase) Execute() (*healthRes.HealthResponse, error) {

	// Aqui podriamos agregar otras validaciones:
	// - Verificar conexión a DB
	// - Verificación con servicios externos
	// Si alguno falla, se retorna error HTTP 503
	return healthRes.NewHealthResponse(), nil
}
