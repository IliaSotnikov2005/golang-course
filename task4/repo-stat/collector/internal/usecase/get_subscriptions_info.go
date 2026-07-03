package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/collector/internal/domain"
	"golang.org/x/sync/errgroup"
)

const maxConcurrentRequests = 10

type GetSubscriptionsInfoUseCase struct {
	log              *slog.Logger
	subscriberClient SubscriberClient
	githubClient     GitHubClient
}

func NewGetSubscriptionsInfoUseCase(log *slog.Logger, sub SubscriberClient, gh GitHubClient) *GetSubscriptionsInfoUseCase {
	return &GetSubscriptionsInfoUseCase{
		log:              log,
		subscriberClient: sub,
		githubClient:     gh,
	}
}

func (uc *GetSubscriptionsInfoUseCase) Execute(ctx context.Context) ([]domain.Repository, error) {
	logger := uc.log.With(slog.String("operation", "GetSubscriptionsInfoUseCase.Execute"))

	subs, err := uc.subscriberClient.GetSubscriptions(ctx)
	if err != nil {
		logger.Error("failed to get user subscriptions", slog.String("error", err.Error()))
		return nil, fmt.Errorf("get subscriptions error: %w", err)
	}

	if len(subs) == 0 {
		logger.Debug("no subscriptions found")
		return []domain.Repository{}, nil
	}

	logger.Debug("fetching repositories", slog.Int("count", len(subs)))

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrentRequests)

	resultCh := make(chan domain.Repository, len(subs))

	var errMu sync.Mutex
	var errList []error

	for _, s := range subs {
		g.Go(func() error {
			log := logger.With(
				slog.String("owner", s.Owner),
				slog.String("repo", s.Repo),
			)
			log.Debug("fetching repository")

			repo, err := uc.githubClient.GetRepository(ctx, s.Owner, s.Repo)
			if err != nil {
				log.Error("failed to fetch repository", slog.String("error", err.Error()))

				errMu.Lock()
				errList = append(errList, err)
				errMu.Unlock()

				return nil
			}

			log.Debug("repository fetched")
			resultCh <- *repo

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		logger.Error("group wait failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("group wait error: %w", err)
	}

	close(resultCh)

	results := make([]domain.Repository, 0, len(resultCh))
	for repo := range resultCh {
		results = append(results, repo)
	}

	if len(errList) > 0 {
		logger.Warn(
			"some repositories failed",
			slog.Int64("failed", int64(len(errList))),
			slog.Int("total", len(subs)),
		)

		if len(results) == 0 {
			return nil, fmt.Errorf("all requests failed: %v", errList)
		}

		return results, nil
	}

	logger.Info(
		"successfully fetched all repositories",
		slog.Int("count", len(results)),
	)

	return results, nil
}
