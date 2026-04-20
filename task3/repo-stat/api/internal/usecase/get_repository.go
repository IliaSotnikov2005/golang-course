package usecase

import (
	"context"
	"fmt"
	"net/url"
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
		return nil, fmt.Errorf("%w: %s", domain.ErrInvalidInput, err.Error())
	}

	return uc.processor.GetRepository(ctx, owner, repo)
}

func parseURL(rawURL string) (owner, repo string, err error) {
	input := strings.TrimSpace(strings.TrimRight(rawURL, "/"))

	if !strings.Contains(input, "://") {
		input = "https://" + input
	}

	u, err := url.Parse(input)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL format")
	}

	path := u.Path
	if u.Host != "" && !strings.Contains(u.Host, "github.com") {
		return "", "", fmt.Errorf("only GitHub repositories are supported")
	}

	path = strings.TrimPrefix(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		return "", "", fmt.Errorf("expected format: owner/repo")
	}

	if len(parts) > 2 {
		return "", "", fmt.Errorf("invalid URL: too many path segments")
	}

	return parts[0], parts[1], nil
}
