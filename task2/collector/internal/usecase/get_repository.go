package usecase

import (
	"context"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task2/collector/internal/domain"
)

type GetRepositoryUseCase struct {
	githubClient domain.GitHubClient
}

func NewGetRepositoryUseCase(gitHubClient domain.GitHubClient) *GetRepositoryUseCase {
	return &GetRepositoryUseCase{
		githubClient: gitHubClient,
	}
}

func (uc *GetRepositoryUseCase) Execute(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if owner == "" {
		return nil, fmt.Errorf("owner cannot be empty")
	}
	if repo == "" {
		return nil, fmt.Errorf("repo cannot be empty")
	}

	repository, err := uc.githubClient.GetRepository(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository %s/%s: %w", owner, repo, err)
	}

	return repository, nil
}
