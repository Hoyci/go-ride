package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserType string

const (
	DRIVER    UserType = "DRIVER"
	PASSENGER UserType = "PASSENGER"
)

type UserModel struct {
	ID             uuid.UUID
	Name           string
	Email          string
	PasswordHashed string
	Type           UserType
}

type UserService interface {
	CreateUser(ctx context.Context, user *UserModel) (*UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *UserModel) (*UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
}
