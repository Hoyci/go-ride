package main

import (
	"context"
	"go-ride/services/api-gateway/internal/controllers"
	httpHandler "go-ride/services/api-gateway/internal/handlers/http"
	"go-ride/services/api-gateway/internal/handlers/ws"
	"go-ride/shared/env"
	"go-ride/shared/jwt"
	"go-ride/shared/messaging"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_clients "go-ride/services/api-gateway/internal/clients/grpc"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

var (
	httpAddr      = env.GetString("HTTP_ADDR", ":8081")
	environment   = env.GetString("ENVIRONMENT", "development")
	userSvcAddr   = env.GetString("USER_SERVICE_ADDR", "user-service:9091")
	driverSvcAddr = env.GetString("DRIVER_SERVICE_ADDR", "driver-service:9092")
	tripSvcAddr   = env.GetString("TRIP_SERVICE_ADDR", "trip-service:9093")
	JWTSecret     = env.GetString("JWT_SECRET", "um-secret-muito-complexo")
	RedisAddr     = env.GetString("REDIS_ADDR", "localhost:6379")
)

func main() {
	log.Println("Starting API Gateway")

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

	tripClient, conn, err := grpc_clients.NewTripServiceClient(tripSvcAddr)
	if err != nil {
		log.Fatalf("could not connect to trip service: %v", err)
	}
	defer conn.Close()

	userClient, conn, err := grpc_clients.NewUserServiceClient(userSvcAddr)
	if err != nil {
		log.Fatalf("could not connect to user service: %v", err)
	}
	defer conn.Close()

	driverClient, conn, err := grpc_clients.NewDriverServiceClient(driverSvcAddr)
	if err != nil {
		log.Fatalf("could not connect to user service: %v", err)
	}
	defer conn.Close()

	connManager := messaging.NewConnectionManager()

	jwtSvc := jwt.NewJWTService(JWTSecret)
	v := validator.New()

	tripController := controllers.NewTripController(v, tripClient)
	userController := controllers.NewUserController(v, userClient)
	driverController := controllers.NewDriverController(v, driverClient)

	driverWSHandler := ws.NewDriverWSHandler(connManager, driverClient)

	handler := httpHandler.NewHTTPHandler(jwtSvc, rdb)
	handler.RegisterRoutes(userController, tripController, driverController, driverWSHandler)
	finalHandler := handler.GetHandler()

	server := &http.Server{
		Addr:    httpAddr,
		Handler: finalHandler,
	}

	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server listening on %s", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		log.Printf("Error starting the server: %v", err)

	case sig := <-shutdown:
		log.Printf("Server is shutting down due to %v signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop server gracefully: %v", err)
			server.Close()
		}
	}
}
