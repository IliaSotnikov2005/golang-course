package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/app"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/internal/config"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/platform/logger"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/platform/must"
)

func main() {
	cfg := must.Do(config.Load())

	log := must.Do(logger.MakeLogger(cfg.LogLevel))

	application := must.Do(app.New(log, cfg))
	application.MustRun()

	log.Info("application started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.Stop()
	log.Info("application stopped")
}
