package usecase

import (
	"context"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/collector/internal/domain"
)

type GetRepositoryUseCase struct {
	githubClient GitHubClient
}

func NewGetRepositoryUseCase(gitHubClient GitHubClient) *GetRepositoryUseCase {
	return &GetRepositoryUseCase{
		githubClient: gitHubClient,
	}
}

func (uc *GetRepositoryUseCase) Execute(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if owner == "" {
		return nil, fmt.Errorf("%w: owner cannot be empty", domain.ErrInvalidInput)
	}
	if repo == "" {
		return nil, fmt.Errorf("%w: repo cannot be empty", domain.ErrInvalidInput)
	}

	repository, err := uc.githubClient.GetRepository(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository %s/%s: %w", owner, repo, err)
	}

	return repository, nil
}
