package http

import (
	"go-ride/services/api-gateway/internal/controllers"
	"go-ride/services/api-gateway/internal/handlers/ws"
	"go-ride/shared/jwt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Handler struct {
	Router     *http.ServeMux
	jwtService *jwt.JWTService
	rdb        *redis.Client
}

func NewHTTPHandler(jwtService *jwt.JWTService, rdb *redis.Client) *Handler {
	router := http.NewServeMux()

	return &Handler{
		Router:     router,
		jwtService: jwtService,
		rdb:        rdb,
	}
}

func (h *Handler) RegisterRoutes(
	userController *controllers.UserController,
	tripController *controllers.TripController,
	driverController *controllers.DriverController,
	driverWSHandler *ws.DriverWSHandler,
) {
	h.registerUserRoutes(userController)
	h.registerTripRoutes(tripController)
	h.registerDriverRoutes(driverController, driverWSHandler)

}

func (h *Handler) registerUserRoutes(userController *controllers.UserController) {
	h.Router.HandleFunc("POST /api/v1/user", userController.HandleCreateUser)
	h.Router.HandleFunc("POST /api/v1/login", userController.HandleLogin)
	h.Router.HandleFunc("POST /api/v1/refresh", userController.HandleRefreshToken)

	h.Router.Handle("POST /api/v1/logout", h.withAuth(userController.HandleLogout))
}

func (h *Handler) registerTripRoutes(tripController *controllers.TripController) {
	h.Router.Handle("POST /api/v1/trip-preview", h.withAuth(tripController.HandleTripPreview))
	h.Router.Handle("POST /api/v1/trip", h.withAuth(tripController.HandleCreateTrip))
}

func (h *Handler) registerDriverRoutes(_ *controllers.DriverController, driverWSHandler *ws.DriverWSHandler) {
	h.Router.Handle("GET /api/v1/driver/stream", h.withAuth(driverWSHandler.HandleConnection))
}

func (h *Handler) withAuth(next http.HandlerFunc) http.Handler {
	return AuthMiddleware(h.jwtService, h.rdb)(next)
}

func (h *Handler) GetHandler() http.Handler {
	finalHandler := Chain(
		Logger,
		Recoverer,
		CORS,
	)(h.Router)

	return finalHandler
}
