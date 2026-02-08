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

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pu.UnimplementedUserServiceServer
	userService domain.UserService
	jwtService  *jwt.JWTService
}

func NewGRPCHandler(server *grpc.Server, userService domain.UserService, jwtService *jwt.JWTService) *gRPCHandler {
	handler := &gRPCHandler{
		userService: userService,
		jwtService:  jwtService,
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

	accessToken, _ := h.jwtService.GenerateToken(user.ID.String(), jwt.ACCESS)
	refreshToken, _ := h.jwtService.GenerateToken(user.ID.String(), jwt.REFRESH)

	return &pu.LoginResponse{
		Id:           user.ID.String(),
		Name:         user.Name,
		Email:        user.Email,
		Type:         types.MapUserTypeDomainToProto(user.Type).String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
