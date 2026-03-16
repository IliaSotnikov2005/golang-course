package app

import (
	"log/slog"

	grpcapp "github.com/IliaSotnikov2005/golang-course/task2/collector/internal/app/grpc"
	"github.com/IliaSotnikov2005/golang-course/task2/collector/internal/config"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	gRPCApp := grpcapp.New(log, cfg.GRPC, cfg.HTTP)

	return &App{
		GRPCServer: gRPCApp,
	}
}
