package httputil

import (
	"microservice/pkg/apperrors"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error     ErrorDetail `json:"error"`
	Timestamp string      `json:"timestamp"`
}

type ErrorDetail struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// MapErrorToHttp convierte errores de dominio a respuestas HTTP
func MapErrorToHttp(err error) (int, ErrorResponse) {
	if domainErr, ok := err.(*apperrors.DomainError); ok {
		statusCode := getHttpStatusCode(domainErr.Type)

		return statusCode, ErrorResponse{
			Error: ErrorDetail{
				Type:    string(domainErr.Type),
				Code:    domainErr.Code,
				Message: domainErr.Message,
				Details: domainErr.Details,
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Error genérico no controlado
	return http.StatusInternalServerError, ErrorResponse{
		Error: ErrorDetail{
			Type:    string(apperrors.ErrorTypeInternal),
			Code:    "INTERNAL_001",
			Message: "Internal server error",
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func getHttpStatusCode(errorType apperrors.ErrorType) int {
	switch errorType {
	case apperrors.ErrorTypeValidation:
		return http.StatusBadRequest
	case apperrors.ErrorTypeNotFound:
		return http.StatusNotFound
	case apperrors.ErrorTypeConflict:
		return http.StatusConflict
	case apperrors.ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case apperrors.ErrorTypeForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
