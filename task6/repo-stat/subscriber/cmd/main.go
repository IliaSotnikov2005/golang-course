package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/platform/logger"
	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/subscriber/internal/app"
	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/subscriber/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		_, _ = os.Stderr.WriteString("config load error: " + err.Error() + "\n")
		os.Exit(1)
	}

	log, err := logger.MakeLogger(cfg.LogLevel)
	if err != nil {
		_, _ = os.Stderr.WriteString("logger init error: " + err.Error() + "\n")
		os.Exit(1)
	}

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	application, err := app.New(appCtx, log, cfg)
	if err != nil {
		_, _ = os.Stderr.WriteString("application init error: " + err.Error() + "\n")
		os.Exit(1)
	}

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
