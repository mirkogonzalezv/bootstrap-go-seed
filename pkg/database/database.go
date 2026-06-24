package database

import "context"

type Database interface {
	Ping(ctx context.Context) error
	Close()
	// Expone la conexión navita para que los repositorios la type-aseen
	Native() interface{}
}
