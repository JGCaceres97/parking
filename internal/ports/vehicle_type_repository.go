package ports

import (
	"context"

	"github.com/JGCaceres97/parking/internal/core/domain"
)

// VehicleTypeRepository define el contrato para la persistencia de la entidad VehicleType.
type VehicleTypeRepository interface {
	// FindByID obtiene la información de un tipo de vehículo por su ULID.
	// Esto es necesario para obtener la tarifa horario aplicada.
	FindByID(ctx context.Context, id string) (*domain.VehicleType, error)

	// ListAll obtiene una lista de todos los tipos de vehículo.
	ListAll(ctx context.Context) ([]domain.VehicleType, error)
}
