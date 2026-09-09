package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	GRPCAddr   string
	RedisAddr  string
	JWTSecret  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

// Option — функция, модифицирующая уже собранный AppConfig
type Option func(*AppConfig)

func WithGRPCAddr(addr string) Option {
	return func(c *AppConfig) {
		c.GRPCAddr = addr
	}
}

func WithRedisAddr(addr string) Option {
	return func(c *AppConfig) {
		c.RedisAddr = addr
	}
}

func WithJWTSecret(secret string) Option {
	return func(c *AppConfig) {
		c.JWTSecret = secret
	}
}

func WithAccessTTL(d time.Duration) Option {
	return func(c *AppConfig) {
		c.AccessTTL = d
	}
}

func WithRefreshTTL(d time.Duration) Option {
	return func(c *AppConfig) {
		c.RefreshTTL = d
	}
}

func Load(path string, opts ...Option) (*AppConfig, error) {
	_ = godotenv.Load(path)

	accessTTL, err := parseDuration("ACCESS_TOKEN_TTL", "15m")
	if err != nil {
		return nil, fmt.Errorf("ACCESS_TOKEN_TTL: %w", err)
	}

	refreshTTL, err := parseDuration("REFRESH_TOKEN_TTL", "1h")
	if err != nil {
		return nil, fmt.Errorf("REFRESH_TOKEN_TTL: %w", err)
	}

	cfg := &AppConfig{
		GRPCAddr:   envOr("GRPC_ADDR", ":50051"),
		RedisAddr:  envOr("REDIS_ADDR", "localhost:6379"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		AccessTTL:  accessTTL,
		RefreshTTL: refreshTTL,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.RedisAddr == "" {
		return nil, fmt.Errorf("REDIS_ADDR is required")
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
