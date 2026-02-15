package grpc

import (
	"context"
	"log"

	"go-ride/services/driver-service/internal/service"
	pd "go-ride/shared/proto/driver"
	"go-ride/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pd.UnimplementedDriverServiceServer
	server        *grpc.Server
	driverService *service.DriverService
}

func NewGRPCHandler(server *grpc.Server, driverService *service.DriverService) *gRPCHandler {
	handler := &gRPCHandler{
		server:        server,
		driverService: driverService,
	}

	pd.RegisterDriverServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) UpdateStatus(ctx context.Context, req *pd.UpdateStatusRequest) (*pd.UpdateStatusResponse, error) {
	driverStatus := types.MapProtoDriverStatusToDomain(req.Status)

	var location *types.Coordinate
	if req.ActualLocation != nil {
		location = &types.Coordinate{
			Latitude:  req.ActualLocation.Latitude,
			Longitude: req.ActualLocation.Longitude,
		}
	}

	err := h.driverService.UpdateDriverStatus(ctx, req.DriverID, driverStatus, location)
	if err != nil {
		log.Printf("Failed to update driver status: %v", err)
		return nil, status.Error(codes.Internal, "failed to update status")
	}

	return &pd.UpdateStatusResponse{
		Success: true,
	}, nil
}
