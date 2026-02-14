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

	// TODO: Checar se a fare em questão já não expirou
	// Se tiver expirada, preciso estimar o novo preço e informar o usuário e ver se ele aprova.
	// Se ele aceitar, crio a viagem e emito o evento para buscar motorista
	// Se ele recusar, então devo somente enviar ele para a home

	trip := &domain.TripModel{
		ID:          uuid.New(),
		PassengerID: passengerID,
		Status:      "pending",
		RideFare:    fare,
		Driver:      &trip.TripDriver{},
	}

	return s.repo.CreateTrip(ctx, trip)
}

// TODO: Adicionar um endpoint que cancela uma viagem.
// Basicamente será pegar o ID do tripModel e atualizar o status para CANCELED
// PS: Não deve ser possível cancelar se já estiver em progresso
// PS: Nao deve ser possível cancelar se já estiver completa
// PS: Se for cancelada pelo passageiro e estiver com o status ACCEPTED, então devo avisar o motorista para que ele possa buscar novas viagens
// PS: Se for cancelada pelo motorista e estiver com o status ACCEPTED, então devo avisar o passageiro para que ele possa buscar novas viagens

// TODO: Adicionar um endpoint para encerrar uma viagem
// Basicamente será pegar o ID do tripModel e atualizar o status para COMPLETED
// PS: Só deve ser possível encerrar se já estiver em progresso
// PS: Nao deve ser possível encerrar se já estiver completa
// PS: Não deve ser possível encerrar se a corrida estiver com o status de ACCEPTED
// PS: Somente o motorista pode encerrar uma viagem
// PS: O usuário precisa ser notificado que a viagem dele foi encerrada
