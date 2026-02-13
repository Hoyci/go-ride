package grpc

import (
	"context"
	"errors"
	"go-ride/services/user-service/internal/domain"
	"go-ride/services/user-service/internal/service"
	"go-ride/shared/jwt"
	pu "go-ride/shared/proto/user"
	"go-ride/shared/types"
	"log"

	jwtLib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pu.UnimplementedUserServiceServer
	userService domain.UserService
	jwtService  *jwt.JWTService
	rdb         *redis.Client
}

func NewGRPCHandler(server *grpc.Server, userService domain.UserService, jwtService *jwt.JWTService, rdb *redis.Client) *gRPCHandler {
	handler := &gRPCHandler{
		userService: userService,
		jwtService:  jwtService,
		rdb:         rdb,
	}

	pu.RegisterUserServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) CreateUser(ctx context.Context, req *pu.CreateUserRequest) (*pu.CreateUserResponse, error) {
	userModel := &domain.UserModel{
		Name:           req.Name,
		Email:          req.Email,
		PasswordHashed: req.Password,
		Type:           types.UserType(req.Type.String()),
	}

	user, err := h.userService.CreateUser(ctx, userModel)

	if err != nil {
		// Se o erro for que o usuário já existe, aplicamos a máscara de segurança
		// Retornamos sucesso falso para evitar User Enumeration
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return &pu.CreateUserResponse{Id: uuid.New().String()}, nil
		}

		log.Printf("create user error: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pu.CreateUserResponse{Id: user.ID.String()}, nil
}

func (h *gRPCHandler) Login(ctx context.Context, req *pu.LoginRequest) (*pu.LoginResponse, error) {
	user, err := h.userService.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("login error: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	accessToken, accessJTI, _ := h.jwtService.GenerateToken(user.ID.String(), jwt.ACCESS)
	refreshToken, refreshJTI, _ := h.jwtService.GenerateToken(user.ID.String(), jwt.REFRESH)

	h.rdb.Set(ctx, "session:"+user.ID.String(), accessJTI, jwt.ACCESS_EXPIRATION)
	h.rdb.Set(ctx, "refresh_session:"+user.ID.String(), refreshJTI, jwt.REFRESH_EXPIRATION)

	return &pu.LoginResponse{
		Id:           user.ID.String(),
		Name:         user.Name,
		Email:        user.Email,
		Type:         types.MapUserTypeDomainToProto(user.Type).String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *gRPCHandler) RefreshToken(ctx context.Context, req *pu.RefreshTokenRequest) (*pu.RefreshTokenResponse, error) {
	token, err := h.jwtService.Validate(req.Token)
	if err != nil || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "refresh token is invalid or expired")
	}

	claims, ok := token.Claims.(jwtLib.MapClaims)
	if !ok || claims["type"] != "REFRESH" {
		return nil, status.Error(codes.Unauthenticated, "invalid token type")
	}

	userID := claims["sub"].(string)
	incomingJTI := claims["jti"].(string)

	currentRefreshJTI, err := h.rdb.Get(ctx, "refresh_session:"+userID).Result()
	if err != nil || currentRefreshJTI != incomingJTI {
		h.rdb.Del(ctx, "session:"+userID, "refresh_session:"+userID)
		return nil, status.Error(codes.Unauthenticated, "refresh token invalidated")
	}

	newAccess, newAccessJTI, _ := h.jwtService.GenerateToken(userID, jwt.ACCESS)
	newRefresh, newRefreshJTI, _ := h.jwtService.GenerateToken(userID, jwt.REFRESH)

	h.rdb.Set(ctx, "session:"+userID, newAccessJTI, jwt.ACCESS_EXPIRATION)
	h.rdb.Set(ctx, "refresh_session:"+userID, newRefreshJTI, jwt.REFRESH_EXPIRATION)

	return &pu.RefreshTokenResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

func (h *gRPCHandler) Logout(ctx context.Context, req *pu.LogoutRequest) (*pu.LogoutResponse, error) {
	err := h.rdb.Del(ctx, "session:"+req.UserID).Err()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete session: %v", err)
	}

	return &pu.LogoutResponse{Success: true}, nil
}
