package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/internal/domain"
)

type GetRepositoryUseCase struct {
	processor RepositoryProvider
}

func NewGetRepositoryUseCase(processor RepositoryProvider) *GetRepositoryUseCase {
	return &GetRepositoryUseCase{
		processor: processor,
	}
}

func (uc *GetRepositoryUseCase) Execute(ctx context.Context, url string) (*domain.Repository, error) {
	owner, repo, err := parseURL(url)
	if err != nil {
		return nil, err
	}

	return uc.processor.GetRepository(ctx, owner, repo)
}

func parseURL(rawURL string) (owner, repo string, err error) {
	parts := strings.Split(strings.TrimRight(rawURL, "/"), "/")
	if len(parts) < 2 {
		return "", "", errors.New("invalid github url")
	}
	return parts[len(parts)-2], parts[len(parts)-1], nil
}
