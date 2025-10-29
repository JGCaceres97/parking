package ports

import (
	"context"

	"github.com/JGCaceres97/parking/internal/core/domain"
)

// UserRepository define el contrato para la persistencia de la entidad User.
type UserRepository interface {
	// Create registra un nuevo usuario en la base de datos.
	Create(ctx context.Context, user *domain.User) error

	// FindByID busca un usuario por su ULID
	FindByID(ctx context.Context, id string) (*domain.User, error)

	// FindByUsername busca un usuario por su nombre de usuario para el login.
	FindByUsername(ctx context.Context, username string) (*domain.User, error)

	// Update actualiza la información del usuario.
	Update(ctx context.Context, user *domain.User) error

	// List lista todos los usuarios.
	List(ctx context.Context) ([]domain.User, error)
}
