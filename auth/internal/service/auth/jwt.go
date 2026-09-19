package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid or expired token")
var ErrTokenExpired = errors.New("token is expired")

type accessClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *authService) generateAccessToken(userID uint64) (token string, jti string, expiresAt time.Time, err error) {
	jti = uuid.NewString()
	now := time.Now()
	expiresAt = now.Add(s.accessTTL)

	claims := accessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   strconv.FormatUint(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}
	return signed, jti, expiresAt, nil
}

func (s *authService) keyFunc(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, ErrInvalidToken
	}
	return s.jwtSecret, nil
}

func (s *authService) parseAccessToken(tokenStr string) (*accessClaims, error) {
	return s.parse(tokenStr)
}

// parseAccessTokenIgnoringExpiry проверяет подпись, но не срок: нужен Logout, чтобы узнать
// пользователя по уже истёкшему access-токену.
func (s *authService) parseAccessTokenIgnoringExpiry(tokenStr string) (*accessClaims, error) {
	return s.parse(tokenStr, jwt.WithoutClaimsValidation())
}

func (s *authService) parse(tokenStr string, opts ...jwt.ParserOption) (*accessClaims, error) {
	claims := &accessClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, s.keyFunc, opts...)
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
