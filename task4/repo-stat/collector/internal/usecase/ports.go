package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/collector/internal/domain"
)

type GitHubClient interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type SubscriberClient interface {
	GetSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}
