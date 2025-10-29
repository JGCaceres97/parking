package middlewares

import (
	"net/http"

	"github.com/JGCaceres97/parking/internal/core/domain"
	"github.com/JGCaceres97/parking/pkg/response"
)

func RoleMiddleware(requiredRole domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extraer el rol
			ctx := r.Context()

			userRole, ok := ctx.Value(UserRoleKey).(string)
			if !ok {
				response.ErrorJSON(w, response.ErrUserIDNotInContext, http.StatusInternalServerError)
				return
			}

			// Verificar si el rol del usuario coincide con el requerido
			if domain.Role(userRole) != requiredRole {
				response.ErrorJSON(w, response.ErrPermissionDenied, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
