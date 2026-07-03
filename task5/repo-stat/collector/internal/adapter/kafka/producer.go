package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/domain"
	dtokafka "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/platform/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Adapter struct {
	client *kgo.Client
	topic  string
	log    *slog.Logger
}

func NewAdapter(client *kgo.Client, topic string, log *slog.Logger) *Adapter {
	return &Adapter{client: client, topic: topic, log: log}
}

func (a *Adapter) Send(ctx context.Context, repo *domain.Repository, opErr error) error {
	log := a.log.With(slog.String("component", "kafka-producer"), slog.String("topic", a.topic))

	resp := a.buildResponse(repo, opErr)

	val, err := json.Marshal(resp)
	if err != nil {
		log.Error("failed to marshall response",
			slog.String("error", err.Error()))
		return fmt.Errorf("marshal error: %w", err)
	}

	log.Debug(
		"sending message to Kafka",
		slog.Int("size", len(val)),
		slog.Bool("has_errors", opErr != nil),
	)

	record := &kgo.Record{
		Topic: a.topic,
		Value: val,
		Headers: []kgo.RecordHeader{
			{Key: "content-type", Value: []byte("application/json")},
		},
	}

	if opErr != nil {
		record.Headers = append(record.Headers, kgo.RecordHeader{
			Key:   "error",
			Value: []byte("true"),
		})
	}

	result := a.client.ProduceSync(ctx, record)
	if err := result.FirstErr(); err != nil {
		a.log.Error(
			"produce failed",
			slog.String("error", err.Error()),
			slog.String("topic", a.topic),
		)
		return fmt.Errorf("kafka produce failed: %w", err)
	}

	return nil
}

func (a *Adapter) buildResponse(repo *domain.Repository, opErr error) dtokafka.RepoResponse {
	if opErr != nil {
		return dtokafka.RepoResponse{Error: opErr.Error()}
	}

	if repo == nil {
		return dtokafka.RepoResponse{Error: "repository is null"}
	}

	return dtokafka.RepoResponse{
		FullName:    repo.FullName,
		Description: repo.Description,
		Stargazers:  repo.Stargazers,
		Forks:       repo.Forks,
		CreatedAt:   repo.CreatedAt,
		HTMLURL:     repo.HTMLURL,
		Error:       "",
	}
}

func (a *Adapter) Close() {
	if a.client != nil {
		a.client.Close()
	}
}
