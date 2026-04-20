package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/platform/grpcserver"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/platform/logger"
	subscriberpb "github.com/IliaSotnikov2005/golang-course/task3/repo-stat/proto/subscriber"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/subscriber/config"
	grpccontroller "github.com/IliaSotnikov2005/golang-course/task3/repo-stat/subscriber/internal/controller/grpc"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/subscriber/internal/usecase"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log, err := logger.MakeLogger(cfg.Logger.LogLevel)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "logger init error:", err)
		return err
	}

	log.Info("starting subscriber server...")
	log.Debug("debug messages are enabled")

	pingUseCase := usecase.NewPing()
	pingServer := grpccontroller.NewServer(log, pingUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	subscriberpb.RegisterSubscriberServer(srv.GRPC(), pingServer)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
}
