package domain

import (
	"context"
	"time"
)

type Repository struct {
	FullName    string
	Description string
	Stargazers  int
	Forks       int
	CreatedAt   time.Time
	HTMLURL     string
}

type RepositoryProvider interface {
	GetRepository(ctx context.Context, owner, repo string) (*Repository, error)
}
