package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockRedisService struct {
	mock.Mock
}

// Interface asli hanya punya 2 parameter untuk Set!
func (m *MockRedisService) Set(key string, value interface{}) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockRedisService) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisService) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockRedisService) Exists(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisService) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	args := m.Called(key, value, expiration)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisService) GetTTL(key string) (time.Duration, error) {
	args := m.Called(key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockRedisService) Incr(key string) (int64, error) {
	args := m.Called(key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisService) Decr(key string) (int64, error) {
	args := m.Called(key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisService) HSet(key string, values ...interface{}) error {
	args := m.Called(key, values)
	return args.Error(0)
}

func (m *MockRedisService) HGet(key, field string) (string, error) {
	args := m.Called(key, field)
	return args.String(0), args.Error(1)
}

func (m *MockRedisService) HGetAll(key string) (map[string]string, error) {
	args := m.Called(key)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockRedisService) HDel(key string, fields ...string) error {
	args := m.Called(key, fields)
	return args.Error(0)
}
