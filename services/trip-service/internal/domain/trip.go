package domain

import (
	"context"
	tripTypes "go-ride/services/trip-service/pkg/types"
	"go-ride/shared/types"

	"github.com/google/uuid"
)

type TripStatus string

const (
	REQUESTED   TripStatus = "REQUESTED"
	ACCEPTED    TripStatus = "ACCEPTED"
	IN_PROGRESS TripStatus = "IN_PROGRESS"
	COMPLETED   TripStatus = "COMPLETED"
	CANCELED    TripStatus = "CANCELED"
)

type TripModel struct {
	ID          uuid.UUID
	PassengerID uuid.UUID
	Status      TripStatus
	RideFare    *RideFareModel
}

// func (t *TripModel) ToProto() *pb.Trip {
// 	return &pb.Trip{
// 		Id:           t.ID.Hex(),
// 		UserId:       t.UserID,
// 		SelectedFare: t.RideFare.ToProto(),
// 		Status:       t.Status,
// 		Driver:       t.Driver,
// 		Route:        t.RideFare.Route.ToProto(),
// 	}
// }

type TripRepository interface {
	// CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, fare *RideFareModel) error
	// GetRideFareByID(ctx context.Context, fareID string) (*RideFareModel, error)
	// GetTripByID(ctx context.Context, tripID string) (*TripModel, error)
	// UpdateTrip(ctx context.Context, tripID, status string, driver *pbd.Driver) error
}

type TripService interface {
	// CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	EstimatePackagesPriceWithRoute(route *tripTypes.OSRMApiResponse) []*RideFareModel
	GenerateTripFares(ctx context.Context, fares []*RideFareModel, userID string, route *tripTypes.OSRMApiResponse) ([]*RideFareModel, error)
	// GetAndValidateFare(ctx context.Context, fareID, userID string) (*RideFareModel, error)
	// GetTripByID(ctx context.Context, tripID string) (*TripModel, error)
	// UpdateTrip(ctx context.Context, tripID, status string, driver *pbd.Driver) error
}

type OSRMService interface {
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OSRMApiResponse, error)
}
