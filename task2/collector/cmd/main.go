package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	grpcapp "github.com/IliaSotnikov2005/golang-course/task2/collector/internal/app/grpc"
	"github.com/IliaSotnikov2005/golang-course/task2/collector/internal/config"
)

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}

	return obj
}

func setupLogger(level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

func getLogLevel(cfg *config.Config) slog.Level {
	switch cfg.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func main() {
	cfg := Must(config.Load())

	log := setupLogger(getLogLevel(cfg))

	gRPCApp := grpcapp.New(log, cfg.GRPC, cfg.HTTP)

	go func() {
		gRPCApp.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	gRPCApp.Stop()
	log.Info("application stopped")
}
