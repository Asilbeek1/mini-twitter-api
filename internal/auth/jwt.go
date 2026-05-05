package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/Asilbeek1/mini-twitter-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type JWTManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWTManager(cfg config.Config) *JWTManager {
	return &JWTManager{
		accessSecret:  []byte(cfg.JWT.AccessSecret),
		refreshSecret: []byte(cfg.JWT.RefreshSecret),
		accessTTL:     cfg.JWT.AccessTTL,
		refreshTTL:    cfg.JWT.RefreshTTL,
	}
}

func (m *JWTManager) GenerateTokenPair(userID int64) (*TokenPair, error) {
	accessToken, err := m.generate(userID, m.accessSecret, m.accessTTL)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}
	refreshToken, err := m.generate(userID, m.refreshSecret, m.refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (m *JWTManager) ParseRefreshToken(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, m.refreshSecret)
}

func (m *JWTManager) ParseAccessToken(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, m.accessSecret)
}

func (m *JWTManager) generate(userID int64, secret []byte, ttl time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (m *JWTManager) parse(tokenStr string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (m *JWTManager) AccessTTL() time.Duration  { return m.accessTTL }
func (m *JWTManager) RefreshTTL() time.Duration { return m.refreshTTL }
