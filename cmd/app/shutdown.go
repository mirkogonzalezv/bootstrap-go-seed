package app

import (
	"microservice/pkg/logger"
	"time"
)

func (a *App) Shutdown() {
	logger.General("Cerrando conexiones y recursos...")

	if a.db != nil {
		a.db.Close()
	}

	// Cerramos conexiones globales
	// Cierre DB, Colas, recursos , etc...

	time.Sleep(100 * time.Millisecond)

	logger.General("Recursos cerrados correctamente")
}
