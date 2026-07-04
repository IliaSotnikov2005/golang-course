package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/processor/internal/db"
	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/processor/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PostgresRepository) GetByFullName(ctx context.Context, fullName string) (*domain.Repository, error) {
	row, err := r.queries.GetRepository(ctx, fullName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return &domain.Repository{
		FullName:    row.FullName,
		Description: row.Description.String,
		Stargazers:  int(row.Stargazers),
		Forks:       int(row.Forks),
		CreatedAt:   row.CreatedAt.Time,
		HTMLURL:     row.HtmlUrl.String,
	}, nil
}

func (r *PostgresRepository) Upsert(ctx context.Context, repo *domain.Repository) error {
	return r.queries.UpsertRepository(ctx, db.UpsertRepositoryParams{
		FullName: repo.FullName,
		Description: pgtype.Text{
			String: repo.Description,
			Valid:  repo.Description != "",
		},
		Stargazers: int32(repo.Stargazers),
		Forks:      int32(repo.Forks),
		CreatedAt: pgtype.Timestamptz{
			Time:  repo.CreatedAt,
			Valid: !repo.CreatedAt.IsZero(),
		},
		HtmlUrl: pgtype.Text{
			String: repo.HTMLURL,
			Valid:  repo.HTMLURL != "",
		},
	})
}

func (r *PostgresRepository) ListAll(ctx context.Context) ([]domain.Repository, error) {
	rows, err := r.queries.ListAllRepositories(ctx)
	if err != nil {
		return nil, err
	}

	repositories := make([]domain.Repository, 0, len(rows))
	for _, repo := range rows {
		repositories = append(repositories, domain.Repository{
			FullName:    repo.FullName,
			Description: repo.Description.String,
			Stargazers:  int(repo.Stargazers),
			Forks:       int(repo.Forks),
			CreatedAt:   repo.CreatedAt.Time,
			HTMLURL:     repo.HtmlUrl.String,
		})
	}

	return repositories, nil
}
