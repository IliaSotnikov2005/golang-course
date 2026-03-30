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

// @title           Repo Stat API
// @version         1.0
// @description     A service for collecting information of GitHub repositories.
// @host            localhost:8080
// @BasePath        /api
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
