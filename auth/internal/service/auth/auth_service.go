package auth

import (
	"github.com/radik-alt/auth_service/auth/internal/config"
	"github.com/radik-alt/auth_service/auth/internal/repository"
	"github.com/radik-alt/auth_service/auth/internal/service"
	"time"
)

type authService struct {
	repo       repository.TokenRepository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(
	repo repository.TokenRepository,
	cfg *config.AppConfig,
) service.AuthService {
	return &authService{
		repo:       repo,
		jwtSecret:  []byte(cfg.JWTSecret),
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
	}
}
