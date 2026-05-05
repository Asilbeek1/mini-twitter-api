package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenStore struct {
	client *redis.Client
}

func NewTokenStore(client *redis.Client) *TokenStore {
	return &TokenStore{client: client}
}

func (s *TokenStore) BlacklistToken(ctx context.Context, tokenStr string, ttl time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", tokenStr)
	return s.client.Set(ctx, key, "1", ttl).Err()
}

func (s *TokenStore) IsBlacklisted(ctx context.Context, tokenStr string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", tokenStr)
	err := s.client.Get(ctx, key).Err()
	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("Error checking blacklist: %w", err)
	}

	return true, nil
}

func (s *TokenStore) SaveRefreshToken(ctx context.Context, userID int64, token string, ttl time.Duration) error {
	key := fmt.Sprintf("refresh:%d", userID)
	return s.client.Set(ctx, key, token, ttl).Err()
}

func (s *TokenStore) GetRefreshToken(ctx context.Context, userID int64) (string, error) {
	key := fmt.Sprintf("refresh:%d", userID)
	token, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("refresh token not found or expired")
	}
	if err != nil {
		return "", fmt.Errorf("get refresh token: %w", err)
	}
	return token, nil
}

func (s *TokenStore) RevokeRefreshToken(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("refresh:%d", userID)
	return s.client.Del(ctx, key).Err()
}
