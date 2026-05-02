package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string         `yaml:"environment" env-default:"local"` // local dev prod
	Postgres   PostgresConfig `yaml:"postgres"`
	HTTPServer HTTPServer     `yaml:"http_server"`
}

type PostgresConfig struct {
	Host   string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port   int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	DBName string `yaml:"db_name" env:"DB_NAME" env-default:"postgres"`

	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`

	Pool PoolConfig `yaml:"pool"`
}

type HTTPServer struct {
	Address      string        `yaml:"address" env:"HTTP_ADDRESS" env-default:":8080"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"4s"`
}

type PoolConfig struct {
	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"25"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"4"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"60m"`
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG FILE does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}
	return &cfg
}
