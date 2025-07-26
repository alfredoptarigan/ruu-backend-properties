package services

import (
	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/repositories"
	"github.com/go-playground/validator/v10"
)

type UserService interface {
	Register(request dtos.UserRegisterRequest) error
}

type userServiceImpl struct {
	userRepository repositories.UserRepository
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
