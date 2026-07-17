package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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
	const query = `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
	           VALUES ($1, $2, $3, $4, $5)`
	if _, err := repo.db.Exec(ctx, query, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.CreatedAt); err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}

func (repo *refreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	const query = `SELECT id, user_id, token_hash, expires_at, created_at
	           FROM refresh_tokens WHERE token_hash = $1`
	var t model.RefreshToken
	err := repo.db.QueryRow(ctx, query, tokenHash).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get refresh token: %w", err)
	}
	return &t, nil
}

func (repo *refreshTokenRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	const query = `DELETE FROM refresh_tokens WHERE token_hash = $1`
	if _, err := repo.db.Exec(ctx, query, tokenHash); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}
