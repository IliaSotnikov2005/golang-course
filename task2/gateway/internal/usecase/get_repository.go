package usecase

import (
	"context"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/domain"
)

type GetRepositoryUseCase struct {
	repoProvider domain.RepositoryProvider
}

func NewGetRepositoryUseCase(repoProvider domain.RepositoryProvider) *GetRepositoryUseCase {
	return &GetRepositoryUseCase{
		repoProvider: repoProvider,
	}
}

func (uc *GetRepositoryUseCase) Execute(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if owner == "" {
		return nil, fmt.Errorf("%w: owner cannot be empty", domain.ErrInvalidInput)
	}
	if repo == "" {
		return nil, fmt.Errorf("%w: repo cannot be empty", domain.ErrInvalidInput)
	}

	repository, err := uc.repoProvider.GetRepository(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return repository, nil
}
