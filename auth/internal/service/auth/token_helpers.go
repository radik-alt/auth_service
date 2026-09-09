package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

func newRefreshTokenValue() string {
	return uuid.NewString()
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
