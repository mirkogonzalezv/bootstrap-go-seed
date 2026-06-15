package app

import (
	"microservice/pkg/logger"
	"time"
)

func (a *App) Shutdown() {
	logger.General("Cerrando conexiones y recursos...")

	if a.container != nil && a.container.DB != nil {
		a.container.DB.Close()
		logger.Success("Pool de DB cerrado")
	}

	// Cerramos conexiones globales
	// Cierre DB, Colas, recursos , etc...

	time.Sleep(100 * time.Millisecond)

	logger.General("Recursos cerrados correctamente")
}
