package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/collector/internal/domain"
)

type GitHubClient interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}
