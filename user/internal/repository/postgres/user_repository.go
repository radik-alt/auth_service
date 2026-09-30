package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/radik-alt/auth_service/user/internal/model"
	"github.com/radik-alt/auth_service/user/internal/repository"
)

const uniqueViolationCode = "23505"

var uniqueIndexErrors = map[string]error{
	"users_email_uniq":    repository.ErrEmailTaken,
	"users_username_uniq": repository.ErrUsernameTaken,
	"users_phone_uniq":    repository.ErrPhoneTaken,
}

const userColumns = `id, email, username, phone, avatar_url, bio, password_hash,
	status, email_verified, created_at, updated_at, deleted_at`

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepo{pool: pool}
}

func (r *userRepo) CreateUser(ctx context.Context, param model.CreateUserParams) (*model.User, error) {
	query := `
		INSERT INTO users (email, username, phone, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING ` + userColumns

	var u model.User
	err := r.pool.QueryRow(ctx, query, param.Email, param.Username, param.Phone, param.PasswordHash).Scan(
		&u.ID, &u.Email, &u.Username, &u.Phone, &u.AvatarURL, &u.Bio, &u.PasswordHash,
		&u.Status, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if mapped := mapUniqueViolation(err); mapped != nil {
			return nil, mapped
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return &u, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, userID string, params model.UpdateUserParams) (*model.User, error) {
	query := `
		UPDATE users
		SET username = COALESCE($1, username),
		    phone = CASE WHEN $2::text IS NULL THEN phone ELSE NULLIF($2::text, '') END,
		    avatar_url = COALESCE($3, avatar_url),
		    bio = COALESCE($4, bio),
		    updated_at = now()
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING ` + userColumns

	var u model.User
	err := r.pool.QueryRow(ctx, query, params.Username, params.Phone, params.AvatarURL, params.Bio, userID).Scan(
		&u.ID, &u.Email, &u.Username, &u.Phone, &u.AvatarURL, &u.Bio, &u.PasswordHash,
		&u.Status, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		if mapped := mapUniqueViolation(err); mapped != nil {
			return nil, mapped
		}

		return nil, fmt.Errorf("update user: %w", err)
	}

	return &u, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT ` + userColumns + `
		FROM users
		WHERE email = $1 AND deleted_at IS NULL`

	var u model.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.Username, &u.Phone, &u.AvatarURL, &u.Bio, &u.PasswordHash,
		&u.Status, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &u, nil
}

func mapUniqueViolation(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != uniqueViolationCode {
		return nil
	}
	return uniqueIndexErrors[pgErr.ConstraintName]
}
