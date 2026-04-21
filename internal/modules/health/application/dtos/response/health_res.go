package dtos

import "time"

type HealthResponse struct {
	Timestamp time.Time `json:"timestamp" example:"2024-01-15T10.30:00Z"`
	Service   string    `json:"service" example:"microservice-seed-go"`
}

// Constructor
func NewHealthResponse() *HealthResponse {
	return &HealthResponse{
		Timestamp: time.Now(),
		Service:   "microservice",
	}
}
