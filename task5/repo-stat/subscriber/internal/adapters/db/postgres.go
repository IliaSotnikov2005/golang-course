package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/subscriber/internal/db"
	"github.com/IliaSotnikov2005/golang-course/task5/repo-stat/subscriber/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *PostgresRepository) Save(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	row, err := r.queries.CreateSubscription(ctx, db.CreateSubscriptionParams{
		Owner: sub.Owner,
		Repo:  sub.Repo,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrSubscriptionAlreadyExists
		}

		return nil, err
	}

	return &domain.Subscription{
		ID:        int(row.ID),
		Owner:     row.Owner,
		Repo:      row.Repo,
		CreatedAt: row.CreatedAt.Time,
	}, nil

}

func (r *PostgresRepository) Delete(ctx context.Context, owner, repo string) error {
	err := r.queries.DeleteSubscription(ctx, db.DeleteSubscriptionParams{
		Owner: owner,
		Repo:  repo,
	})

	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]domain.Subscription, error) {
	rows, err := r.queries.ListSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	subscriptions := make([]domain.Subscription, 0, len(rows))

	for _, row := range rows {
		subscriptions = append(subscriptions, domain.Subscription{
			ID:        int(row.ID),
			Owner:     row.Owner,
			Repo:      row.Repo,
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return subscriptions, nil
}
