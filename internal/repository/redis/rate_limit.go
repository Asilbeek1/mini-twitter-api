package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/Asilbeek1/mini-twitter-api/internal/config"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRateLimiter(client *redis.Client, cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{client: client, limit: cfg.Limit, window: cfg.Window}
}

func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	redisKey := fmt.Sprintf("rate:%s", key)

	pipe := rl.client.Pipeline()
	incr := pipe.Incr(ctx, redisKey)
	pipe.Expire(ctx, redisKey, rl.window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("rate limiter pipeline: %w", err)
	}

	return incr.Val() <= int64(rl.limit), nil
}

func (rl *RateLimiter) Remaining(ctx context.Context, key string) (int, error) {
	redisKey := fmt.Sprintf("rate:%s", key)
	count, err := rl.client.Get(ctx, redisKey).Int()
	if err == redis.Nil {
		return rl.limit, nil
	}
	if err != nil {
		return 0, err
	}
	remaining := rl.limit - count
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}
