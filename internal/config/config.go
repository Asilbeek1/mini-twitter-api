package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string          `yaml:"environment" env-default:"local"` // local dev prod
	Postgres   PostgresConfig  `yaml:"postgres"`
	HTTPServer HTTPServer      `yaml:"http_server"`
	Redis      RedisConfig     `yaml:"redis"`
	JWT        JWTConfig       `yaml:"jwt"`
	RateLimit  RateLimitConfig `yaml:"rate_limit"`
}

type PostgresConfig struct {
	Host       string     `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost" env-required:"true"`
	Port       string     `yaml:"port" env:"POSTGRES_PORT" env-default:"5432" env-required:"true"`
	User       string     `env:"DB_USER" env-required:"true"`
	Password   string     `env:"DB_PASSWORD" env-required:"true"`
	Name       string     `env:"DB_NAME" env-required:"true"`
	PoolConfig PoolConfig `yaml:"pool"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" env:"REDIS_ADDR" env-default:"localhost:6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
	DB       int    `yaml:"db" env:"REDIS_DB" env-default:"0"`
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

type JWTConfig struct {
	AccessSecret  string        `yaml:"access_secret"`
	RefreshSecret string        `yaml:"refresh_secret"`
	AccessTTL     time.Duration `yaml:"access_ttl"`
	RefreshTTL    time.Duration `yaml:"refresh_ttl"`
}
type RateLimitConfig struct {
	Limit  int           `yaml:"limit" env-default:"10"`
	Window time.Duration `yaml:"window" env-default:"1m"`
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
