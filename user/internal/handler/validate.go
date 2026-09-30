package handler

import (
	"net/mail"
	"net/url"
	"regexp"
	"unicode/utf8"
)

const (
	minPasswordLen  = 8
	maxPasswordLen  = 72 // bcrypt учитывает только первые 72 байта
	maxEmailLen     = 254
	maxAvatarURLLen = 2048
	maxBioLen       = 500 // в символах, не в байтах
)

var (
	usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_.]{3,32}$`)
	phoneRe    = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`) // E.164
	uuidRe     = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

// validateUserID нужен до похода в БД: невалидный UUID Postgres вернёт как ошибку запроса (Internal), а не NotFound
func validateUserID(userID string) error {
	if userID == "" {
		return requiredField("user_id")
	}
	if !uuidRe.MatchString(userID) {
		return invalidField("user_id", "must be a UUID")
	}
	return nil
}

// validateAvatarURL допускает пустую строку — это удаление аватара
func validateAvatarURL(avatarURL string) error {
	if avatarURL == "" {
		return nil
	}
	if len(avatarURL) > maxAvatarURLLen {
		return invalidField("avatar_url", "too long")
	}
	u, err := url.Parse(avatarURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return invalidField("avatar_url", "must be an absolute http(s) URL")
	}
	return nil
}

func validateBio(bio string) error {
	if utf8.RuneCountInString(bio) > maxBioLen {
		return invalidField("bio", "must be at most 500 characters")
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return requiredField("email")
	}
	if len(email) > maxEmailLen {
		return invalidField("email", "too long")
	}
	// ParseAddress принимает и "Name <a@b.c>", поэтому сверяем, что это голый адрес
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return invalidField("email", "invalid format")
	}
	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return requiredField("username")
	}
	if !usernameRe.MatchString(username) {
		return invalidField("username", "must be 3-32 characters: letters, digits, '_' or '.'")
	}
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return requiredField("password")
	}
	if len(password) < minPasswordLen {
		return invalidField("password", "must be at least 8 characters")
	}
	if len(password) > maxPasswordLen {
		return invalidField("password", "must be at most 72 bytes")
	}
	return nil
}

func validatePhone(phone *string) error {
	if phone == nil {
		return nil
	}
	if !phoneRe.MatchString(*phone) {
		return invalidField("phone", "must be in E.164 format, e.g. +79991234567")
	}
	return nil
}
