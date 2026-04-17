package middleware

import (
	"context"
	"ethno/internal/auth"
	"net/http"

	"github.com/sirupsen/logrus"
)

type contextKey string
const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	UsernameKey  contextKey = "username"
	UserRoleKey  contextKey = "role"
)

func AuthFromCookie(jwtProvider auth.Provider, logger *logrus.Logger, cookieName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cookieName)
			if err != nil || cookie.Value == "" {
				logger.Warn("No auth cookie found")
				http.Error(w, `{"error":"Не авторизован"}`, http.StatusUnauthorized)
				return
			}

			claims, err := jwtProvider.ParseJWT(cookie.Value)
			if err != nil {
				logger.Warnf("Invalid JWT: %v", err)
				http.Error(w, `{"error":"Неверный токен"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, UsernameKey, claims.Username)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) string {
	if v := r.Context().Value(UserIDKey); v != nil {
		return v.(string)
	}
	return ""
}

func GetUsername(r *http.Request) string {
	if v := r.Context().Value(UsernameKey); v != nil {
		return v.(string)
	}
	return ""
}

func GetUserEmail(r *http.Request) string {
	if v := r.Context().Value(UserEmailKey); v != nil {
		return v.(string)
	}
	return ""
}

func GetUserRole(r *http.Request) string {
	if v := r.Context().Value(UserRoleKey); v != nil {
		return v.(string)
	}
	return ""
}
