package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(environment string) {
	appEnv := os.Getenv("ENVIRONMENT")

	if appEnv == "" {
		appEnv = "dev"
	}

	if appEnv == "dev" || appEnv == "qa" {
		filename := ".env." + appEnv

		if _, err := os.Stat(filename); err == nil {
			fmt.Printf("Cargando %s...", filename)
			if err := godotenv.Load(filename); err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Archivo %s no encontrado, usando variables del sistema\n", filename)
		}
	}

	fmt.Println("ENVIRONMENT: ", appEnv)

}
