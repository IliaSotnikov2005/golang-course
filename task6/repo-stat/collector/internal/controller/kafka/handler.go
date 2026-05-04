package kafkacontroller

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/collector/internal/usecase"
	dtokafka "github.com/IliaSotnikov2005/golang-course/task6/repo-stat/platform/kafka"
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
	for {
		fetches := client.PollFetches(ctx)
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()

			var req dtokafka.RepoRequest
			json.Unmarshal(record.Value, &req)

			repo, err := h.useCase.Execute(ctx, req.Owner, req.Repo)

			if sendErr := h.sender.Send(ctx, repo, err); sendErr != nil {
				h.log.Error("failed to send result", "err", sendErr)
			}
		}
	}
}
