package ports

import (
	"context"
)

// InitService define el contrato para el manejo de los procesos de inicio.
type InitService interface {
	// CreateAdmin configura el primer usuario administrador del sistema.
	CreateAdmin(ctx context.Context, password string) error
}
