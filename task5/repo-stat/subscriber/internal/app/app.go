package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/platform/interceptors"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/subscriber/internal/adapters/db"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/subscriber/internal/adapters/github"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/subscriber/internal/adapters/kafka"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/subscriber/internal/config"
	grpccontroller "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/subscriber/internal/controller/grpc"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/subscriber/internal/usecase"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/grpc"

	pb "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/proto/subscriber"
)

type App struct {
	log         *slog.Logger
	gRPCServer  *grpc.Server
	pool        *pgxpool.Pool
	kafkaClient *kgo.Client
	port        string
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) (*App, error) {
	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseDSN,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to run up migrations: %w", err)
	}
	log.Info("migrations applied successfully")

	pool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	repo := db.NewPostgresRepository(pool)

	httpClient := &http.Client{Timeout: cfg.Github.Timeout}
	ghClient := github.NewClient(httpClient, cfg.Github.BaseURL, cfg.Github.UserAgent, log)

	kafkaClient, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Kafka.Brokers...),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka client error: %w", err)
	}

	eventSender := kafka.NewEventSender(kafkaClient, cfg.Kafka.Topic)

	subscribeUC := usecase.NewSubscribeUseCase(repo, ghClient, eventSender)
	unsubscribeUC := usecase.NewUnsubscribeUseCase(repo)
	listUC := usecase.NewListUseCase(repo)
	pingUC := usecase.NewPingUseCase()

	grpcHandler := grpccontroller.NewServer(log, subscribeUC, unsubscribeUC, listUC, pingUC)

	grpcServer := grpc.NewServer(grpc.ConnectionTimeout(cfg.GRPC.Timeout), grpc.ChainUnaryInterceptor(interceptors.LoggingInterceptor(log)))
	pb.RegisterSubscriberServer(grpcServer, grpcHandler)

	return &App{
		log:         log,
		gRPCServer:  grpcServer,
		pool:        pool,
		kafkaClient: kafkaClient,
		port:        cfg.GRPC.Port,
	}, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", a.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", a.port, err)
	}

	a.log.Info("grpc server is running",
		slog.String("addr", lis.Addr().String()),
	)

	if err := a.gRPCServer.Serve(lis); err != nil {
		if errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}

		return fmt.Errorf("serve error: %w", err)
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping subscriber app")
	a.gRPCServer.GracefulStop()
	if a.kafkaClient != nil {
		a.kafkaClient.Close()
	}

	a.pool.Close()
}
