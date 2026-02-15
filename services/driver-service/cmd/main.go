package main

import (
	"context"
	"go-ride/services/driver-service/internal/infrastructure/grpc"
	"go-ride/services/driver-service/internal/repository" // Import do novo repositório
	"go-ride/services/driver-service/internal/service"
	"go-ride/shared/env"
	"go-ride/shared/messaging"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	grpcserver "google.golang.org/grpc"
)

var (
	GrpcAddr    = env.GetString("GRPC_ADDR", ":9092")
	AMQPAddr    = env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672")
	RedisAddr   = env.GetString("REDIS_ADDR", "localhost:6379")
	environment = env.GetString("ENVIRONMENT", "development")
)

func main() {
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

	log.Printf("Connecting to Redis at %s", RedisAddr)
	rdb := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
	})

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := rdb.Ping(pingCtx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer rdb.Close()
	log.Println("connected to Redis successfully")

	rabbitmq, err := messaging.NewRabbitMQ(AMQPAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()
	log.Println("starting rabbitmq connection")

	driverRepo := repository.NewRedisRepository(rdb)
	driverService := service.NewDriverService(driverRepo)

	grpcServer := grpcserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, driverService)

	log.Printf("starting GRPC driver service on port %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("shutting down the server...")

	grpcServer.GracefulStop()
}
