package handler

import (
	"log/slog"

	repo "github.com/Asilbeek1/mini-twitter-api/internal/repository"
)

type Handler struct {
	repo *repo.Repository
	log  *slog.Logger
}

func New(repo *repo.Repository, log *slog.Logger) *Handler {
	return &Handler{repo: repo, log: log}
}
