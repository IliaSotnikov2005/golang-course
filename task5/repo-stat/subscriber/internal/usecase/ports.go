package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/subscriber/internal/domain"
)

type SubscriptionRepository interface {
	Save(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	Delete(ctx context.Context, owner, repo string) error
	List(ctx context.Context) ([]domain.Subscription, error)
}

type GithubClient interface {
	Exists(ctx context.Context, owner, repo string) (bool, error)
}
