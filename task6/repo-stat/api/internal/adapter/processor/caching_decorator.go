package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/internal/domain"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

type CachingDecorator struct {
	base  *Client
	cache Cache
	log   *slog.Logger
	ttl   time.Duration
}

func NewCachingDecorator(base *Client, cache Cache, log *slog.Logger, ttl time.Duration) *CachingDecorator {
	return &CachingDecorator{
		base:  base,
		cache: cache,
		log:   log,
		ttl:   ttl,
	}
}

func (cd *CachingDecorator) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	cacheKey := fmt.Sprintf("repo:%s:%s", owner, repo)

	if data, err := cd.cache.Get(ctx, cacheKey); err == nil {
		var r domain.Repository
		if err := json.Unmarshal(data, &r); err == nil {
			cd.log.Debug("cache hit", slog.String("key", cacheKey))
			return &r, nil
		}
	}

	cd.log.Debug("cache miss", slog.String("key", cacheKey))

	repoData, err := cd.base.GetRepository(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	go func() {
		bytes, err := json.Marshal(repoData)
		if err != nil {
			cd.log.Error("failed to marshal repo data for cache", slog.Any("error", err))
			return
		}

		writeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()

		if err := cd.cache.Set(writeCtx, cacheKey, bytes, cd.ttl); err != nil {
			cd.log.Error("failed to save repo to cache", slog.Any("error", err))
		}
	}()

	return repoData, nil
}

func (cd *CachingDecorator) GetSubscriptionsInfo(ctx context.Context) ([]domain.Repository, error) {
	cacheKey := "subscriptions:info:all"

	if data, err := cd.cache.Get(ctx, cacheKey); err == nil {
		var repos []domain.Repository
		if err := json.Unmarshal(data, &repos); err == nil {
			cd.log.Debug("cache hit for subscriptions info", slog.String("key", cacheKey))
			return repos, nil
		}
	}

	cd.log.Debug("cache miss for subscriptions info", slog.String("key", cacheKey))

	repos, err := cd.base.GetSubscriptionsInfo(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		bytes, err := json.Marshal(repos)
		if err != nil {
			cd.log.Error("failed to marshal repos for cache", slog.Any("error", err))
			return
		}

		writeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()

		if err := cd.cache.Set(writeCtx, cacheKey, bytes, cd.ttl); err != nil {
			cd.log.Error("failed to save subs info to cache", slog.Any("error", err))
		}
	}()

	return repos, nil
}
