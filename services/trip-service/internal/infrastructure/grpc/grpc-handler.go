package grpc

import (
	"context"
	"go-ride/services/trip-service/internal/domain"
	"go-ride/services/trip-service/internal/events"
	pb "go-ride/shared/proto/trip"
	"go-ride/shared/types"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	tripService domain.TripService
	OSRMService domain.OSRMService
	publisher   *events.TripEventPublisher
}

func NewGRPCHandler(server *grpc.Server, tripService domain.TripService, OSRMService domain.OSRMService, publisher *events.TripEventPublisher) *gRPCHandler {
	handler := &gRPCHandler{
		tripService: tripService,
		OSRMService: OSRMService,
		publisher:   publisher,
	}

	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetStartLocation()
	destination := req.GetEndLocation()

	pickupCoord := &types.Coordinate{
		Latitude:  pickup.Latitude,
		Longitude: pickup.Longitude,
	}

	destinationCoord := &types.Coordinate{
		Latitude:  destination.Latitude,
		Longitude: destination.Longitude,
	}

	route, err := h.OSRMService.GetRoute(ctx, pickupCoord, destinationCoord)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	estimatedFares := h.tripService.EstimatePackagesPriceWithRoute(route)
	fares, err := h.tripService.GenerateTripFares(ctx, estimatedFares, req.PassengerID, route)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to generate the ride fares: %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     route.ToProto(),
		RideFares: domain.ToRideFaresProto(fares),
	}, nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	userID := req.GetUserID()
	fareID := req.GetRideFareID()

	rideFare, err := h.tripService.GetAndValidateFare(ctx, fareID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate the fare: %v", err)
	}

	trip, err := h.tripService.CreateTrip(ctx, rideFare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create trip: %v", err)
	}

	if err := h.publisher.PublishTripCreated(ctx, trip); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish the trip created event message: %v", err)
	}

	return &pb.CreateTripResponse{
		TripID: trip.ID.String(),
	}, nil
}
