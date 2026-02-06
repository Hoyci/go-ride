package main

import (
	"context"
	"go-ride/services/api-gateway/internal/controllers"
	httpHandler "go-ride/services/api-gateway/internal/handlers/http"
	"go-ride/shared/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_clients "go-ride/services/api-gateway/internal/clients/grpc"

	"github.com/go-playground/validator/v10"
)

var (
	httpAddr    = env.GetString("HTTP_ADDR", ":8081")
	environment = env.GetString("ENVIRONMENT", "development")
	tripSvcAddr = env.GetString("TRIP_SERVICE_ADDR", "trip-service:9093")
)

func main() {
	log.Println("Starting API Gateway")

	v := validator.New()
	tripClient, conn, err := grpc_clients.NewTripServiceClient(tripSvcAddr)
	if err != nil {
		log.Fatalf("could not connect to trip service: %v", err)
	}
	defer conn.Close()

	tripController := controllers.NewTripController(v, tripClient)

	handler := httpHandler.NewHTTPHandler()
	handler.RegisterRoutes(tripController)
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
