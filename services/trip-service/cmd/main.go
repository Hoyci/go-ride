package main

import (
	"context"
	"go-ride/services/trip-service/internal/events"
	"go-ride/services/trip-service/internal/infrastructure/grpc"
	"go-ride/services/trip-service/internal/repository"
	"go-ride/services/trip-service/internal/service"
	"go-ride/shared/env"
	"go-ride/shared/messaging"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var (
	GrpcAddr    = env.GetString("GRPC_ADDR", ":9093")
	AMQPAddr    = env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672")
	environment = env.GetString("ENVIRONMENT", "development")
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	osrmSvc := service.NewOSRMService()
	tripSvc := service.NewTripService(inmemRepo)

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

	// Starting the rabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(AMQPAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	log.Println("starting rabbitmq connection")

	publisher := events.NewTripEventPublisher(rabbitmq)

	grpcServer := grpcserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, tripSvc, osrmSvc, publisher)

	go func() {
		log.Printf("starting GRPC trip service on port %s", lis.Addr().String())
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
