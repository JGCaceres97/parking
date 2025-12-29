package parking

import (
	"context"

	"github.com/JGCaceres97/parking/internal/domain"
)

type Service interface {
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

type Repository interface {
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
