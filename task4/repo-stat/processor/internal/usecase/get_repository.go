package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/processor/internal/domain"
)

type GetRepositoryUseCase struct {
	collector RepositoryProvider
}

func NewGetRepositoryUseCase(collector RepositoryProvider) *GetRepositoryUseCase {
	return &GetRepositoryUseCase{
		collector: collector,
	}
}

func (u *GetRepositoryUseCase) Execute(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	return u.collector.GetRepository(ctx, owner, repo)
}
