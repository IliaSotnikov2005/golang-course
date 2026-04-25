package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/platform/logger"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/processor/internal/app"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/processor/internal/config"
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

	initCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	application, err := app.New(initCtx, log, cfg.GRPCServer, cfg.DatabaseDSN, cfg.KafkaConfig)
	if err != nil {
		log.Error("app init failed", slog.Any("error", err))
		os.Exit(1)
	}

	errChan := make(chan error, 1)
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	go func() {
		if err := application.Run(appCtx); err != nil {
			errChan <- err
		}
	}()

	log.Info("processor service started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Error("application error", slog.Any("error", err))
	case sig := <-stop:
		log.Info("received signal", slog.String("signal", sig.String()))
	}

	appCancel()
	application.Stop()
	log.Info("application gracefully stopped")
}
