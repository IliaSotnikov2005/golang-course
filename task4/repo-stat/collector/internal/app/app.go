package app

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/collector/internal/adapter/github"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/collector/internal/config"
	grpccontroller "github.com/IliaSotnikov2005/golang-course/task4/repo-stat/collector/internal/controller/grpc"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/collector/internal/usecase"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/platform/must"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/proto/collector"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(
	log *slog.Logger,
	cfgGRPC config.GRPCServer,
	cfgGithub config.GithubConfig,
) *App {

	httpClient := http.Client{
		Timeout: cfgGithub.Timeout,
	}

	githubClient := github.NewClient(
		&httpClient,
		cfgGithub.BaseURL,
		cfgGithub.UserAgent,
		log.With(slog.String("component", "github-client")),
	)

	getRepoUseCase := usecase.NewGetRepositoryUseCase(githubClient)
	pingUseCase := usecase.NewPingUseCase()

	grpcHandler := grpccontroller.NewHandler(log, getRepoUseCase, pingUseCase)

	gRPCServer := grpc.NewServer(
		grpc.ConnectionTimeout(cfgGRPC.Timeout),
		grpc.ChainUnaryInterceptor(grpccontroller.LoggingInterceptor(log)))
	collector.RegisterCollectorServiceServer(gRPCServer, grpcHandler)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       cfgGRPC.Port,
	}
}

func (a *App) MustRun() {
	go func() {
		must.NotError(a.Run())
	}()
}

func (a *App) Run() error {
	const operation = "app.Run"

	log := a.log.With(
		slog.String("operation", operation),
		slog.String("port", a.port),
	)

	lis, err := net.Listen("tcp", a.port)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("grpc server is running", slog.String("addr", lis.Addr().String()))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func (a *App) Stop() {
	const operation = "app.Stop"

	a.log.With(slog.String("operation", operation)).Info("stopping gRPC server", slog.String("port", a.port))

	a.gRPCServer.GracefulStop()
}
