package service

import (
	"context"
	"errors"
	"fmt"
	"go-ride/services/user-service/internal/domain"
	"go-ride/services/user-service/pkg"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type userService struct {
	repo domain.UserRepository
}

func NewUserSerivce(repo domain.UserRepository) *userService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *domain.UserModel) (*domain.UserModel, error) {
	existingUser, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := pkg.HashPassword(user.PasswordHashed)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}
	user.PasswordHashed = hashedPassword

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.UserModel, error) {
	return s.repo.GetUserByEmail(ctx, email)
}
