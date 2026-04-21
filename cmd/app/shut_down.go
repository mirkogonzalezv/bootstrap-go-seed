package app

import (
	"microservice/pkg/logger"
	"time"
)

func (a *App) ShutDown() {
	logger.General("Cerrando conexiones y recursos...")

	// Cerramos conexiones globales
	// Cierre DB, Colas, recursos , etc...

	time.Sleep(100 * time.Millisecond)

	logger.General("Recursos cerrados correctamente")
}
