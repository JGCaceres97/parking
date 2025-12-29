package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/JGCaceres97/parking/internal/application/auth"
	"github.com/JGCaceres97/parking/pkg/response"
)

type ContextKey string

const (
	UserIDKey   ContextKey = "userID"
	UserRoleKey ContextKey = "userRole"
)

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return "", response.ErrUserIDNotInContext
	}

	return userID, nil
}

func AuthMiddleware(service auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extraer token de la cabecera
			header := r.Header.Get("Authorization")
			if header == "" {
				response.ErrorJSON(w, response.ErrMissingToken, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(header, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.ErrorJSON(w, response.ErrInvalidTokenFormat, http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Validar el token
			userID, role, err := service.ParseToken(token)
			if err != nil {
				if errors.Is(err, auth.ErrExpiredToken) {
					response.ErrorJSON(w, response.ErrTokenExpired, http.StatusUnauthorized)
					return
				}

				if errors.Is(err, auth.ErrInvalidToken) {
					response.ErrorJSON(w, response.ErrInvalidToken, http.StatusUnauthorized)
					return
				}

				response.ErrorJSON(w, response.ErrTokenValidationFailed, http.StatusInternalServerError)
				return
			}

			// Inyectar la información del usuario
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, role)

			// Pasar el control
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
