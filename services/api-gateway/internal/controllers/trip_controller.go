package controllers

import (
	"encoding/json"
	"go-ride/services/api-gateway/internal/dto"
	"go-ride/shared/contracts"
	"go-ride/shared/responses"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type TripController struct {
	validator *validator.Validate
}

func NewTripController(v *validator.Validate) *TripController {
	return &TripController{
		validator: v,
	}
}

func (s *TripController) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
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

	// fluxo feliz 🚀
	responses.WriteJSON(w, http.StatusOK, contracts.APIResponse{
		Data: "trip preview ok",
	})
}
