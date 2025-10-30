package ports

import (
	"context"

	"github.com/JGCaceres97/parking/internal/core/domain"
)

// UserService define el contrato para la gestión de usuarios, aplicando reglas de rol.
type UserService interface {
	// -- Admin

	// Create crea un nuevo usuario.
	Create(ctx context.Context, user *domain.User) (*domain.User, error)

	// Update actualiza la información de un usuario específico.
	Update(ctx context.Context, id string, userUpdate *domain.User) (*domain.User, error)

	// ToggleActive bloquea o desbloquea (elimina lógicamente) a un usuario.
	ToggleActive(ctx context.Context, id string, isActive bool) (*domain.User, error)

	// Delete eliminar un usuario.
	Delete(ctx context.Context, id string) error

	// ListAll lista todos los usuarios.
	ListAll(ctx context.Context, id string) ([]domain.User, error)

	// -- Common

	// UpdateUsername permite a un usuario editar únicamente su propio username.
	UpdateUsername(ctx context.Context, id string, newUsername string) (*domain.User, error)
}
