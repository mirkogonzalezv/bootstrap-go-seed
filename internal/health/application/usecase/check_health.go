package usecase

import "microservice/internal/health/application/dto"

type HealthUseCase struct{}

// Constructor del Caso de Uso
func NewHealthUseCase() *HealthUseCase {
	return &HealthUseCase{}
}

func (h *HealthUseCase) Execute() (*dto.HealthResponse, error) {

	// Aqui podriamos agregar otras validaciones:
	// - Verificar conexión a DB
	// - Verificación con servicios externos
	// Si alguno falla, se retorna error HTTP 503
	return dto.NewHealthResponse(), nil
}
