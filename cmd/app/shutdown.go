package app

import (
	"microservice/pkg/logger"
	"time"
)

func (a *App) Shutdown() {

	// Cerramos conexiones globales
	// Cierre DB, Colas, recursos , etc...

	logger.General("Cerrando conexiones y recursos...")

	if a.db != nil {
		a.db.Close()
	}

	if a.ps != nil {
		a.ps.Close()
	}

	time.Sleep(100 * time.Millisecond)

	logger.General("Recursos cerrados correctamente")
}
