package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/radik-alt/auth_service/gateway/internal/client"
	"google.golang.org/grpc/codes"
)

// TokenValidator — проверка access-токена (реализует client.AuthClient).
type TokenValidator interface {
	ValidateToken(ctx context.Context, accessToken string) (userID string, err error)
}

type userIDKey struct{}

// UserIDFromContext возвращает user_id, проверенный middleware Auth.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey{}).(string)
	return userID, ok && userID != ""
}

// Auth пропускает публичные пути как есть, а для остальных требует
// заголовок "Authorization: Bearer <access>" и кладёт user_id владельца токена в контекст.
// Закрыто всё, что не перечислено в isPublic: новая ручка по умолчанию требует токен.
func Auth(validator TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// "/api/v1/auth/../users/x" и подобные не должны проскочить мимо проверки по префиксу
			if path.Clean(r.URL.Path) != r.URL.Path {
				writeError(w, http.StatusNotFound, codes.NotFound, "not found")
				return
			}

			if isPublic(r) {
				next.ServeHTTP(w, r)
				return
			}

			token, ok := bearerToken(r)
			if !ok {
				writeError(w, http.StatusUnauthorized, codes.Unauthenticated, "missing access token")
				return
			}

			userID, err := validator.ValidateToken(r.Context(), token)
			switch {
			case err == nil:
			case errors.Is(err, client.ErrInvalidToken):
				writeError(w, http.StatusUnauthorized, codes.Unauthenticated, "invalid or expired token")
				return
			case errors.Is(err, client.ErrUnavailable):
				log.Printf("auth middleware: %v", err)
				writeError(w, http.StatusServiceUnavailable, codes.Unavailable, "service temporarily unavailable")
				return
			default:
				log.Printf("auth middleware: %v", err)
				writeError(w, http.StatusInternalServerError, codes.Internal, "internal error")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isPublic(r *http.Request) bool {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/users":
		return true
	case strings.HasPrefix(r.URL.Path, "/api/v1/auth/"):
		return true
	}
	return false
}

func bearerToken(r *http.Request) (string, bool) {
	scheme, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}
	token = strings.TrimSpace(token)
	return token, token != ""
}

// writeError отвечает в том же формате, что и ошибки grpc-gateway
func writeError(w http.ResponseWriter, httpStatus int, code codes.Code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":    code,
		"message": message,
		"details": []any{},
	})
}
