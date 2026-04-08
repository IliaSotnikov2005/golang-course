package app

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/platform/must"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/adapter/collector"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/config"
	grpccontroller "github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/controller/grpc"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/usecase"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/proto/processor"
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
	collectorAddress string,
) (*App, error) {
	collectorClient, err := collector.NewCollectorAdapter(collectorAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create collector client: %w", err)
	}

	getRepositoryUseCase := usecase.NewGetRepositoryUseCase(collectorClient)
	pingUseCase := usecase.NewPingUseCase()

	gRPCHandler := grpccontroller.NewHandler(log, getRepositoryUseCase, pingUseCase)

	gRPCServer := grpc.NewServer(grpc.ConnectionTimeout(cfgGRPC.Timeout))
	processor.RegisterProcessorServiceServer(gRPCServer, gRPCHandler)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       cfgGRPC.Port,
	}, nil
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
