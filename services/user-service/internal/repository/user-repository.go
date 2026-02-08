package repository

import (
	"context"
	"go-ride/services/user-service/internal/domain"
)

type inmemRepository struct {
	users map[string]*domain.UserModel
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		users: make(map[string]*domain.UserModel),
	}
}

func (r *inmemRepository) CreateUser(ctx context.Context, user *domain.UserModel) (*domain.UserModel, error) {
	r.users[user.ID.String()] = user
	return user, nil
}

func (r *inmemRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserModel, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, nil
}
