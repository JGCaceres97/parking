package vehicle_type

import (
	"context"

	"github.com/JGCaceres97/parking/internal/domain"
)

type Service interface {
	// FindByID obtiene un tipo de vehículo por su identificador.
	FindByID(ctx context.Context, id string) (*domain.VehicleType, error)

	// ListAll obtiene una lista de todos los tipos de vehículo disponibles.
	ListAll(ctx context.Context) ([]domain.VehicleType, error)
}

type Repository interface {
	// FindByID obtiene la información de un tipo de vehículo por su ULID.
	// Esto es necesario para obtener la tarifa horario aplicada.
	FindByID(ctx context.Context, id string) (*domain.VehicleType, error)

	// ListAll obtiene una lista de todos los tipos de vehículo.
	ListAll(ctx context.Context) ([]domain.VehicleType, error)
}
