package grpc_clients

import (
	pu "go-ride/shared/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserServiceClient(addr string) (pu.UserServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	return pu.NewUserServiceClient(conn), conn, nil
}
