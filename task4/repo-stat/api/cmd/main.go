package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/api/app"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/api/internal/config"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/platform/logger"
)

// @title           Repo Stat API
// @version         1.0
// @description     A service for collecting information of GitHub repositories.
// @host            localhost:8080
// @BasePath        /api
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

	application, err := app.New(log, cfg)
	if err != nil {
		log.Error("app init failed", slog.Any("error", err))
		os.Exit(1)
	}

	application.MustRun()

	log.Info("application started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.Stop()
	log.Info("application stopped")
}
