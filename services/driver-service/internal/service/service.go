package service

import (
	"context"
	"go-ride/services/driver-service/internal/domain"
	"go-ride/shared/types"
)

type DriverService struct {
	repo domain.DriverRepository
}

func NewDriverService(repo domain.DriverRepository) *DriverService {
	return &DriverService{
		repo: repo,
	}
}

func (s *DriverService) UpdateDriverStatus(ctx context.Context, driverID string, status types.DriverStatus, location *types.Coordinate) error {
	if status == types.OFFLINE {
		s.repo.RemoveLocation(ctx, driverID)
		return s.repo.RemoveStatus(ctx, driverID)
	}

	if err := s.repo.SetStatus(ctx, driverID, status); err != nil {
		return err
	}

	if location != nil {
		if err := s.repo.UpdateLocation(ctx, driverID, location); err != nil {
			return err
		}
	}

	return nil
}
