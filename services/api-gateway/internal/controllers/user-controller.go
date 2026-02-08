package controllers

import (
	"context"
	"encoding/json"
	"go-ride/services/api-gateway/internal/dto"
	"go-ride/shared/contracts"
	pu "go-ride/shared/proto/user"
	"go-ride/shared/responses"
	"go-ride/shared/types"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type UserController struct {
	validator   *validator.Validate
	userService pu.UserServiceClient
}

func NewUserController(v *validator.Validate, us pu.UserServiceClient) *UserController {
	return &UserController{
		validator:   v,
		userService: us,
	}
}

func (s *UserController) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.WriteJSON(w, http.StatusBadRequest, contracts.APIResponse{
			Error: &contracts.APIError{
				Code:    http.StatusBadRequest,
				Message: "invalid JSON payload",
			},
		})
		return
	}

	if err := s.validator.Struct(req); err != nil {
		responses.WriteJSON(w, http.StatusUnprocessableEntity, contracts.APIResponse{
			Error: &contracts.APIError{
				Code:    http.StatusUnprocessableEntity,
				Message: "validation failed",
				Details: responses.ParseValidationErrors(err),
			},
		})
		return
	}

	grpcRes, err := s.userService.CreateUser(
		ctx,
		&pu.CreateUserRequest{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
			Type:     types.MapUserTypeToProto(req.UserType),
		},
	)

	if err != nil {
		log.Printf("failed to call create user: %v", err)
		responses.WriteJSON(w, http.StatusInternalServerError, contracts.APIResponse{
			Error: &contracts.APIError{Message: "failed to contact user service"},
		})
		return
	}

	responses.WriteJSON(w, http.StatusOK, contracts.APIResponse{
		Data: grpcRes,
	})
}
