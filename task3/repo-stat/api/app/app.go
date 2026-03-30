package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/internal/adapter/processor"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/internal/adapter/subscriber"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/internal/config"
	v1 "github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/internal/controller/http/v1"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/internal/usecase"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
	port       string
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) (*App, error) {
	processorClient, err := processor.NewClient(cfg.Services.Processor, log)
	if err != nil {
		return nil, fmt.Errorf("failed to init processor client: %w", err)
	}

	subscriberClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		return nil, fmt.Errorf("failed to init subscriber client: %w", err)
	}

	getRepoUC := usecase.NewGetRepositoryUseCase(processorClient)
	pingUC := usecase.NewPingUseCase(processorClient, subscriberClient)

	handler := v1.NewHandler(log, getRepoUC, pingUC)
	router := handler.Router()

	srv := &http.Server{
		Addr:         cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.Timeout,
		WriteTimeout: cfg.HTTP.Timeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	return &App{
		log:        log,
		httpServer: srv,
		port:       cfg.HTTP.Port,
	}, nil
}

func (a *App) MustRun() {
	go func() {
		if err := a.Run(); err != nil {
			panic(err)
		}
	}()
}

func (a *App) Run() error {
	a.log.Info("HTTP server is running", slog.String("addr", a.port))
	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server error: %w", err)
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
