package server

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/Asilbeek1/mini-twitter-api/internal/config"
	handler "github.com/Asilbeek1/mini-twitter-api/internal/http/handlers"
	repo "github.com/Asilbeek1/mini-twitter-api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(cfg *config.Config, db *sql.DB, log *slog.Logger) *http.Server {
	r := chi.NewRouter()

	// --- Global middleware ---
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// --- Wire dependencies ---
	repo := repo.New(db)
	_ = handler.New(repo, log)

	// --- Routes ---
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("Ok"))
		})
	})

	return &http.Server{
		Addr:         ":" + cfg.HTTPServer.Address,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
