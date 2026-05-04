package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/processor/internal/domain"
)

type DataStorage interface {
	GetByFullName(ctx context.Context, fullName string) (*domain.Repository, error)
	Upsert(ctx context.Context, repo *domain.Repository) error
	ListAll(ctx context.Context) ([]domain.Repository, error)
}

type EventPublisher interface {
	PublishFetchRequest(ctx context.Context, owner, repo string) error
}
