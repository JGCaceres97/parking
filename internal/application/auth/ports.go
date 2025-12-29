package auth

import (
	"context"

	"github.com/JGCaceres97/parking/internal/domain"
)

type LoginInput struct {
	Username string
	Password string
}

type LoginOutput struct {
	Token     string
	TokenType string
	ExpiresIn int64
	Role      domain.Role
}

type Service interface {
	// CreateAdmin configura el primer usuario administrador del sistema.
	CreateAdmin(ctx context.Context, password string) error

	// Login verifica las credenciales y, si son válidas, genera un token JWT.
	// Retorna un LoginResponse que incluye el token y su expiración.
	Login(ctx context.Context, req LoginInput) (*LoginOutput, error)

	// ParseToken verifica la validez del JWT y extrae los claims. Retorna el ID de usuario y el rol.
	ParseToken(tokenStr string) (userId string, role domain.Role, err error)
}
