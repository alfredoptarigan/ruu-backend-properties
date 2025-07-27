package mocks

import (
	"github.com/stretchr/testify/mock"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/models"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(request dtos.UserRegisterRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockUserService) FindUserByUuid(uuid string) (*models.User, error) {
	args := m.Called(uuid)
	return args.Get(0).(*models.User), args.Error(1)
}
