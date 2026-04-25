package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/platform/interceptors"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/processor/internal/adapter/db"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/processor/internal/adapter/kafka"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/processor/internal/config"
	grpccontroller "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/processor/internal/controller/grpc"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/processor/internal/usecase"
	pb "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/proto/processor"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/grpc"
)

type App struct {
	log            *slog.Logger
	gRPCServer     *grpc.Server
	port           string
	pgPool         *pgxpool.Pool
	kafkaClient    *kgo.Client
	resultConsumer *kafka.ResultConsumer
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfgGRPC config.GRPCServer,
	databaseDSN string,
	cfgKafka config.KafkaConfig,
) (*App, error) {
	m, err := migrate.New(
		"file://migrations",
		databaseDSN,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to run up migrations: %w", err)
	}

	log.Info("migrations applied successfully")

	pool, err := pgxpool.New(ctx, databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	kClient, err := kgo.NewClient(
		kgo.SeedBrokers(cfgKafka.Brokers...),
		kgo.ConsumeTopics(cfgKafka.ResultsTopic),
		kgo.ConsumerGroup("processor-group"),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka client error: %w", err)
	}

	storage := db.NewPostgresRepository(pool)
	publisher := kafka.NewPublisher(kClient, cfgKafka.RequestsTopic)
	resultsConsumer := kafka.NewResultConsumer(kClient, storage, log)

	getRepoUC := usecase.NewGetRepositoryUseCase(storage, publisher)
	infoUC := usecase.NewGetSubscriptionsInfoUseCase(storage)
	pingUC := usecase.NewPingUseCase()

	grpcHandler := grpccontroller.NewHandler(log, getRepoUC, infoUC, pingUC)
	gRPCServer := grpc.NewServer(grpc.ConnectionTimeout(cfgGRPC.Timeout), grpc.ChainUnaryInterceptor(interceptors.LoggingInterceptor(log)))
	pb.RegisterProcessorServiceServer(gRPCServer, grpcHandler)

	return &App{
		log:            log,
		gRPCServer:     gRPCServer,
		port:           cfgGRPC.Port,
		pgPool:         pool,
		kafkaClient:    kClient,
		resultConsumer: resultsConsumer,
	}, nil
}

func (a *App) MustRun(ctx context.Context) {
	go func() {
		if err := a.Run(ctx); err != nil {
			panic(err)
		}
	}()
}

func (a *App) Run(ctx context.Context) error {
	go a.resultConsumer.Start(ctx)

	lis, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		return fmt.Errorf("Processor start error: %w", err)
	}

	a.log.Info("grpc server is running", slog.String("port", a.port))

	return a.gRPCServer.Serve(lis)
}

func (a *App) Stop() {
	a.log.Info("stopping application")
	a.gRPCServer.GracefulStop()
	if a.kafkaClient != nil {
		a.kafkaClient.Close()
	}

	a.pgPool.Close()
}
