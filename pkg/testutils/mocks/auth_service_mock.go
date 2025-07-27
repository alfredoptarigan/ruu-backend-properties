package mocks

import (
	"github.com/stretchr/testify/mock"

	"alfredo/ruu-properties/pkg/dtos"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(request dtos.LoginRequest) (dtos.LoginResponse, error) {
	args := m.Called(request)
	return args.Get(0).(dtos.LoginResponse), args.Error(1)
}
