package main

import (
	"fmt"
	"os"

	"github.com/Asilbeek1/mini-twitter-api/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if err := config.LoadConfig(); err != nil {
		os.Exit(1)
	}
	fmt.Println("Config was succesfully loaded")

}
