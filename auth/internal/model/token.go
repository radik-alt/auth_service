package model

import "time"

// TokenPair возвращается из сервиса в хендлер
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    uint64 // секунды до истечения access-токена
}

// RefreshToken — DB-модель хранения refresh-токенов
type RefreshToken struct {
	ID        string
	UserID    string // UUID из user-сервиса
	TokenHash string // SHA-256 от сырого токена
	ExpiresAt time.Time
	CreatedAt time.Time

	// парный access-токен: нужен, чтобы отозвать его при refresh и LogoutAll
	AccessJTI       string
	AccessExpiresAt time.Time
}

func (t *RefreshToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
