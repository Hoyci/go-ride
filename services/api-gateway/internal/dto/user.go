package dto

import "go-ride/shared/types"

type CreateUserRequest struct {
	Name     string         `json:"name" validate:"required"`
	Email    string         `json:"email" validate:"required,email"`
	Password string         `json:"password" validate:"required"`
	UserType types.UserType `json:"user_type" validate:"required"`
}
