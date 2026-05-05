package server

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/Asilbeek1/mini-twitter-api/internal/auth"
	"github.com/Asilbeek1/mini-twitter-api/internal/config"
	handler "github.com/Asilbeek1/mini-twitter-api/internal/http/handlers"
	postgresdb "github.com/Asilbeek1/mini-twitter-api/internal/repository/postgres"
	redisdb "github.com/Asilbeek1/mini-twitter-api/internal/repository/redis"
	"github.com/Asilbeek1/mini-twitter-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

func New(cfg *config.Config, db *sql.DB, rdb *redis.Client, log *slog.Logger) *http.Server {
	userRepo := postgresdb.NewUserRepository(db)
	postRepo := postgresdb.NewPostRepository(db)

	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo)

	jwtManager := auth.NewJWTManager(*cfg)
	tokenStore := redisdb.NewTokenStore(rdb)
	rateLimiter := redisdb.NewRateLimiter(rdb, cfg.RateLimit)

	authHandler := handler.NewAuthHandler(userService, jwtManager, tokenStore, log)
	userHandler := handler.NewUserHandler(userService, log)
	postHandler := handler.NewPostHandler(postService, log)

	mw := auth.NewMiddleware(jwtManager, tokenStore, rateLimiter, log)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.RateLimit)
			r.Post("/register", userHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.RateLimit)
			r.Use(mw.Authenticate)
			r.Post("/logout", authHandler.Logout)
			r.Post("/refresh", authHandler.Refresh)
			r.Get("/users/{id}", userHandler.GetProfile)
			r.Put("/users/{id}", userHandler.UpdateProfile)
			r.Delete("/users/{id}", userHandler.DeleteProfile)
			r.Get("/posts/feed", postHandler.GetFeed)
			r.Post("/posts", postHandler.CreatePost)
			r.Get("/posts/{id}", postHandler.GetPost)
			r.Put("/posts/{id}", postHandler.UpdatePost)
			r.Delete("/posts/{id}", postHandler.DeletePost)
			r.Get("/users/{id}/posts", postHandler.GetUserPosts)
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
