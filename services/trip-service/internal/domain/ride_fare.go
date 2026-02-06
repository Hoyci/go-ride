package domain

import (
	"time"

	pb "go-ride/shared/proto/trip"

	tripTypes "go-ride/services/trip-service/pkg/types"

	"github.com/google/uuid"
)

type PackageSlug string

const (
	UberX PackageSlug = "UBERX"
	Black PackageSlug = "BLACK"
)

type RideFareModel struct {
	ID                uuid.UUID
	PassengerID       string
	PackageSlug       PackageSlug
	TotalPriceInCents float64
	ExpiresAt         time.Time
	Route             *tripTypes.OSRMApiResponse
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:                r.ID.String(),
		PassengerID:       r.PassengerID,
		PackageSlug:       toProtoPackageSlug(r.PackageSlug),
		TotalPriceInCents: r.TotalPriceInCents,
	}
}

func toProtoPackageSlug(s PackageSlug) pb.PackageSlug {
	switch s {
	case UberX:
		return pb.PackageSlug_UBERX
	case Black:
		return pb.PackageSlug_BLACK
	default:
		return pb.PackageSlug_PACKAGE_SLUG_UNSPECIFIED
	}
}

func ToRideFaresProto(fares []*RideFareModel) []*pb.RideFare {
	var protoFares []*pb.RideFare
	for _, fare := range fares {
		protoFares = append(protoFares, fare.ToProto())
	}

	return protoFares
}
