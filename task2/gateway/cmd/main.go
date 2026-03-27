package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/app"
	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/config"
	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/logger/sl"
)

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}

	return obj
}

// @title           GitHub Collector API
// @version         1.0
// @description     API Gateway for collecting GitHub repository information.
// @host            localhost:8080
// @BasePath        /
func main() {
	cfg := Must(config.Load())

	log := sl.SetupLogger(cfg)

	application := Must(app.New(log, cfg))

	application.HTTPServer.MustRun()

	log.Info("gateway started", slog.String("port", cfg.HTTP.Port))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
	log.Info("gateway stopped")
}
