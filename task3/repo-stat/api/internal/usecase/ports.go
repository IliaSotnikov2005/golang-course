package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/internal/domain"
)

type RepositoryProvider interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
}
