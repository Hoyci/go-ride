package types

import (
	pd "go-ride/shared/proto/driver"
	pu "go-ride/shared/proto/user"
)

type UserType string
type DriverStatus string

const (
	DRIVER    UserType = "DRIVER"
	PASSENGER UserType = "PASSENGER"
)

const (
	ONLINE  DriverStatus = "ONLINE"
	OFFLINE DriverStatus = "OFFLINE"
)

func MapUserTypeDomainToProto(t UserType) pu.UserType {
	switch t {
	case DRIVER:
		return pu.UserType_DRIVER
	case PASSENGER:
		return pu.UserType_PASSENGER
	default:
		return pu.UserType_USER_TYPE_UNSPECIFIED
	}
}

func MapProtoToUserTypeDomain(t pu.UserType) UserType {
	switch t {
	case pu.UserType_DRIVER:
		return DRIVER
	case pu.UserType_PASSENGER:
		return PASSENGER
	default:
		return PASSENGER
	}
}

func MapDriverStatusDomainToProto(s DriverStatus) pd.DriverStatusType {
	switch s {
	case ONLINE:
		return pd.DriverStatusType_ONLINE
	case OFFLINE:
		return pd.DriverStatusType_OFFLINE
	default:
		return pd.DriverStatusType_STATUS_TYPE_UNSPECIFIED
	}
}

func MapProtoDriverStatusToDomain(s pd.DriverStatusType) DriverStatus {
	switch s {
	case pd.DriverStatusType_ONLINE:
		return ONLINE
	case pd.DriverStatusType_OFFLINE:
		return OFFLINE
	default:
		return OFFLINE
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
