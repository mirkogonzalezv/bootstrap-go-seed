package app

import (
	"context"
	"microservice/cmd/config"
	"microservice/cmd/container"
	"microservice/cmd/routes"
	"microservice/pkg/database"
	"microservice/pkg/logger"
	"microservice/pkg/middleware"
	"os"

	_ "microservice/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type App struct {
	config    *config.Configuration
	container *container.Container
	router    *gin.Engine
	log       *zap.Logger
	db        database.Database
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init() error {
	env := os.Getenv("ENVIRONMENT")

	if env == "" {
		env = "dev"
	}

	config.LoadEnv(env)

	logger.Init(env)
	a.log = logger.L()
	logger.General("==== Microservice Start ====")

	cfg, err := config.LoadVars()
	if err != nil {
		return err
	}

	a.config = cfg

	logger.Success("Configuración cargada correctamente")
	logger.General("Servidor Configurado")

	router := gin.New()

	if a.config.AppEnv != "prod" {
		router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Aplicamos Middlewares
	router.Use(gin.Logger())

	router.Use(middleware.ErrorHandlerMiddleware(a.log)) // Global Error Handler

	router.Use(gin.Recovery())

	// Middleware Cors
	router.Use(cors.Default())

	// Middleware Helmet
	router.Use(middleware.SecurityHeaders())

	// Otros middlewares
	//...

	// Base de datos
	db, err := database.NewPostgresDatabase(context.Background(), database.PostgresConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		DBName:   cfg.DBName,
		Password: cfg.DBPassword,
		SSLMode:  cfg.DBSSLMode,
	})

	if err != nil {
		return err
	}

	a.db = db

	c, err := container.NewContainer(a.config, db, a.log)
	if err != nil {
		return err
	}
	a.container = c

	apiRouter := routes.NewAPIRouter(a.container)
	apiRouter.RegisterRoutes(router)

	a.router = router

	return nil
}
