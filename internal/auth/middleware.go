package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"log/slog"

	redisdb "github.com/Asilbeek1/mini-twitter-api/internal/repository/redis"
)

type contextKey string

const (
	ContextUserID contextKey = "user_id"
)

type Middleware struct {
	jwt        *JWTManager
	tokenStore *redisdb.TokenStore
	limiter    *redisdb.RateLimiter
	log        *slog.Logger
}

func NewMiddleware(
	jwt *JWTManager,
	tokenStore *redisdb.TokenStore,
	limiter *redisdb.RateLimiter,
	log *slog.Logger,
) *Middleware {
	return &Middleware{jwt: jwt, tokenStore: tokenStore, limiter: limiter, log: log}
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := ExtractBearer(r)
		if tokenStr == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwt.ParseAccessToken(tokenStr)
		if err != nil {
			if errors.Is(err, ErrExpiredToken) {
				http.Error(w, "token expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "invalid token", http.StatusUnauthorized)
			}
			return
		}

		blacklisted, err := m.tokenStore.IsBlacklisted(r.Context(), tokenStr)
		if err != nil {
			m.log.ErrorContext(r.Context(), "blacklist check failed", slog.Any("err", err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if blacklisted {
			http.Error(w, "token revoked", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		allowed, err := m.limiter.Allow(r.Context(), ip)
		if err != nil {
			m.log.ErrorContext(r.Context(), "rate limiter error", slog.Any("err", err))
			next.ServeHTTP(w, r)
			return
		}

		if !allowed {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ExtractBearer(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}

func UserIDFromCtx(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(ContextUserID).(int64)
	return id, ok
}
