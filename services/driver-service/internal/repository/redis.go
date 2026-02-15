package repository

import (
	"context"
	"fmt"
	"time"

	"go-ride/services/driver-service/internal/domain"
	"go-ride/shared/types"

	"github.com/redis/go-redis/v9"
)

const (
	driversLocationKey = "drivers:locations"
	driversStatusKey   = "drivers:status"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) domain.DriverRepository {
	return &redisRepository{
		client: client,
	}
}

func (r *redisRepository) SetStatus(ctx context.Context, driverID string, status types.DriverStatus) error {
	key := fmt.Sprintf("%s:%s", driversStatusKey, driverID)

	return r.client.Set(ctx, key, string(status), 30*time.Second).Err()
}

func (r *redisRepository) RemoveStatus(ctx context.Context, driverID string) error {
	key := fmt.Sprintf("%s:%s", driversStatusKey, driverID)
	return r.client.Del(ctx, key).Err()
}

func (r *redisRepository) UpdateLocation(ctx context.Context, driverID string, location *types.Coordinate) error {
	if location == nil {
		return nil
	}

	return r.client.GeoAdd(ctx, driversLocationKey, &redis.GeoLocation{
		Name:      driverID,
		Longitude: location.Longitude,
		Latitude:  location.Latitude,
	}).Err()
}

func (r *redisRepository) RemoveLocation(ctx context.Context, driverID string) error {
	return r.client.ZRem(ctx, driversLocationKey, driverID).Err()
}
