package app

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/platform/interceptors"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/adapter/collector"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/config"
	grpccontroller "github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/controller/grpc"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/usecase"
	collectorpb "github.com/IliaSotnikov2005/golang-course/task4/repo-stat/proto/collector"
	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/proto/processor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
	conn       *grpc.ClientConn
}

func New(
	log *slog.Logger,
	cfgGRPC config.GRPCServer,
	collectorAddress string,
) (*App, error) {
	conn, err := grpc.NewClient(collectorAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	collectorRawClient := collectorpb.NewCollectorServiceClient(conn)

	collectorAdapter := collector.NewCollectorAdapter(collectorRawClient)

	getRepositoryUseCase := usecase.NewGetRepositoryUseCase(collectorAdapter)
	getSubscribtionsInfoUseCase := usecase.NewGetSubscriptionsInfoUseCase(collectorAdapter)
	pingUseCase := usecase.NewPingUseCase()

	gRPCHandler := grpccontroller.NewHandler(log, getRepositoryUseCase, getSubscribtionsInfoUseCase, pingUseCase)

	gRPCServer := grpc.NewServer(grpc.ConnectionTimeout(cfgGRPC.Timeout), grpc.ChainUnaryInterceptor(interceptors.LoggingInterceptor(log)))
	processor.RegisterProcessorServiceServer(gRPCServer, gRPCHandler)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       cfgGRPC.Port,
		conn:       conn,
	}, nil
}

func (a *App) MustRun() {
	go func() {
		if err := a.Run(); err != nil {
			panic(err)
		}
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

	a.conn.Close()
}
