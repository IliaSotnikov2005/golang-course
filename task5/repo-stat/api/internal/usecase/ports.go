package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/api/internal/domain"
)

type RepositoryProvider interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
	GetSubscriptionsInfo(ctx context.Context) ([]domain.Repository, error)
}

type Subscriber interface {
	Subscribe(ctx context.Context, owner, repo string) error
	Unsubscribe(ctx context.Context, owner, repo string) error
	List(ctx context.Context) ([]domain.Subscription, error)
}

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
}
