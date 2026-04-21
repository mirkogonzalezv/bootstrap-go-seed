package app

import (
	"context"
	"microservice/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func (a *App) Run() error {
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(a.config.Port),
		Handler: a.router,
	}

	go func() {
		logger.Success("Microservice iniciado correctamente en el puerto: " + strconv.Itoa(a.config.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error al iniciar servidor: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	sig := <-quit
	logger.General("Terminando microservice..." + sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.ShutDownTimeOut)*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Error al cerrar servidor: " + err.Error())
		return err
	}

	a.ShutDown()

	logger.L().Sync()
	logger.Success("Servidor terminado correctamente")
	return nil
}
