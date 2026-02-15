package domain

import (
	"context"
	"go-ride/shared/types"
)

type DriverService interface {
	UpdateDriverStatus(ctx context.Context, driverID string, status types.DriverStatus, location *types.Coordinate) error
}

type DriverRepository interface {
	SetStatus(ctx context.Context, driverID string, status types.DriverStatus) error
	RemoveStatus(ctx context.Context, driverID string) error

	UpdateLocation(ctx context.Context, driverID string, location *types.Coordinate) error
	RemoveLocation(ctx context.Context, driverID string) error
}
