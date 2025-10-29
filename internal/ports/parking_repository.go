package ports

import (
	"context"

	"github.com/JGCaceres97/parking/internal/core/domain"
)

// ParkingRepository define el contrato para la persistencia de la entidad ParkingRecord.
type ParkingRepository interface {
	// CreateEntry registra la entrada de un vehículo.
	CreateEntry(ctx context.Context, record *domain.ParkingRecord) error

	// FindByID busca un registro de estacionamiento por su identificador.
	FindByID(ctx context.Context, id string) (*domain.ParkingRecord, error)

	// FindOpenByLicensePlate busca un registro de estacionamiento activo (exit_time IS NULL)
	// para una placa específica.
	FindOpenByLicensePlate(ctx context.Context, licensePlate string) (*domain.ParkingRecord, error)

	// UpdateExit completa un registro de estacionamiento al registra la salida y el cobro.
	UpdateExit(ctx context.Context, record *domain.ParkingRecord) error

	// ListCurrent lista todos los vehículos que aún están estacionados (exit_time IS NULL).
	ListCurrent(ctx context.Context) ([]domain.ParkingRecord, error)

	// ListHistory lista todos los registros de estacionamiento (historial).
	ListHistory(ctx context.Context) ([]domain.ParkingRecord, error)
}
