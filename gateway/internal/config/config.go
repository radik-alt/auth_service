package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	HTTPAddr     string
	AuthGRPCAddr string
}

func Load(path string) (*AppConfig, error) {
	_ = godotenv.Load(path)

	cfg := &AppConfig{
		HTTPAddr:     envOr("HTTP_ADDR", ":8080"),
		AuthGRPCAddr: envOr("AUTH_GRPC_ADDR", "localhost:50051"),
	}

	return cfg, nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
