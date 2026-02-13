package main

import (
	"context"
	"go-ride/services/user-service/internal/infrastructure/grpc"
	"go-ride/services/user-service/internal/repository"
	"go-ride/services/user-service/internal/service"
	"go-ride/shared/env"
	"go-ride/shared/jwt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	grpcserver "google.golang.org/grpc"
)

var (
	GrpcAddr    = env.GetString("GRPC_ADDR", ":9091")
	environment = env.GetString("ENVIRONMENT", "development")
	JWTSecret   = env.GetString("JWT_SECRET", "um-secret-muito-complexo")
	RedisAddr   = env.GetString("REDIS_ADDR", "localhost:6379")
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
	})
	inmemRepo := repository.NewInmemRepository()
	userSvc := service.NewUserSerivce(inmemRepo)
	jwtSvc := jwt.NewJWTService(JWTSecret)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signCh := make(chan os.Signal, 1)
		signal.Notify(signCh, os.Interrupt, syscall.SIGTERM)
		<-signCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpcserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, userSvc, jwtSvc, rdb)

	go func() {
		log.Printf("starting GRPC user service on port %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	log.Println("shutting down the server...")
	grpcServer.GracefulStop()
}
