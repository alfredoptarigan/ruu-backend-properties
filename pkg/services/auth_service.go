package services

import (
	"fmt"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/helpers"
	"alfredo/ruu-properties/pkg/repositories"
)

type AuthService interface {
	Login(request dtos.LoginRequest) (response dtos.LoginResponse, err error)
}

type authServiceImpl struct {
	jwtService     JwtService
	userService    UserService
	userRepository repositories.UserRepository
	redisService   RedisService
}

func (a *authServiceImpl) Login(request dtos.LoginRequest) (response dtos.LoginResponse, err error) {
	// Find user by email
	user, err := a.userRepository.FindUserByEmail(request.Email)
	if err != nil {
		return response, fmt.Errorf("user not found")
	}

	// Check if the password same with the hash password
	if ok, err := helpers.CheckPasswordHashWithArgon2(request.Password, user.Password); !ok || err != nil {
		return response, fmt.Errorf("invalid email or password")
	}

	generateToken := helpers.GenerateToken(32)
	token, err := a.jwtService.GenerateToken(user.UUID, generateToken)
	if err != nil {
		return response, fmt.Errorf("failed to generate token")
	}

	return dtos.LoginResponse{
		TokenType:    "Bearer",
		ExpiresIn:    int64(token.ExpiresIn),
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Email:        user.Email,
		UserUuid:     user.UUID,
		Name:         user.Name,
	}, nil
}

func NewAuthService(
	jwtService JwtService,
	userService UserService,
	userRepository repositories.UserRepository,
	redisService RedisService,
) AuthService {
	return &authServiceImpl{
		jwtService:     jwtService,
		userService:    userService,
		userRepository: userRepository,
		redisService:   redisService,
	}
}
