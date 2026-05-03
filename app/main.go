package main

import (
	"log/slog"

	"github.com/Asilbeek1/mini-twitter-api/internal/config"
	logger "github.com/Asilbeek1/mini-twitter-api/pkg"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()

	// TODO : init logger: slog
	log := logger.SetupLogger(cfg.Env)

	log.Info("MINI REDIS APP is starting", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled", slog.String("env", cfg.Env))

	//TODO: init database: postgres

	//TODO: init router: chi

}
