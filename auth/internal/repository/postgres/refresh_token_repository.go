package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/radik-alt/auth_service/auth/internal/model"
	"github.com/radik-alt/auth_service/auth/internal/repository"
)

type refreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepo(db *pgxpool.Pool) repository.TokenRepository {
	return &refreshTokenRepository{db: db}
}

func (repo *refreshTokenRepository) SaveByTokenHash(ctx context.Context, token *model.RefreshToken) error {
	return nil
}

func (repo *refreshTokenRepository) GetByTokenHash(ctx context.Context, userID string) (string, error) {
	return "", nil
}

func (repo *refreshTokenRepository) DeleteByTokenHash(ctx context.Context, token string) error {
	return nil
}
