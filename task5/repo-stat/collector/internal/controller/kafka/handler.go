package kafkacontroller

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/usecase"
	dtokafka "github.com/IliaSotnikov2005/golang-course/task5/repo-stat/platform/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Handler struct {
	log     *slog.Logger
	useCase *usecase.GetRepositoryUseCase
	sender  usecase.RepositoryResultSender
}

func NewHandler(log *slog.Logger, getRepositoryUseCase *usecase.GetRepositoryUseCase, sender usecase.RepositoryResultSender) *Handler {
	return &Handler{
		log:     log,
		useCase: getRepositoryUseCase,
		sender:  sender,
	}
}

func (h *Handler) Run(ctx context.Context, client *kgo.Client) {
	log := h.log.With(slog.String("component", "kafka-handler"))
	log.Info("starting Kafka handler")

	for {
		select {
		case <-ctx.Done():
			log.Info("context cancelled, stopping handler")
			return
		default:
		}

		ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)

		fetches := client.PollFetches(ctxWithTimeout)
		cancel()

		if err := fetches.Err(); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				log.Debug("poll timeout or cancelled", slog.String("error", err.Error()))
				continue
			}

			log.Error("poll fetch failed", slog.String("error", err.Error()))
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			h.processRecord(ctx, record)
		}
	}
}

func (h *Handler) processRecord(ctx context.Context, record *kgo.Record) {
	log := h.log.With(
		slog.String("topic", record.Topic),
		slog.Int64(
			"partition", int64(record.Partition),
		),
		slog.Int64("offset", record.Offset),
	)

	log.Debug("processing message")

	var req dtokafka.RepoRequest
	if err := json.Unmarshal(record.Value, &req); err != nil {
		log.Error("failed to unmarshal message",
			slog.String("error", err.Error()),
			slog.String("value", string(record.Value)))
	}

	log = log.With(
		slog.String("owner", req.Owner),
		slog.String("repo", req.Repo),
	)

	repo, err := h.useCase.Execute(ctx, req.Owner, req.Repo)

	if sendErr := h.sender.Send(ctx, repo, err); sendErr != nil {
		log.Error("failed to send result", slog.String("err", sendErr.Error()))
		return
	}

	if err != nil {
		log.Warn("repository fetch failed, error sent to response topic",
			slog.String("error", err.Error()))
		return
	}

	log.Info("message processed successfully", slog.String("repository", repo.FullName))
}
