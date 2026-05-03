package main

import (
	"os"

	"github.com/Asilbeek1/mini-twitter-api/internal/config"
	server "github.com/Asilbeek1/mini-twitter-api/internal/http"
	"github.com/Asilbeek1/mini-twitter-api/internal/storage"
	logger "github.com/Asilbeek1/mini-twitter-api/pkg"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// 1. Config first — everything depends on it
	cfg := config.LoadConfig()

	// 2. Logger second — needs config (log level, format, etc.)
	log := logger.SetupLogger(cfg.Env)

	// 3. DB third — needs config (DSN), runs migrations
	db, err := storage.New(cfg.SqLite.StoragePath)
	if err != nil {
		log.Error("db init failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// 4. Server last — needs everything
	srv := server.New(cfg, db, log)

	log.Info("starting on", "port", cfg.HTTPServer.Address)
	srv.ListenAndServe()
}
