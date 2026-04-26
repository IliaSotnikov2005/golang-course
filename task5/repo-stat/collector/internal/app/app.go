package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/adapter/github"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/adapter/kafka"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/adapter/subscriber"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/config"
	kafkacontroller "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/controller/kafka"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/usecase"
	subscriberpb "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/proto/subscriber"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	log               *slog.Logger
	kafkaClient       *kgo.Client
	kafkaHandler      *kafkacontroller.Handler
	taskDispatcher    *kafka.Dispatcher
	subscriberAdapter *subscriber.Client
	subscriberConn    *grpc.ClientConn
}

func New(
	log *slog.Logger,
	cfgGithub config.GithubConfig,
	cfgSubscriber config.SubscriberConfig,
	cfgKafka config.KafkaConfig,
) *App {

	githubClient := github.NewClient(
		&http.Client{Timeout: cfgGithub.Timeout},
		cfgGithub.BaseURL,
		cfgGithub.UserAgent,
		log.With(slog.String("component", "github-client")),
	)

	conn, err := grpc.NewClient(cfgSubscriber.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	subscriberRawClient := subscriberpb.NewSubscriberClient(conn)
	subscriberAdapter := subscriber.NewClient(subscriberRawClient)

	kafkaClient, err := kgo.NewClient(
		kgo.SeedBrokers(cfgKafka.Brokers...),
		kgo.ConsumerGroup(cfgKafka.GroupID),
		kgo.ConsumeTopics(cfgKafka.RequestTopic),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create kafka client: %w", err))
	}

	resultProducer := kafka.NewAdapter(kafkaClient, cfgKafka.ResponseTopic)
	taskDispatcher := kafka.NewDispatcher(kafkaClient, cfgKafka.RequestTopic)

	getRepoUC := usecase.NewGetRepositoryUseCase(githubClient)

	handler := kafkacontroller.NewHandler(log, getRepoUC, resultProducer)

	return &App{
		log:               log,
		kafkaClient:       kafkaClient,
		kafkaHandler:      handler,
		taskDispatcher:    taskDispatcher,
		subscriberAdapter: subscriberAdapter,
		subscriberConn:    conn,
	}
}

func (a *App) MustRun(ctx context.Context) {
	go func() {
		if err := a.Run(ctx); err != nil {
			panic(err)
		}
	}()
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("collector application is starting")

	go a.runBackgroundUpdater(ctx)

	a.kafkaHandler.Run(ctx, a.kafkaClient)

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping application")
	a.subscriberConn.Close()
	a.kafkaClient.Close()
}
