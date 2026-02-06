package grpc_clients

import (
	pb "go-ride/shared/proto/trip"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewTripServiceClient(addr string) (pb.TripServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	return pb.NewTripServiceClient(conn), conn, nil
}
