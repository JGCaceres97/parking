package ports

import (
	"context"

	"github.com/JGCaceres97/parking/internal/core/domain"
)

// UserService define el contrato para la gestión de usuarios, aplicando reglas de rol.
type UserService interface {
	// -- Admin

	// CreateUser crea un nuevo usuario. Solo para admins.
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)

	// UpdateUser actualiza la información de un usuario específico. Solo para admins.
	UpdateUser(ctx context.Context, id string, userUpdate *domain.User) (*domain.User, error)

	// ToggleActive bloquea o desbloquea (elimina lógicamente) a un usuario. Solo para admins.
	ToggleActive(ctx context.Context, id string, isActive bool) (*domain.User, error)

	// ListUsers lista todos los usuarios.
	ListUsers(ctx context.Context) ([]domain.User, error)

	// -- Common

	// UpdateUsername permite a un usuario editar únicamente su propio username.
	UpdateUsername(ctx context.Context, id string, newUsername string) (*domain.User, error)
}
