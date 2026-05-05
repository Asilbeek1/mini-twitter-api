package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Asilbeek1/mini-twitter-api/internal/auth"
	redisdb "github.com/Asilbeek1/mini-twitter-api/internal/repository/redis"
	"github.com/Asilbeek1/mini-twitter-api/internal/service"
)

type AuthHandler struct {
	userService *service.UserService
	jwt         *auth.JWTManager
	tokenStore  *redisdb.TokenStore
	log         *slog.Logger
}

func NewAuthHandler(userService *service.UserService, jwt *auth.JWTManager, tokenStore *redisdb.TokenStore, log *slog.Logger) *AuthHandler {
	return &AuthHandler{userService: userService, jwt: jwt, tokenStore: tokenStore, log: log}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.Login(input.Email, input.Password)
	if errors.Is(err, service.ErrInvalidCredentials) {
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	pair, err := h.jwt.GenerateTokenPair(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := h.tokenStore.SaveRefreshToken(r.Context(), user.ID, pair.RefreshToken, h.jwt.RefreshTTL()); err != nil {
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	respondJSON(w, http.StatusOK, pair)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claims, err := h.jwt.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	stored, err := h.tokenStore.GetRefreshToken(r.Context(), claims.UserID)
	if err != nil || stored != input.RefreshToken {
		respondError(w, http.StatusUnauthorized, "refresh token mismatch")
		return
	}

	pair, err := h.jwt.GenerateTokenPair(claims.UserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.tokenStore.SaveRefreshToken(r.Context(), claims.UserID, pair.RefreshToken, h.jwt.RefreshTTL())
	respondJSON(w, http.StatusOK, pair)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	tokenStr := auth.ExtractBearer(r)
	userID, ok := auth.UserIDFromCtx(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	h.tokenStore.BlacklistToken(r.Context(), tokenStr, h.jwt.AccessTTL())
	h.tokenStore.RevokeRefreshToken(r.Context(), userID)

	w.WriteHeader(http.StatusNoContent)
}
