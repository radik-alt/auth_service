package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	GRPCAddr string
	DB       string
}

func Load(path string) (*AppConfig, error) {
	_ = godotenv.Load(path)

	cfg := &AppConfig{
		GRPCAddr: envOr("GRPC_ADDR", ":50051"),
		DB:       os.Getenv("DATABASE_URL"),
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
