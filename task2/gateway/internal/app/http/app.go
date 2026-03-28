package httpapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	v1 "github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/adapter/rest/v1"
	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/config"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
	port       string
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	handler *v1.Handler,
) *App {
	httpServer := &http.Server{
		Addr:         cfg.HTTP.Port,
		Handler:      handler.Router(),
		ReadTimeout:  cfg.HTTP.Timeout,
		WriteTimeout: cfg.HTTP.Timeout,
		IdleTimeout:  120 * time.Second,
	}

	return &App{
		log:        log,
		httpServer: httpServer,
		port:       cfg.HTTP.Port,
	}
}

func (a *App) MustRun() {
	go func() {
		if err := a.Run(); err != nil {
			panic(err)
		}
	}()
}

func (a *App) Run() error {
	const operation = "httpapp.Run"

	log := a.log.With(
		slog.String("operation", operation),
		slog.String("port", a.port),
	)

	log.Info("HTTP server is starting")

	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func (a *App) Stop() {
	const operation = "httpapp.Stop"

	a.log.With(slog.String("operation", operation)).Info("stopping HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.log.Error("failed to shutdown HTTP server", slog.String("error", err.Error()))
	}
}
