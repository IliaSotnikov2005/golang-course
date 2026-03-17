package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/IliaSotnikov2005/golang-course/task2/collector/internal/adapter/github"
	grpcadapter "github.com/IliaSotnikov2005/golang-course/task2/collector/internal/adapter/grpc"
	collectorpb "github.com/IliaSotnikov2005/golang-course/task2/collector/internal/api/proto/gen"
	"github.com/IliaSotnikov2005/golang-course/task2/collector/internal/config"
	"github.com/IliaSotnikov2005/golang-course/task2/collector/internal/usecase"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(
	log *slog.Logger,
	cfgGRPC config.GRPCConfig,
	cfgHTTP config.HTTPConfig,
) *App {

	httpClient := http.Client{
		Timeout: cfgHTTP.Timeout,
	}

	githubClient := github.NewClient(
		httpClient,
		cfgHTTP.BaseURL,
		cfgHTTP.UserAgent,
		log.With(slog.String("component", "github-client")),
	)

	getRepoUseCase := usecase.NewGetRepositoryUseCase(githubClient)
	grpcHandler := grpcadapter.NewHandler(getRepoUseCase)

	gRPCServer := grpc.NewServer(grpc.ConnectionTimeout(cfgGRPC.Timeout))
	collectorpb.RegisterCollectorServiceServer(gRPCServer, grpcHandler)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       cfgGRPC.Port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)

	lis, err := net.Listen("tcp", a.port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", lis.Addr().String()))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.String("port", a.port))

	a.gRPCServer.GracefulStop()
}
