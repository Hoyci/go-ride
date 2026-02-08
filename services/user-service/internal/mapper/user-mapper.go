// mapper/user_mapper.go
package mapper

import (
	"fmt"

	"go-ride/services/user-service/internal/domain"
	"go-ride/services/user-service/internal/dto"
	pu "go-ride/shared/proto/user"
)

func DomainToDTO(u *domain.UserModel) (*dto.UserDTO, error) {
	if u == nil {
		return nil, fmt.Errorf("user is nil")
	}

	return &dto.UserDTO{
		ID:             u.ID.String(),
		Name:           u.Name,
		Email:          u.Email,
		PasswordHashed: u.PasswordHashed,
		Type:           string(u.Type),
	}, nil
}

func DTOToProto(d *dto.UserDTO) (*pu.User, error) {
	switch d.Type {
	case "DRIVER":
		return &pu.User{
			Id:             d.ID,
			Name:           d.Name,
			Email:          d.Email,
			PasswordHashed: d.PasswordHashed,
			Type:           pu.UserType_DRIVER,
		}, nil
	case "PASSENGER":
		return &pu.User{
			Id:             d.ID,
			Name:           d.Name,
			Email:          d.Email,
			PasswordHashed: d.PasswordHashed,
			Type:           pu.UserType_PASSENGER,
		}, nil
	default:
		return nil, fmt.Errorf("invalid user type: %s", d.Type)
	}
}
