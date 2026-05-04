package redis

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed lua/ratelimit.lua
var limitScriptSource string

var limitScript = redis.NewScript(limitScriptSource)

type RedisLimiter struct {
	client *redis.Client
}

func NewRedisLimiter(client *redis.Client) *RedisLimiter {
	return &RedisLimiter{client: client}
}

func (rl *RedisLimiter) Allow(ctx context.Context, ip string, rps float64, burst int) (bool, error) {
	result, err := limitScript.Run(ctx, rl.client, []string{"ratelimit:" + ip},
		rps,
		burst,
		time.Now().UnixNano(),
	).Int()

	if err != nil {
		return true, err
	}

	return result == 1, nil
}
