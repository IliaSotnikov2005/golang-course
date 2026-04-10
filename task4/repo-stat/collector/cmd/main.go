package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/collector/internal/app"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/collector/internal/config"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/platform/logger"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/platform/must"
)

func main() {
	cfg := must.Do(config.Load())

	log := must.Do(logger.MakeLogger(cfg.LogLevel))

	application := app.New(log, cfg.GRPC, cfg.Github, cfg.Subscriber)

	application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.Stop()
	log.Info("application stopped")
}
