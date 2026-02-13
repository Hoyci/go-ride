package controllers

import (
	"context"
	"encoding/json"
	"go-ride/services/api-gateway/internal/dto"
	"go-ride/shared/contracts"
	"go-ride/shared/responses"
	"log"
	"net/http"
	"time"

	pb "go-ride/shared/proto/trip"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type TripController struct {
	validator   *validator.Validate
	tripService pb.TripServiceClient
}

func NewTripController(v *validator.Validate, ts pb.TripServiceClient) *TripController {
	return &TripController{
		validator:   v,
		tripService: ts,
	}
}

func (s *TripController) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req dto.PreviewTripRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.WriteJSON(w, http.StatusBadRequest, contracts.APIResponse{
			Error: &contracts.APIError{
				Code:    http.StatusBadRequest,
				Message: "invalid JSON payload",
			},
		})
		return
	}

	if err := s.validator.Struct(req); err != nil {
		responses.WriteJSON(w, http.StatusUnprocessableEntity, contracts.APIResponse{
			Error: &contracts.APIError{
				Code:    http.StatusUnprocessableEntity,
				Message: "validation failed",
				Details: responses.ParseValidationErrors(err),
			},
		})
		return
	}

	// Call GRPC trip_service method to get trip preview
	grpcRes, err := s.tripService.PreviewTrip(ctx, &pb.PreviewTripRequest{
		PassengerID: req.PassengerID,
		StartLocation: &pb.Coordinate{
			Latitude:  req.Origin.Latitude,
			Longitude: req.Origin.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:  req.Destination.Latitude,
			Longitude: req.Destination.Longitude,
		},
	}, grpc.WaitForReady(true))

	if err != nil {
		log.Printf("failed to call preview trip: %v", err)
		responses.WriteJSON(w, http.StatusInternalServerError, contracts.APIResponse{
			Error: &contracts.APIError{Message: "failed to contact trip service"},
		})
		return
	}

	responses.WriteJSON(w, http.StatusOK, contracts.APIResponse{
		Data: grpcRes,
	})
}

func (s *TripController) HandleCreateTrip(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req dto.CreateTripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.WriteJSON(w, http.StatusBadRequest, contracts.APIResponse{
			Error: &contracts.APIError{
				Code:    http.StatusBadRequest,
				Message: "invalid JSON payload",
			},
		})
		return
	}

	if err := s.validator.Struct(req); err != nil {
		responses.WriteJSON(w, http.StatusUnprocessableEntity, contracts.APIResponse{
			Error: &contracts.APIError{
				Code:    http.StatusUnprocessableEntity,
				Message: "validation failed",
				Details: responses.ParseValidationErrors(err),
			},
		})
		return
	}

	grpcRes, err := s.tripService.CreateTrip(ctx, &pb.CreateTripRequest{
		RideFareID: req.RideFareID,
		UserID:     req.UserID,
	}, grpc.WaitForReady(true))

	if err != nil {
		log.Printf("failed to call create trip: %v", err)
		responses.WriteJSON(w, http.StatusInternalServerError, contracts.APIResponse{
			Error: &contracts.APIError{Message: "failed to contact trip service"},
		})
		return
	}

	responses.WriteJSON(w, http.StatusCreated, contracts.APIResponse{
		Data: grpcRes,
	})
}
