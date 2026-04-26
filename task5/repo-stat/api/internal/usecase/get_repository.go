package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/api/internal/domain"
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
		return nil, fmt.Errorf("%w: %s", domain.ErrInvalidInput, err.Error())
	}

	return uc.processor.GetRepository(ctx, owner, repo)
}

func parseURL(rawURL string) (string, string, error) {
	input := strings.TrimSpace(strings.TrimRight(rawURL, "/"))

	if strings.Contains(input, "://") {
		parts := strings.SplitN(input, "://", 2)
		input = parts[1]
	}

	input = strings.TrimPrefix(input, "github.com")
	input = strings.TrimPrefix(input, "www.github.com")
	input = strings.Trim(input, "/")

	parts := strings.Split(input, "/")

	if len(parts) < 2 {
		return "", "", fmt.Errorf("expected format: repository URL or owner/repo, got: %s", rawURL)
	}

	if len(parts) > 2 {
		return "", "", fmt.Errorf("invalid URL: too many path segments, expected owner/repo")
	}

	return parts[0], parts[1], nil
}
