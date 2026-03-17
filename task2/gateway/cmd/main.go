package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/app"
	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/config"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	application, err := app.New(log, cfg)
	if err != nil {
		log.Error("failed to create app", slog.String("error", err.Error()))
		os.Exit(1)
	}

	go func() {
		application.HTTPServer.MustRun()
	}()

	log.Info("gateway started", slog.String("port", cfg.HTTP.Port))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
	log.Info("application stopped")
}
