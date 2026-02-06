package dto

import "go-ride/shared/types"

type PreviewTripRequest struct {
	PassengerID string            `json:"passenger_id" validate:"required,uuid4"`
	Origin      *types.Coordinate `json:"origin" validate:"required"`
	Destination *types.Coordinate `json:"destination" validate:"required"`
}
