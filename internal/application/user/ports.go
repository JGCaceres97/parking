package user

import (
	"context"

	"github.com/JGCaceres97/parking/internal/domain"
)

type Service interface {
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

type Repository interface {
	// Create registra un nuevo usuario en la base de datos.
	Create(ctx context.Context, user *domain.User) error

	// FindByID busca un usuario por su ULID
	FindByID(ctx context.Context, id string) (*domain.User, error)

	// FindByUsername busca un usuario por su nombre de usuario para el login.
	FindByUsername(ctx context.Context, username string) (*domain.User, error)

	// ExistsUsername busca si existe un nombre de usuario.
	ExistsUsername(ctx context.Context, username string) bool

	// Update actualiza la información del usuario.
	Update(ctx context.Context, user *domain.User) error

	// Delete eliminar un usuario.
	Delete(ctx context.Context, id string) error

	// ListAll lista todos los usuarios, excepto a ti mismo.
	ListAll(ctx context.Context, id string) ([]domain.User, error)
}
