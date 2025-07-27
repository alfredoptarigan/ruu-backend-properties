package services

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/models"
	"alfredo/ruu-properties/pkg/repositories"
)

type UserService interface {
	Register(request dtos.UserRegisterRequest) error
	FindUserByUuid(uuid string) (*models.User, error)
}

type userServiceImpl struct {
	userRepository repositories.UserRepository
}

// FindUserByUuid implements UserService.
func (u *userServiceImpl) FindUserByUuid(uuid string) (*models.User, error) {
	user, err := u.userRepository.FindUserByUuid(uuid)
	if err != nil {
		return nil, fmt.Errorf("%s", "user not found.")
	}
	return &user, nil
}

func (u userServiceImpl) Register(request dtos.UserRegisterRequest) error {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return err
	}

	if err := u.userRepository.Register(request); err != nil {
		return err
	}

	return nil
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userServiceImpl{userRepository: userRepository}
}
