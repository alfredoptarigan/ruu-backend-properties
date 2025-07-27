package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/testutils"
	"alfredo/ruu-properties/pkg/testutils/mocks"
)

type UserControllerTestSuite struct {
	suite.Suite
	app              *fiber.App
	userController   UserController
	mockUserService  *mocks.MockUserService
	mockRedisService *mocks.MockRedisService
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.app = testutils.CreateTestApp()
	suite.mockUserService = new(mocks.MockUserService)
	suite.mockRedisService = new(mocks.MockRedisService)

	suite.userController = NewUserController(
		suite.mockRedisService,
		suite.mockUserService,
	)

	// Setup routes
	user := suite.app.Group("/user")
	suite.userController.Router(user)
}

func (suite *UserControllerTestSuite) TestRegister_Success() {
	// Arrange
	suite.mockUserService.On("Register", mock.AnythingOfType("dtos.UserRegisterRequest")).Return(nil)

	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fields := map[string]string{
		"name":                  "John Doe",
		"email":                 "john@example.com",
		"password":              "password123",
		"confirmation_password": "password123",
		"phone_number":          "+1234567890",
		"role":                  "user",
	}

	for key, value := range fields {
		writer.WriteField(key, value)
	}
	writer.Close()

	// Act
	req := httptest.NewRequest("POST", "/user/register", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := suite.app.Test(req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)

	var response dtos.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&response)

	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "User registered successfully", response.Message)

	suite.mockUserService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestRegister_ServiceError() {
	// Arrange
	suite.mockUserService.On("Register", mock.AnythingOfType("dtos.UserRegisterRequest")).Return(errors.New("email already exists"))

	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fields := map[string]string{
		"name":                  "John Doe",
		"email":                 "existing@example.com",
		"password":              "password123",
		"confirmation_password": "password123",
		"phone_number":          "+1234567890",
		"role":                  "user",
	}

	for key, value := range fields {
		writer.WriteField(key, value)
	}
	writer.Close()

	// Act
	req := httptest.NewRequest("POST", "/user/register", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := suite.app.Test(req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusInternalServerError, resp.StatusCode)

	var response dtos.ErrorResponseDTO
	json.NewDecoder(resp.Body).Decode(&response)

	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "email already exists", response.Message)

	suite.mockUserService.AssertExpectations(suite.T())
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
