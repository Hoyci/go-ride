package http

import (
	"go-ride/services/api-gateway/internal/controllers"
	"go-ride/shared/jwt"
	"net/http"

	"github.com/redis/go-redis/v9"
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
	rdb *redis.Client,
) {
	h.registerUserRoutes(userController, jwtService, rdb)
	h.registerTripRoutes(tripController, jwtService, rdb)
}

func (h *Handler) registerUserRoutes(userController *controllers.UserController, jwtSvc *jwt.JWTService, rdb *redis.Client) {
	h.Router.HandleFunc("POST /api/v1/user", userController.HandleCreateUser)
	h.Router.HandleFunc("POST /api/v1/login", userController.HandleLogin)
	h.Router.HandleFunc("POST /api/v1/refresh", userController.HandleRefreshToken)

	logoutHandler := http.HandlerFunc(userController.HandleLogout)
	protectedRoute := AuthMiddleware(jwtSvc, rdb)(logoutHandler)

	h.Router.Handle("POST /api/v1/logout", protectedRoute)
}

func (h *Handler) registerTripRoutes(tripController *controllers.TripController, jwtSvc *jwt.JWTService, rdb *redis.Client) {
	tripHandler := http.HandlerFunc(tripController.HandleTripPreview)
	protectedRoute := AuthMiddleware(jwtSvc, rdb)(tripHandler)

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
