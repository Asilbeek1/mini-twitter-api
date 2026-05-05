package handler

import (
	"log/slog"

	repo "github.com/Asilbeek1/mini-twitter-api/internal/repository"
)

type Handler struct {
	postRepo *repo.PostRepository
	userRepo *repo.UserRepository
	log      *slog.Logger
}

func New(log *slog.Logger) *Handler {
	return &Handler{log: log}
}
