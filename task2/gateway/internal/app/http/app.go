package httpapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/adapter/rest"
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
	handler *rest.Handler,
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
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)

	log.Info("HTTP server is starting")

	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "httpapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.log.Error("failed to shutdown HTTP server", slog.String("error", err.Error()))
	}
}
