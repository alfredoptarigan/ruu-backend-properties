package mocks

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"

	"alfredo/ruu-properties/pkg/dtos"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) IsTokenExpired(token string) bool {
	args := m.Called(token)
	return args.Bool(0)
}

func (m *MockJWTService) Revoke(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockJWTService) IsTokenRevoked(token string) bool {
	args := m.Called(token)
	return args.Bool(0)
}

func (m *MockJWTService) GenerateToken(userUuid string, tokens string) (dtos.GenerateTokenResponse, error) {
	args := m.Called(userUuid, tokens)
	if args.Get(0) == nil {
		return dtos.GenerateTokenResponse{}, args.Error(1)
	}
	return args.Get(0).(dtos.GenerateTokenResponse), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (*jwt.Token, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func (m *MockJWTService) GetUserIdFromToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GetRoleFromToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}
