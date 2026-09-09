package auth

import (
	"time"

	"github.com/radik-alt/auth_service/auth/internal/config"
	"github.com/radik-alt/auth_service/auth/internal/repository"
	"github.com/radik-alt/auth_service/auth/internal/service"
)

type authService struct {
	refreshRepo repository.TokenRepository
	blacklist   repository.TokenBlacklist
	jwtSecret   []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func NewAuthService(
	refreshRepo repository.TokenRepository,
	blacklist repository.TokenBlacklist,
	cfg *config.AppConfig,
) service.AuthService {
	return &authService{
		refreshRepo: refreshRepo,
		blacklist:   blacklist,
		jwtSecret:   []byte(cfg.JWTSecret),
		accessTTL:   cfg.AccessTTL,
		refreshTTL:  cfg.RefreshTTL,
	}
}
