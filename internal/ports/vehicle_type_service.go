package ports

import (
	"context"

	"github.com/JGCaceres97/parking/internal/core/domain"
)

// VehicleTypeService define el contrato para la gestión de tipos de vehículo.
type VehicleTypeService interface {
	// FindByID obtiene un tipo de vehículo por su identificador.
	FindByID(ctx context.Context, id string) (*domain.VehicleType, error)

	// ListAll obtiene una lista de todos los tipos de vehículo disponibles.
	ListAll(ctx context.Context) ([]domain.VehicleType, error)
}
