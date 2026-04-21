package main

import (
	"log"
	"microservice/cmd/app"
)

func main() {
	app := app.NewApp()

	if err := app.Initializer(); err != nil {
		log.Fatal("Error al iniciar: ", err)
	}

	if err := app.Run(); err != nil {
		log.Fatal("Error al ejecutar aplicativo: ", err)
	}
}
