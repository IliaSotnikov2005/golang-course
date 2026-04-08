package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/platform/logger"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/platform/must"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/subscriber/internal/app"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/subscriber/internal/config"
)

func main() {
	cfg := must.Do(config.Load())
	log := must.Do(logger.MakeLogger(cfg.LogLevel))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	application := app.New(ctx, log, cfg)

	errChan := make(chan error, 1)

	go func() {
		if err := application.Run(); err != nil {
			errChan <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	log.Info("subscriber service started")

	select {
	case err := <-errChan:
		log.Error("critical error, stopping...", slog.String("err", err.Error()))
	case sig := <-stop:
		log.Info("received stop signal", slog.String("signal", sig.String()))
	}

	application.Stop()
	log.Info("subscriber service gracefully stopped")
}
