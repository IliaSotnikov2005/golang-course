package main

import (
	"os"
	"os/signal"
	"syscall"

	grpcapp "github.com/IliaSotnikov2005/golang-course/task2/collector/internal/app/grpc"
	"github.com/IliaSotnikov2005/golang-course/task2/collector/internal/config"
	"github.com/IliaSotnikov2005/golang-course/task2/collector/internal/logger/sl"
)

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}

	return obj
}

func main() {
	cfg := Must(config.Load())

	log := sl.SetupLogger(cfg)

	gRPCApp := grpcapp.New(log, cfg.GRPC, cfg.HTTP)

	gRPCApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	gRPCApp.Stop()
	log.Info("application stopped")
}
