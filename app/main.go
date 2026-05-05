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

	cfg := config.LoadConfig()


	log := logger.SetupLogger(cfg.Env)

	
	db, err := storage.New(cfg.SqLite.StoragePath)
	if err != nil {
		log.Error("db init failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	srv := server.New(cfg, db, log)

	log.Info("starting on", "port", cfg.HTTPServer.Address)
	srv.ListenAndServe()
}
