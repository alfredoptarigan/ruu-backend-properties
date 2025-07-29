package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"alfredo/ruu-properties/config"
	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/router"
)

type FeatureIntegrationTestSuite struct {
	suite.Suite
	app   *fiber.App
	db    *gorm.DB
	token string
}

func (suite *FeatureIntegrationTestSuite) SetupSuite() {
	suite.db = config.InitTestDatabase()

	suite.app = fiber.New()
	router.SetupRoutes(suite.app)
}

func (suite *FeatureIntegrationTestSuite) SetupTest() {
	// Clean database before each test
	suite.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	suite.db.Exec("TRUNCATE TABLE features RESTART IDENTITY CASCADE")

	// Setup auth token after cleaning database
	suite.setupAuthToken()
}

func (suite *FeatureIntegrationTestSuite) TearDownSuite() {
	// Clean up
	suite.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	suite.db.Exec("TRUNCATE TABLE features RESTART IDENTITY CASCADE")

	// Close database connection
	db, _ := suite.db.DB()
	db.Close()
}

// setupAuthToken creates a user and gets authentication token
func (suite *FeatureIntegrationTestSuite) setupAuthToken() {
	// Generate unique email for each test run
	timestamp := time.Now().UnixNano()
	email := fmt.Sprintf("integration-%d@test.com", timestamp)

	// First, register a user
	registerData := map[string]string{
		"name":                  "Integration Test User",
		"email":                 email, // Use unique email
		"password":              "password123",
		"confirmation_password": "password123",
		"phone_number":          fmt.Sprintf("+123456789%d", timestamp%1000), // Unique phone too
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
	fmt.Println("Register Response Status:", registerResp.StatusCode)

	// Read response body for debugging
	respBody, _ := io.ReadAll(registerResp.Body)
	fmt.Println("Register Response Body:", string(respBody))

	// Reset body for further reading
	registerResp.Body = io.NopCloser(bytes.NewBuffer(respBody))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, registerResp.StatusCode)

	// Login to get token
	loginData := dtos.LoginRequest{
		Email:    email,
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

	// Extract token from response
	if data, ok := loginResponse.Data.(map[string]interface{}); ok {
		if token, ok := data["access_token"].(string); ok {
			suite.token = token
		}
	}

	assert.NotEmpty(suite.T(), suite.token, "Token should not be empty")
}

func (suite *FeatureIntegrationTestSuite) TestCreateFeature_Success() {
	featureData := dtos.FeatureRequest{
		Name:        "Kolam Renang",
		Description: "Kolam Renang dibuka selama 24 jam",
	}

	clientBody, _ := json.Marshal(featureData)
	req := httptest.NewRequest("POST", "/api/v1/features", bytes.NewBuffer(clientBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusCreated, resp.StatusCode)

	var response dtos.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&response)

	// Println response from backend api
	fmt.Println("Create Feature Response:", response)

	assert.True(suite.T(), response.Success)
	assert.NotNil(suite.T(), response.Data)

	// Verify response data
	// Convert response data to map for easier assertion
	responseData := response.Data.(map[string]interface{})

	// Assert feature data matches what was sent
	assert.Equal(suite.T(), featureData.Name, responseData["name"])
	assert.Equal(suite.T(), featureData.Description, responseData["description"])
	assert.NotEmpty(suite.T(), responseData["uuid"])
}

func TestFeatureIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FeatureIntegrationTestSuite))
}
