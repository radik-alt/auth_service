package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/radik-alt/auth_service/auth/internal/model"
)

func (s *authService) GenerateTokens(ctx context.Context, userID string, email string) (*model.TokenPair, error) {
	accessToken, err := s.issueAccessToken(userID, email)
	if err != nil {
		return nil, err
	}

	rawRefresh, err := s.issueRefreshToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    int64(s.accessTTL.Seconds()),
	}, nil
}

func (s *authService) issueAccessToken(userID, email string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"jti": uuid.NewString(),
		"exp": now.Add(s.accessTTL).Unix(),
		"iat": now.Unix(),
	}
	if email != "" {
		claims["email"] = email
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}
	return accessToken, nil
}

func (s *authService) issueRefreshToken(ctx context.Context, userID string) (string, error) {
	rawRefresh := make([]byte, 32)
	if _, err := rand.Read(rawRefresh); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	rawRefreshStr := base64.RawURLEncoding.EncodeToString(rawRefresh)
	now := time.Now()

	refreshToken := &model.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    userID,
		TokenHash: hashToken(rawRefreshStr),
		ExpiresAt: now.Add(s.refreshTTL),
		CreatedAt: now,
	}
	if err := s.repo.SaveByTokenHash(ctx, refreshToken); err != nil {
		return "", fmt.Errorf("save refresh token: %w", err)
	}

	return rawRefreshStr, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
