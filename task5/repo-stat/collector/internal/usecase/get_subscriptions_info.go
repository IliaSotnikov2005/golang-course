package usecase

import (
	"context"
	"sync"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/collector/internal/domain"
	"golang.org/x/sync/errgroup"
)

const maxConcurrentRequests = 10

type GetSubscriptionsInfoUseCase struct {
	subscriberClient SubscriberClient
	githubClient     GitHubClient
}

func NewGetSubscriptionsInfoUseCase(sub SubscriberClient, gh GitHubClient) *GetSubscriptionsInfoUseCase {
	return &GetSubscriptionsInfoUseCase{
		subscriberClient: sub,
		githubClient:     gh,
	}
}

func (uc *GetSubscriptionsInfoUseCase) Execute(ctx context.Context) ([]domain.Repository, error) {
	subs, err := uc.subscriberClient.GetSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrentRequests)

	results := make([]domain.Repository, 0, len(subs))
	var mu sync.Mutex

	for _, s := range subs {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			repo, err := uc.githubClient.GetRepository(ctx, s.Owner, s.Repo)
			if err != nil {
				return nil
			}

			mu.Lock()
			results = append(results, *repo)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}
