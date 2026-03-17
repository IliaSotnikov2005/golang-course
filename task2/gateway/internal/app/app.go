package app

import (
	"log/slog"

	grpcclient "github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/adapter/grpc"
	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/adapter/rest"
	httpapp "github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/app/http"
	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/config"
	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/usecase"
)

type App struct {
	HTTPServer *httpapp.App
	GRPCClient *grpcclient.Client
}

func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	grpcClient, err := grpcclient.NewClient(cfg.Collector.Address, cfg.Collector.Timeout)
	if err != nil {
		return nil, err
	}

	getRepoUseCase := usecase.NewGetRepositoryUseCase(grpcClient)
	restHandler := rest.NewHandler(getRepoUseCase)
	httpServer := httpapp.New(log, cfg, restHandler)

	return &App{
		HTTPServer: httpServer,
		GRPCClient: grpcClient,
	}, nil
}

func (a *App) Stop() {
	if a.GRPCClient != nil {
		a.GRPCClient.Close()
	}
	if a.HTTPServer != nil {
		a.HTTPServer.Stop()
	}
}
