package model

import "time"

type UserStatus string

const (
	UserStatusActive  UserStatus = "active"
	UserStatusBlocked UserStatus = "blocked"
	UserStatusDeleted UserStatus = "deleted"
)

// User — DB-модель пользователя (таблица users)
type User struct {
	ID            string // UUID
	Email         string // хранится в нижнем регистре
	Username      string
	Phone         *string // E.164, nil — не указан
	AvatarURL     string
	Bio           string
	PasswordHash  string // bcrypt; наружу никогда не отдаётся
	Status        UserStatus
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time // заполнен при мягком удалении
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// CreateUserInput — данные регистрации из хендлера в сервис (пароль в открытом виде)
type CreateUserInput struct {
	Email    string
	Username string
	Phone    *string
	Password string
}

// CreateUserParams — данные для регистрации (пароль уже захеширован в сервисе)
type CreateUserParams struct {
	Email        string
	Username     string
	Phone        *string
	PasswordHash string
}

// UpdateUserParams — частичное обновление профиля: nil означает «не менять».
// Собирается в хендлере из update_mask.
type UpdateUserParams struct {
	Username  *string
	Phone     *string // указатель на "" — удалить телефон (в БД NULL)
	AvatarURL *string
	Bio       *string
}

// ListUsersParams — курсорная пагинация по (created_at, id)
type ListUsersParams struct {
	Limit int
	// курсор: последний элемент предыдущей страницы, nil — первая страница
	AfterCreatedAt *time.Time
	AfterID        string
}
