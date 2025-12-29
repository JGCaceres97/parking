package vehicle_type

import (
	"context"

	"github.com/JGCaceres97/parking/internal/domain"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) FindByID(ctx context.Context, id string) (*domain.VehicleType, error) {
	vehicleType, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return vehicleType, nil
}

func (s *service) ListAll(ctx context.Context) ([]domain.VehicleType, error) {
	vehicleTypes, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return vehicleTypes, nil
}
