package types

import (
	pu "go-ride/shared/proto/user"
)

type UserType string

const (
	DRIVER    UserType = "DRIVER"
	PASSENGER UserType = "PASSENGER"
)

func MapUserTypeToProto(t UserType) pu.UserType {
	switch t {
	case DRIVER:
		return pu.UserType_DRIVER
	case PASSENGER:
		return pu.UserType_PASSENGER
	default:
		return pu.UserType_PASSENGER
	}
}

type Route struct {
	Distance float64     `json:"distance"`
	Duration float64     `json:"duration"`
	Geometry []*Geometry `json:"geometry"`
}

type Geometry struct {
	Coordinates []*Coordinate `json:"coordinates"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}
