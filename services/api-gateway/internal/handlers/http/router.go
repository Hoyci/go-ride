package http

import (
	"go-ride/services/api-gateway/internal/controllers"
	"net/http"
)

type Handler struct {
	Router *http.ServeMux
}

func NewHTTPHandler() *Handler {
	router := http.NewServeMux()

	return &Handler{
		Router: router,
	}
}

func (h *Handler) RegisterRoutes(tripController *controllers.TripController) {
	h.registerTripRoutes(tripController)
}

func (h *Handler) registerTripRoutes(tripController *controllers.TripController) {
	h.Router.HandleFunc("POST /api/v1/trip-preview", tripController.HandleTripPreview)
}

func (h *Handler) GetHandler() http.Handler {
	finalHandler := Chain(
		Logger,
		Recoverer,
		CORS,
	)(h.Router)

	return finalHandler
}
