package http

import (
	"go-ride/services/api-gateway/internal/controllers"
	"go-ride/shared/jwt"
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

func (h *Handler) RegisterRoutes(
	userController *controllers.UserController,
	tripController *controllers.TripController,
	jwtService *jwt.JWTService,
) {
	h.registerUserRoutes(userController)
	h.registerTripRoutes(tripController, jwtService)
}

func (h *Handler) registerUserRoutes(userController *controllers.UserController) {
	h.Router.HandleFunc("POST /api/v1/user", userController.HandleCreateUser)
	h.Router.HandleFunc("POST /api/v1/login", userController.HandleLogin)
}

func (h *Handler) registerTripRoutes(tripController *controllers.TripController, jwtSvc *jwt.JWTService) {
	tripHandler := http.HandlerFunc(tripController.HandleTripPreview)
	protectedRoute := AuthMiddleware(jwtSvc)(tripHandler)

	h.Router.Handle("POST /api/v1/trip-preview", protectedRoute)
}

func (h *Handler) GetHandler() http.Handler {
	finalHandler := Chain(
		Logger,
		Recoverer,
		CORS,
	)(h.Router)

	return finalHandler
}
