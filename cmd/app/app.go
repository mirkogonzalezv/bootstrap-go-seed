package app

import (
	"microservice/cmd/config"
	"microservice/cmd/container"
	"microservice/cmd/routes"
	"microservice/pkg/logger"
	"microservice/pkg/middlewares"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	config *config.Configuration
	router *gin.Engine
	log    *zap.Logger
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initializer() error {
	env := os.Getenv("ENVIRONMENT")

	if env == "" {
		env = "dev"
	}

	config.LoadEnv(env)

	logger.Init(env)
	logger.General("==== Microservice Start ====")

	cfg, err := config.LoadVars()
	if err != nil {
		return err
	}

	a.config = cfg

	logger.Success("Configuración cargada correctamente")
	logger.General("Servidor Configurado")

	router := gin.New()

	// Aplicamos Middlewares
	router.Use(gin.Logger())

	router.Use(middlewares.ErrorHandlerMiddleware()) // Global Error Handler

	router.Use(gin.Recovery())

	// Middleware Cors
	router.Use(cors.Default())

	// Otros middlewares
	//...

	container := container.NewContainer(a.config)
	apiRouter := routes.NewAPIRouter(container)
	apiRouter.RegisterRoutes(router)

	a.router = router

	return nil
}
