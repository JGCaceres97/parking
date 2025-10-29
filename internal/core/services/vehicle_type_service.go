package services

import (
	"context"

	"github.com/JGCaceres97/parking/internal/core/domain"
	"github.com/JGCaceres97/parking/internal/ports"
)

type VehicleTypeService struct {
	repo ports.VehicleTypeRepository
}

func NewVehicleTypeService(repo ports.VehicleTypeRepository) ports.VehicleTypeService {
	return &VehicleTypeService{repo: repo}
}

func (s *VehicleTypeService) FindByID(ctx context.Context, id string) (*domain.VehicleType, error) {
	vehicleType, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return vehicleType, nil
}

func (s *VehicleTypeService) ListAll(ctx context.Context) ([]domain.VehicleType, error) {
	vehicleTypes, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return vehicleTypes, nil
}
