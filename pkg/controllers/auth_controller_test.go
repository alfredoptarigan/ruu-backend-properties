package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/testutils"
	"alfredo/ruu-properties/pkg/testutils/mocks"
)

type AuthControllerTestSuite struct {
	suite.Suite
	app              *fiber.App
	authController   AuthController
	mockAuthService  *mocks.MockAuthService
	mockRedisService *mocks.MockRedisService
	mockJWTService   *mocks.MockJWTService
}

func (suite *AuthControllerTestSuite) SetupTest() {
	suite.app = testutils.CreateTestApp()
	suite.mockAuthService = new(mocks.MockAuthService)
	suite.mockRedisService = new(mocks.MockRedisService)
	suite.mockJWTService = new(mocks.MockJWTService)

	suite.authController = NewAuthController(
		suite.mockAuthService,
		suite.mockRedisService,
		suite.mockJWTService,
	)

	// Setup routes - this was missing!
	auth := suite.app.Group("/auth")
	suite.authController.Router(auth)
}

func (suite *AuthControllerTestSuite) TestLogin_Success() {
	// Arrange
	loginRequest := dtos.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedResponse := dtos.LoginResponse{
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}

	suite.mockAuthService.On("Login", loginRequest).Return(expectedResponse, nil)

	// Act
	reqBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)

	var response dtos.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&response)

	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Login successful", response.Message)
	assert.NotNil(suite.T(), response.Data)

	suite.mockAuthService.AssertExpectations(suite.T())
}

func (suite *AuthControllerTestSuite) TestLogin_InvalidRequestBody() {
	// Act
	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)

	var response dtos.ErrorResponseDTO
	json.NewDecoder(resp.Body).Decode(&response)

	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "Invalid request body", response.Message)
}

func (suite *AuthControllerTestSuite) TestLogin_ServiceError() {
	// Arrange
	loginRequest := dtos.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	suite.mockAuthService.On("Login", loginRequest).Return(dtos.LoginResponse{}, errors.New("invalid credentials"))

	// Act
	reqBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusInternalServerError, resp.StatusCode)

	var response dtos.ErrorResponseDTO
	json.NewDecoder(resp.Body).Decode(&response)

	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "invalid credentials", response.Message)

	suite.mockAuthService.AssertExpectations(suite.T())
}

func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}
