package grpc

import (
	"context"
	"errors"
	"go-ride/services/user-service/internal/domain"
	"go-ride/services/user-service/internal/service"
	pu "go-ride/shared/proto/user"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pu.UnimplementedUserServiceServer
	userService domain.UserService
}

func NewGRPCHandler(server *grpc.Server, userService domain.UserService) *gRPCHandler {
	handler := &gRPCHandler{
		userService: userService,
	}

	pu.RegisterUserServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) CreateUser(ctx context.Context, req *pu.CreateUserRequest) (*pu.CreateUserResponse, error) {
	userModel := &domain.UserModel{
		Name:           req.Name,
		Email:          req.Email,
		PasswordHashed: req.Password,
		Type:           domain.UserType(req.Type.String()),
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
