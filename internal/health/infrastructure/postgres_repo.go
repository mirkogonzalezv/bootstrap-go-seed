package infrastructure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresHealthRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresHealthRepository(pool *pgxpool.Pool) *PostgresHealthRepository {
	return &PostgresHealthRepository{pool: pool}
}

// Implementamos la interfaz del repository
func (r *PostgresHealthRepository) Ping(ctx context.Context) error {
	if err := r.pool.Ping(ctx); err != nil {
		return fmt.Errorf("Health check DB: %w", err)
	}
	return nil
}
