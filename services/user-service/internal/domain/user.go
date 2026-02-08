package domain

import (
	"context"
	"go-ride/shared/types"

	"github.com/google/uuid"
)

type UserModel struct {
	ID             uuid.UUID
	Name           string
	Email          string
	PasswordHashed string
	Type           types.UserType
}

type UserService interface {
	CreateUser(ctx context.Context, user *UserModel) (*UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
	Authenticate(ctx context.Context, email, password string) (*UserModel, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *UserModel) (*UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
}
