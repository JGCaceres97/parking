package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/JGCaceres97/parking/internal/core/domain"
	"github.com/JGCaceres97/parking/internal/ports"
	"github.com/JGCaceres97/parking/pkg/ulid"
)

type ParkingService struct {
	repo        ports.ParkingRepository
	vehicleRepo ports.VehicleTypeRepository
}

func NewParkingService(repo ports.ParkingRepository, vehicleRepo ports.VehicleTypeRepository) ports.ParkingService {
	return &ParkingService{
		repo:        repo,
		vehicleRepo: vehicleRepo,
	}
}

func (s *ParkingService) RecordEntry(ctx context.Context, userID, vehicleTypeID, licensePlate string) (*domain.ParkingRecord, error) {
	// Verificar si ya existe registro abierto para la placa.
	_, err := s.repo.FindOpenByLicensePlate(ctx, licensePlate)
	if err == nil {
		return nil, ports.ErrActiveParkingExists
	}

	if !errors.Is(err, ports.ErrParkingRecordNotFound) {
		return nil, fmt.Errorf("error al verificar registro abierto: %w", err)
	}

	// Verificar que el tipo de vehículo sea válido.
	_, err = s.vehicleRepo.FindByID(ctx, vehicleTypeID)
	if err != nil {
		if errors.Is(err, ports.ErrVehicleTypeNotFound) {
			return nil, ports.ErrVehicleTypeNotFound
		}

		return nil, fmt.Errorf("error al buscar tipo de vehículo: %w", err)
	}

	record := domain.ParkingRecord{
		ID:            ulid.GenerateNewULID(),
		UserID:        userID,
		VehicleTypeID: vehicleTypeID,
		LicensePlate:  licensePlate,
		EntryTime:     time.Now().UTC().Truncate(time.Second),
	}

	if err = s.repo.CreateEntry(ctx, &record); err != nil {
		return nil, fmt.Errorf("error al guardar registro de entrada: %w", err)
	}

	return &record, nil
}

func (s *ParkingService) RecordExit(ctx context.Context, userID, licensePlate string) (*domain.ParkingRecord, error) {
	record, err := s.repo.FindOpenByLicensePlate(ctx, licensePlate)
	if err != nil {
		if errors.Is(err, ports.ErrParkingRecordNotFound) {
			return nil, ports.ErrActiveParkingNotFound
		}

		return nil, fmt.Errorf("error al buscar registro abierto: %w", err)
	}

	vehicleType, err := s.vehicleRepo.FindByID(ctx, record.VehicleTypeID)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la tarifa: %w", err)
	}

	exitTime := time.Now().UTC()
	hours, charge := calculateCharge(record.EntryTime, exitTime, vehicleType.HourlyRate)

	truncatedExitTime := exitTime.Truncate(time.Second)

	record.ExitTime = &truncatedExitTime
	record.TotalCharge = &charge
	record.CalculatedHours = &hours

	if err = s.repo.UpdateExit(ctx, record); err != nil {
		return nil, fmt.Errorf("error al actualizar registro de salida: %w", err)
	}

	return record, nil
}

func (s *ParkingService) GetCurrentlyParked(ctx context.Context) ([]domain.ParkingRecord, error) {
	return s.repo.ListCurrent(ctx)
}

func (s *ParkingService) GetHistory(ctx context.Context) ([]domain.ParkingRecord, error) {
	return s.repo.ListHistory(ctx)
}

func (s *ParkingService) GetRecordByID(ctx context.Context, id string) (*domain.ParkingRecord, error) {
	record, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func calculateCharge(entryTime, exitTime time.Time, hourlyRate float64) (int, float64) {
	if hourlyRate == 0.00 {
		duration := exitTime.Sub(entryTime)

		calculatedHours := int(math.Ceil(duration.Hours()))
		if calculatedHours == 0 {
			calculatedHours = 1 // Mínimo 1 hora
		}

		return calculatedHours, 0.00
	}

	minutes := exitTime.Sub(entryTime).Minutes()

	if minutes <= 0 {
		return 1, hourlyRate // Mínimo 1 hora
	}

	fullHours := int(minutes / 60.0)
	remainingMinutes := minutes - (float64(fullHours) * 60.0)

	calculatedHours := fullHours
	if remainingMinutes >= 30.0 {
		calculatedHours += 1
	}

	if calculatedHours == 0 {
		calculatedHours = 1
	}

	totalCharge := float64(calculatedHours) * hourlyRate

	return calculatedHours, totalCharge
}
