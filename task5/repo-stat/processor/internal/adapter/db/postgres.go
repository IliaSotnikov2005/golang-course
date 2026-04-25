package db

import (
	"context"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/processor/internal/db"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/processor/internal/domain"
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
		return nil, fmt.Errorf("failed to get repoistory: %w", err)
	}

}
