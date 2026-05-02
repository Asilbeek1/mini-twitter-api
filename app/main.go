package main

import (
	"fmt"

	"github.com/Asilbeek1/mini-twitter-api/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()
	fmt.Println(cfg)
}
