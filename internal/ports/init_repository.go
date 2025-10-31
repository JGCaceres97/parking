package ports

import (
	"context"

	"github.com/JGCaceres97/parking/internal/core/domain"
)

// InitRepository define el contrato para la persistencia de los procesos de inicio.
type InitRepository interface {
	// CreateAdmin registra el primer usuario administrador del sistema.
	CreateAdmin(ctx context.Context, admin *domain.User) error

	// ExistsUsername busca si existe un nombre de usuario.
	ExistsUsername(ctx context.Context, username string) bool
}
