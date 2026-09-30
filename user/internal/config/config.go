package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	GRPCAddr   string
	HTTPAddr   string
	DB         string
	RedisAddr  string
	JWTSecret  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func Load(path string) (*AppConfig, error) {
	_ = godotenv.Load(path)

	accessTTL, err := parseDuration("ACCESS_TOKEN_TTL", "15m")
	if err != nil {
		return nil, fmt.Errorf("ACCESS_TOKEN_TTL: %w", err)
	}
	refreshTTL, err := parseDuration("REFRESH_TOKEN_TTL", "720h")
	if err != nil {
		return nil, fmt.Errorf("REFRESH_TOKEN_TTL: %w", err)
	}

	cfg := &AppConfig{
		GRPCAddr:   envOr("GRPC_ADDR", ":50051"),
		HTTPAddr:   envOr("HTTP_ADDR", ":8080"),
		DB:         os.Getenv("DATABASE_URL"),
		RedisAddr:  envOr("REDIS_ADDR", "localhost:6380"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		AccessTTL:  accessTTL,
		RefreshTTL: refreshTTL,
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.DB == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseDuration(key, def string) (time.Duration, error) {
	v := envOr(key, def)
	return time.ParseDuration(v)
}
