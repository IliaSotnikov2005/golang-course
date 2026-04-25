package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/app"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/config"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/platform/logger"
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

	application := app.New(log, cfg.Github, cfg.Subscriber, cfg.Kafka)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := application.Run(ctx); err != nil {
			log.Error("application run error", "err", err)
			cancel()
		}
	}()

	log.Info("collector is running")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Info("shutting down...")
	cancel()
	application.Stop()

	log.Info("application stopped")
}
