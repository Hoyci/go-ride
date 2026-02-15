package grpc_clients

import (
	pd "go-ride/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewDriverServiceClient(addr string) (pd.DriverServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	return pd.NewDriverServiceClient(conn), conn, nil
}
