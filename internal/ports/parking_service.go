package ports

import (
	"context"

	"github.com/JGCaceres97/parking/internal/core/domain"
)

// ParkingService define el contrato para el control del estacionamiento y el cálculo de tarifas.
type ParkingService interface {
	// RecordEntry registra la entrada de un vehículo.
	RecordEntry(ctx context.Context, userID string, vehicleTypeID string, licensePlate string) (*domain.ParkingRecord, error)

	// RecordExit registra la salida del vehículo, calcula el tiempo y el cobro.
	RecordExit(ctx context.Context, userID string, licensePlate string) (*domain.ParkingRecord, error)

	// GetCurrentlyParked lista todos los vehículos que tienen registro de entrada abierto.
	GetCurrentlyParked(ctx context.Context) ([]domain.ParkingRecord, error)

	// GetHistory lista todos los registros, incluyendo los cerrados.
	GetHistory(ctx context.Context) ([]domain.ParkingRecord, error)

	// GetRecordByID obtiene un registro específico.
	GetRecordByID(ctx context.Context, id string) (*domain.ParkingRecord, error)
}
