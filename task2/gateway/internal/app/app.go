package app

import (
	"log/slog"

	grpcclient "github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/adapter/grpc"
	v1 "github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/adapter/rest/v1"
	httpapp "github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/app/http"
	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/config"
	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/usecase"
)

type App struct {
	HTTPServer *httpapp.App
	GRPCClient *grpcclient.Client
	log        *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	grpcClient, err := grpcclient.NewClient(cfg.GRPC.Address, cfg.GRPC.Timeout, log)
	if err != nil {
		return nil, err
	}

	getRepoUseCase := usecase.NewGetRepositoryUseCase(grpcClient)
	restHandler := v1.NewHandler(getRepoUseCase)
	httpServer := httpapp.New(log, cfg, restHandler)

	return &App{
		HTTPServer: httpServer,
		GRPCClient: grpcClient,
		log:        log,
	}, nil
}

func (a *App) Stop() {
	if a.GRPCClient != nil {
		if err := a.GRPCClient.Close(); err != nil {
			a.log.Error("failed to close gRPC client: %v", slog.String("error", err.Error()))
		}
	}

	if a.HTTPServer != nil {
		a.HTTPServer.Stop()
	}
}
