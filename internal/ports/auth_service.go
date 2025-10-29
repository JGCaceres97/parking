package ports

import (
	"context"

	"github.com/JGCaceres97/parking/internal/api/dtos"
	"github.com/JGCaceres97/parking/internal/core/domain"
)

// AuthService define el contrato para el manejo de la autenticación y JWT.
type AuthService interface {
	// Login verifica las credenciales y, si son válidas, genera un token JWT.
	// Retorna un LoginResponse que incluye el token y su expiración.
	Login(ctx context.Context, req dtos.LoginRequest) (*dtos.LoginResponse, error)

	// ParseToken verifica la validez del JWT y extrae los claims. Retorna el ID de usuario y el rol.
	ParseToken(tokenStr string) (userId string, role domain.Role, err error)
}
