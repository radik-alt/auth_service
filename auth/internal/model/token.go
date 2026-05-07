package model

import "time"

// TokenPair возвращается из сервиса в хендлер
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // секунды до истечения access-токена
}

// RefreshToken — DB-модель хранения refresh-токенов
type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string // SHA-256 от сырого токена
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (t *RefreshToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
