package main

import (
	"os"

	"github.com/Asilbeek1/mini-twitter-api/internal/config"
	server "github.com/Asilbeek1/mini-twitter-api/internal/http"
	"github.com/Asilbeek1/mini-twitter-api/internal/storage/postgres"
	"github.com/Asilbeek1/mini-twitter-api/internal/storage/redis"
	logger "github.com/Asilbeek1/mini-twitter-api/pkg"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	// Load config
	cfg := config.LoadConfig()
	// Setup logger
	log := logger.SetupLogger(cfg.Env)
	// Setup databases
	db, err := postgres.New(*cfg)
	if err != nil {
		log.Error("postgres db init failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	// Setup Redis
	rdb, err := redis.NewClient(*cfg)
	if err != nil {
		log.Error("redis db init failed", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()

	// Start server
	srv := server.New(cfg, db, rdb, log)

	log.Info("starting on", "port", cfg.HTTPServer.Address)
	srv.ListenAndServe()
}
