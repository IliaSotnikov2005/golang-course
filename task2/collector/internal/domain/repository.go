package domain

import (
	"context"
	"time"
)

type Repository struct {
	Name        string
	Description string
	Stargazers  int
	Forks       int
	CreatedAt   time.Time
	HTMLURL     string
}

type GitHubClient interface {
	GetRepository(ctx context.Context, owner, repo string) (*Repository, error)
}

type RepositoryProvider interface {
	GetRepository(ctx context.Context, owner, repo string) (*Repository, error)
}
