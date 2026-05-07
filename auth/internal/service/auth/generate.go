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
	jti := uuid.NewString()
	now := time.Now()
	expiresAt := now.Add(s.accessTTL)

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"jti":   jti,
		"exp":   expiresAt.Unix(),
		"iat":   now.Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	rawRefresh := make([]byte, 32)
	if _, err := rand.Read(rawRefresh); err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}
	rawRefreshStr := base64.RawURLEncoding.EncodeToString(rawRefresh)
	tokenHash := hashToken(rawRefreshStr)

	refreshToken := &model.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(s.refreshTTL),
		CreatedAt: now,
	}
	if err := s.repo.SaveByTokenHash(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshStr,
		ExpiresIn:    int64(s.accessTTL.Seconds()),
	}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
