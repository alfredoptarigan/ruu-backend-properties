package integration

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"alfredo/ruu-properties/config"
	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/router"
)

type AuthIntegrationTestSuite struct {
	suite.Suite
	app *fiber.App
	db  *gorm.DB
}

func (suite *AuthIntegrationTestSuite) SetupSuite() {
	// Setup test database
	suite.db = config.InitTestDatabase()

	// Setup app
	suite.app = fiber.New()
	router.SetupRoutes(suite.app)
}

func (suite *AuthIntegrationTestSuite) SetupTest() {
	// Clean database before each test
	suite.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
}

func (suite *AuthIntegrationTestSuite) TearDownSuite() {
	// Close database connection
	db, _ := suite.db.DB()
	db.Close()
}

func (suite *AuthIntegrationTestSuite) TestRegisterAndLogin_Success() {
	// First, register a user
	registerData := map[string]string{
		"name":                  "Integration Test User",
		"email":                 "integration2@test.com",
		"password":              "password123",
		"confirmation_password": "password123",
		"phone_number":          "+1234567890",
		"role":                  "user",
	}

	// Create multipart form for registration
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, value := range registerData {
		writer.WriteField(key, value)
	}
	writer.Close()

	// Register user
	registerReq := httptest.NewRequest("POST", "/api/v1/user/register", body)
	registerReq.Header.Set("Content-Type", writer.FormDataContentType())

	registerResp, err := suite.app.Test(registerReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, registerResp.StatusCode)

	// Now login with the registered user
	loginData := dtos.LoginRequest{
		Email:    "integration@test.com",
		Password: "password123",
	}

	loginBody, _ := json.Marshal(loginData)
	loginReq := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")

	loginResp, err := suite.app.Test(loginReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, loginResp.StatusCode)

	var loginResponse dtos.SuccessResponse
	json.NewDecoder(loginResp.Body).Decode(&loginResponse)

	assert.True(suite.T(), loginResponse.Success)
	assert.Equal(suite.T(), "Login successful", loginResponse.Message)
	assert.NotNil(suite.T(), loginResponse.Data)
}

func TestAuthIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}
