package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/collector/internal/domain"
)

type GitHubClient interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type SubscriberClient interface {
	GetSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}

type RepositoryResultSender interface {
	Send(ctx context.Context, repo *domain.Repository, err error) error
}
