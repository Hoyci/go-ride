package service

import (
	"context"
	"fmt"
	"go-ride/services/trip-service/internal/domain"
	tripTypes "go-ride/services/trip-service/pkg/types"
	"go-ride/shared/proto/trip"

	"github.com/google/uuid"
)

type tripService struct {
	repo domain.TripRepository
}

func NewTripService(repo domain.TripRepository) *tripService {
	return &tripService{repo: repo}
}

func (s *tripService) EstimatePackagesPriceWithRoute(route *tripTypes.OSRMApiResponse) []*domain.RideFareModel {
	categories := []domain.PackageSlug{domain.UBERX, domain.BLACK}
	estimatedFares := make([]*domain.RideFareModel, len(categories))

	for i, slug := range categories {
		estimatedFares[i] = s.calculateFareBySlug(slug, route)
	}

	return estimatedFares
}

func (s *tripService) calculateFareBySlug(slug domain.PackageSlug, route *tripTypes.OSRMApiResponse) *domain.RideFareModel {
	cfg := s.getPricingConfig(slug)

	// OSRM: metros -> km | segundos -> minutos
	distanceInKm := route.Routes[0].Distance / 1000.0
	durationInMinutes := route.Routes[0].Duration / 60.0

	distanceTotal := distanceInKm * float64(cfg.PricePerUnitOfDistance)
	timeTotal := durationInMinutes * float64(cfg.PricePerUnitOfTime)

	totalPrice := float64(cfg.BaseFare) + distanceTotal + timeTotal

	if totalPrice < float64(cfg.MinimumFare) {
		totalPrice = float64(cfg.MinimumFare)
	}

	return &domain.RideFareModel{
		PackageSlug:       slug,
		TotalPriceInCents: totalPrice,
	}
}

func (s *tripService) getPricingConfig(slug domain.PackageSlug) tripTypes.PricingConfig {
	switch slug {
	case domain.BLACK:
		return tripTypes.PricingConfig{
			BaseFare:               500,
			PricePerUnitOfDistance: 250,
			PricePerUnitOfTime:     60,
			MinimumFare:            1500,
		}
	default: // UBERX
		return tripTypes.PricingConfig{
			BaseFare:               350,
			PricePerUnitOfDistance: 160,
			PricePerUnitOfTime:     30,
			MinimumFare:            800,
		}
	}
}

func (s *tripService) GenerateTripFares(
	ctx context.Context,
	rideFares []*domain.RideFareModel,
	passengerID string,
	route *tripTypes.OSRMApiResponse,
) ([]*domain.RideFareModel, error) {
	fares := make([]*domain.RideFareModel, len(rideFares))

	for i, fare := range rideFares {
		id := uuid.New()
		newFare := &domain.RideFareModel{
			ID:                id,
			PassengerID:       passengerID,
			TotalPriceInCents: fare.TotalPriceInCents,
			PackageSlug:       fare.PackageSlug,
			Route:             route,
		}

		if err := s.repo.SaveRideFare(ctx, newFare); err != nil {
			return nil, fmt.Errorf("failed to save trip fare: %v", err)
		}
		fares[i] = newFare
	}
	return fares, nil
}

func (s *tripService) GetAndValidateFare(ctx context.Context, fareID, userID string) (*domain.RideFareModel, error) {
	fare, err := s.repo.GetRideFareByID(ctx, fareID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip fare: %v", err)
	}

	if fare == nil {
		return nil, fmt.Errorf("fare does not exists")
	}

	if fare.PassengerID != userID {
		return nil, fmt.Errorf("fare %s does not belong to the user %s", fare.ID, userID)
	}

	return fare, nil
}

func (s *tripService) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	passengerID, err := uuid.Parse(fare.PassengerID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse userID as uuid")
	}

	trip := &domain.TripModel{
		ID:          uuid.New(),
		PassengerID: passengerID,
		Status:      "pending",
		RideFare:    fare,
		Driver:      &trip.TripDriver{},
	}

	return s.repo.CreateTrip(ctx, trip)
}
