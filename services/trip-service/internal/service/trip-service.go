package service

import (
	"context"
	"fmt"
	"go-ride/services/trip-service/internal/domain"
	tripTypes "go-ride/services/trip-service/pkg/types"

	"github.com/google/uuid"
)

type tripService struct {
	repo domain.TripRepository
}

func NewTripService(repo domain.TripRepository) *tripService {
	return &tripService{repo: repo}
}

func (s *tripService) EstimatePackagesPriceWithRoute(route *tripTypes.OSRMApiResponse) []*domain.RideFareModel {
	baseFares := getBaseFares()
	estimatedFares := make([]*domain.RideFareModel, len(baseFares))

	for i, fare := range baseFares {
		estimatedFares[i] = s.estimateFareRoute(fare, route)
	}

	return estimatedFares
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

		fare = &domain.RideFareModel{
			ID:                id,
			PassengerID:       passengerID,
			TotalPriceInCents: fare.TotalPriceInCents,
			PackageSlug:       fare.PackageSlug,
			Route:             route,
		}

		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("failed to save trip fare: %v", err)
		}

		fares[i] = fare
	}

	return fares, nil
}

func (s *tripService) estimateFareRoute(fare *domain.RideFareModel, route *tripTypes.OSRMApiResponse) *domain.RideFareModel {
	pricingCfg := tripTypes.DefaultPricingConfig()
	carPackagePrice := fare.TotalPriceInCents

	distanceInKm := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration

	distanceFare := distanceInKm * pricingCfg.PricePerUnitOfDistance
	timeFare := durationInMinutes * pricingCfg.PricePerUnitOfTime

	totalPrice := carPackagePrice + distanceFare + timeFare

	return &domain.RideFareModel{
		PackageSlug:       fare.PackageSlug,
		TotalPriceInCents: totalPrice,
	}
}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug:       domain.UBERX,
			TotalPriceInCents: 200,
		},
		{
			PackageSlug:       domain.BLACK,
			TotalPriceInCents: 350,
		},
	}
}
