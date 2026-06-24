package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (c PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

type postgresDatabase struct {
	pool *pgxpool.Pool
}

func (d *postgresDatabase) Ping(ctx context.Context) error {
	return d.pool.Ping(ctx)
}

func (d *postgresDatabase) Close() {
	d.pool.Close()
}

func (d *postgresDatabase) Native() interface{} {
	return d.pool
}

func NewPostgresDatabase(ctx context.Context, cfg PostgresConfig) (Database, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN())

	if err != nil {
		return nil, fmt.Errorf("Error creando pool de DB: %w", err)
	}
	return &postgresDatabase{pool: pool}, nil
}
