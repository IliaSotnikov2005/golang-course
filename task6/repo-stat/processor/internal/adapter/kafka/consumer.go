package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	dtokafka "github.com/IliaSotnikov2005/golang-course/task6/repo-stat/platform/kafka"
	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/processor/internal/domain"
	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/processor/internal/usecase"
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

func (c *ResultConsumer) Start(ctx context.Context) {
	for {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			c.log.Error("kafka poll errors", slog.Any("errors", errs))
		}

		fetches.EachRecord(func(record *kgo.Record) {
			var resp dtokafka.RepoResponse
			if err := json.Unmarshal(record.Value, &resp); err != nil {
				c.log.Error("failed to unmarshal kafka record", slog.Any("error", err))
				return
			}

			if resp.Error != "" {
				c.log.Warn("received error response from collector", slog.String("error", resp.Error))
				return
			}

			repo := &domain.Repository{
				FullName:    resp.FullName,
				Description: resp.Description,
				Stargazers:  resp.Stargazers,
				Forks:       resp.Forks,
				CreatedAt:   resp.CreatedAt,
				HTMLURL:     resp.HTMLURL,
			}

			if err := c.storage.Upsert(ctx, repo); err != nil {
				c.log.Error("failed to upsert repo from kafka", slog.String("repo", repo.FullName), slog.Any("error", err))
			}
		})
	}
}
