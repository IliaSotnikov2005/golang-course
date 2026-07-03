package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	dtokafka "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/platform/kafka"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/processor/internal/domain"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/processor/internal/usecase"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ResultConsumer struct {
	client  *kgo.Client
	storage usecase.DataStorage
	log     *slog.Logger
}

func NewResultConsumer(client *kgo.Client, storage usecase.DataStorage, log *slog.Logger) *ResultConsumer {
	return &ResultConsumer{
		client:  client,
		storage: storage,
		log:     log,
	}
}

func (c *ResultConsumer) Run(ctx context.Context) {
	log := c.log.With(slog.String("component", "kafka-result-consumer"))
	log.Info("starting Kafka consumer")

	for {
		select {
		case <-ctx.Done():
			log.Info("context cancelled, stopping consumer")
			return
		default:
		}

		ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
		fetches := c.client.PollFetches(ctxWithTimeout)
		cancel()

		if errs := fetches.Errors(); len(errs) > 0 {
			c.log.Error("kafka poll errors", slog.Any("errors", errs))
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			c.processRecord(ctx, record)
		}
	}
}

func (c *ResultConsumer) processRecord(ctx context.Context, record *kgo.Record) {
	log := c.log.With(
		slog.String("topic", record.Topic),
		slog.Int64("partition", int64(record.Partition)),
		slog.Int64("offset", record.Offset),
	)

	log.Debug("processing message")

	var response dtokafka.RepoResponse
	if err := json.Unmarshal(record.Value, &response); err != nil {
		log.Error("failed to unmarshal message", slog.String("error", err.Error()), slog.String("value", string(record.Value)))
	}

	if response.Error != "" {
		c.log.Warn("received error response from collector", slog.String("error", response.Error))
		return
	}

	repo := &domain.Repository{
		FullName:    response.FullName,
		Description: response.Description,
		Stargazers:  response.Stargazers,
		Forks:       response.Forks,
		CreatedAt:   response.CreatedAt,
		HTMLURL:     response.HTMLURL,
	}

	if err := c.storage.Upsert(ctx, repo); err != nil {
		c.log.Error("failed to upsert repo from kafka", slog.String("repo", repo.FullName), slog.String("error", err.Error()))
		return
	}

	log.Info("message processed successfully")
}
