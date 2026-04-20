package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/collector/internal/app"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/collector/internal/config"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/platform/logger"
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

	application := app.New(log, cfg.GRPC, cfg.Github, cfg.Subscriber)

	application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.Stop()
	log.Info("application stopped")
}
