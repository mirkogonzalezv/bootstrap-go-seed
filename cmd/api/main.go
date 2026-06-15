// @title Microservice API
// @version 1.0
// @description API de ejemplo basada en Clean Architecture
// @host localhost:3200
// @BasePath /api/v1
package main

import (
	"log"
	"microservice/cmd/app"
)

func main() {
	app := app.NewApp()

	if err := app.Init(); err != nil {
		log.Fatal("Error al iniciar: ", err)
	}

	if err := app.Run(); err != nil {
		log.Fatal("Error al ejecutar aplicativo: ", err)
	}
}
