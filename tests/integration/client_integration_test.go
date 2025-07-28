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

type ClientIntegrationTestSuite struct {
	suite.Suite
	app   *fiber.App
	db    *gorm.DB
	token string
}

func (suite *ClientIntegrationTestSuite) SetupSuite() {
	// Setup test database
	suite.db = config.InitTestDatabase()

	// Setup app
	suite.app = fiber.New()
	router.SetupRoutes(suite.app)

	// Clean database initially
	suite.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	suite.db.Exec("TRUNCATE TABLE clients RESTART IDENTITY CASCADE")
}

func (suite *ClientIntegrationTestSuite) SetupTest() {
	// Clean database before each test
	suite.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	suite.db.Exec("TRUNCATE TABLE clients RESTART IDENTITY CASCADE")

	// Setup auth token after cleaning database
	suite.setupAuthToken()
}

func (suite *ClientIntegrationTestSuite) TearDownSuite() {
	// Clean up
	suite.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	suite.db.Exec("TRUNCATE TABLE clients RESTART IDENTITY CASCADE")

	// Close database connection
	db, _ := suite.db.DB()
	db.Close()
}

// setupAuthToken creates a user and gets authentication token
func (suite *ClientIntegrationTestSuite) setupAuthToken() {
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

func (suite *ClientIntegrationTestSuite) TestCreateClient_Success() {
	clientData := dtos.ClientRequest{
		Name:          "PT. Test Company",
		Email:         "test@company.com",
		PhoneNumber:   "+628123456789",
		Address:       "Jl. Test No. 123, Jakarta",
		ContactPerson: "John Doe",
	}

	clientBody, _ := json.Marshal(clientData)
	req := httptest.NewRequest("POST", "/api/v1/clients", bytes.NewBuffer(clientBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)

	var response dtos.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&response)

	assert.True(suite.T(), response.Success)
	assert.NotNil(suite.T(), response.Data)

	// Verify response data
	if clientResp, ok := response.Data.(map[string]interface{}); ok {
		assert.Equal(suite.T(), clientData.Name, clientResp["name"])
		assert.Equal(suite.T(), clientData.Email, clientResp["email"])
		assert.Equal(suite.T(), clientData.PhoneNumber, clientResp["phone_number"])
		assert.Equal(suite.T(), clientData.Address, clientResp["address"])
		assert.Equal(suite.T(), clientData.ContactPerson, clientResp["contact_person"])
		assert.NotEmpty(suite.T(), clientResp["uuid"])
		assert.NotEmpty(suite.T(), clientResp["created_at"])
		assert.NotEmpty(suite.T(), clientResp["updated_at"])
	}
}

func (suite *ClientIntegrationTestSuite) TestCreateClient_Unauthorized() {
	clientData := dtos.ClientRequest{
		Name:          "PT. Test Company",
		Email:         "test@company.com",
		PhoneNumber:   "+628123456789",
		Address:       "Jl. Test No. 123, Jakarta",
		ContactPerson: "John Doe",
	}

	clientBody, _ := json.Marshal(clientData)
	req := httptest.NewRequest("POST", "/api/v1/clients", bytes.NewBuffer(clientBody))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusUnauthorized, resp.StatusCode)
}

func (suite *ClientIntegrationTestSuite) TestCreateClient_InvalidToken() {
	clientData := dtos.ClientRequest{
		Name:          "PT. Test Company",
		Email:         "test@company.com",
		PhoneNumber:   "+628123456789",
		Address:       "Jl. Test No. 123, Jakarta",
		ContactPerson: "John Doe",
	}

	clientBody, _ := json.Marshal(clientData)
	req := httptest.NewRequest("POST", "/api/v1/clients", bytes.NewBuffer(clientBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid_token")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusUnauthorized, resp.StatusCode)
}

func (suite *ClientIntegrationTestSuite) TestCreateClient_ValidationError() {
	// Test with missing required fields
	clientData := dtos.ClientRequest{
		// Missing required fields
		Email: "invalid-email", // Invalid email format
	}

	clientBody, _ := json.Marshal(clientData)
	req := httptest.NewRequest("POST", "/api/v1/clients", bytes.NewBuffer(clientBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)

	var response dtos.ErrorResponseDTO
	json.NewDecoder(resp.Body).Decode(&response)

	assert.False(suite.T(), response.Success)
	assert.Contains(suite.T(), response.Message, "Invalid request body")
}

func (suite *ClientIntegrationTestSuite) TestGetAllClients_Success() {
	// First create a client
	clientData := dtos.ClientRequest{
		Name:          "PT. Test Company",
		Email:         "test@company.com",
		PhoneNumber:   "+628123456789",
		Address:       "Jl. Test No. 123, Jakarta",
		ContactPerson: "John Doe",
	}

	clientBody, _ := json.Marshal(clientData)
	createReq := httptest.NewRequest("POST", "/api/v1/clients", bytes.NewBuffer(clientBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	createResp, err := suite.app.Test(createReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, createResp.StatusCode)

	// Now get all clients
	getReq := httptest.NewRequest("GET", "/api/v1/clients", nil)
	getReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	getResp, err := suite.app.Test(getReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, getResp.StatusCode)

	var response dtos.PaginatedSuccessResponse
	json.NewDecoder(getResp.Body).Decode(&response)

	assert.True(suite.T(), response.Success)
	assert.NotNil(suite.T(), response.Data)
	assert.NotNil(suite.T(), response.Meta)

	// Check if our created client is in the list
	if clients, ok := response.Data.([]interface{}); ok {
		assert.GreaterOrEqual(suite.T(), len(clients), 1)
	}
}

func (suite *ClientIntegrationTestSuite) TestGetAllClients_WithPagination() {
	// Create multiple clients
	for i := 1; i <= 3; i++ {
		clientData := dtos.ClientRequest{
			Name:          fmt.Sprintf("PT. Test Company %d", i),
			Email:         fmt.Sprintf("test%d@company.com", i),
			PhoneNumber:   fmt.Sprintf("+62812345678%d", i),
			Address:       fmt.Sprintf("Jl. Test No. %d, Jakarta", i),
			ContactPerson: fmt.Sprintf("John Doe %d", i),
		}

		clientBody, _ := json.Marshal(clientData)
		createReq := httptest.NewRequest("POST", "/api/v1/clients", bytes.NewBuffer(clientBody))
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

		createResp, err := suite.app.Test(createReq)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusOK, createResp.StatusCode)
	}

	// Test pagination
	getReq := httptest.NewRequest("GET", "/api/v1/clients?page=1&limit=2", nil)
	getReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	getResp, err := suite.app.Test(getReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, getResp.StatusCode)

	var response dtos.PaginatedSuccessResponse
	json.NewDecoder(getResp.Body).Decode(&response)

	assert.True(suite.T(), response.Success)
	assert.NotNil(suite.T(), response.Data)
	assert.NotNil(suite.T(), response.Meta)

	// Check pagination meta
	meta := response.Meta
	assert.Equal(suite.T(), 1, meta.Page)
	assert.Equal(suite.T(), 2, meta.Limit)
	assert.GreaterOrEqual(suite.T(), meta.Total, 3)
}

func (suite *ClientIntegrationTestSuite) TestGetAllClients_Unauthorized() {
	req := httptest.NewRequest("GET", "/api/v1/clients", nil)
	// No Authorization header

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusUnauthorized, resp.StatusCode)
}

func (suite *ClientIntegrationTestSuite) TestGetClientByID_Success() {
	// First create a client
	clientData := dtos.ClientRequest{
		Name:          "PT. Test Company",
		Email:         "test@company.com",
		PhoneNumber:   "+628123456789",
		Address:       "Jl. Test No. 123, Jakarta",
		ContactPerson: "John Doe",
	}

	clientBody, _ := json.Marshal(clientData)
	createReq := httptest.NewRequest("POST", "/api/v1/clients", bytes.NewBuffer(clientBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	createResp, err := suite.app.Test(createReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, createResp.StatusCode)

	var createResponse dtos.SuccessResponse
	json.NewDecoder(createResp.Body).Decode(&createResponse)

	// Extract client UUID
	var clientUUID string
	if clientResp, ok := createResponse.Data.(map[string]interface{}); ok {
		clientUUID = clientResp["uuid"].(string)
	}

	assert.NotEmpty(suite.T(), clientUUID)

	// Now get client by ID
	getReq := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/clients/%s", clientUUID), nil)
	getReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	getResp, err := suite.app.Test(getReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, getResp.StatusCode)

	var response dtos.SuccessResponse
	json.NewDecoder(getResp.Body).Decode(&response)

	assert.True(suite.T(), response.Success)
	assert.NotNil(suite.T(), response.Data)

	// Verify response data
	if clientResp, ok := response.Data.(map[string]interface{}); ok {
		assert.Equal(suite.T(), clientUUID, clientResp["uuid"])
		assert.Equal(suite.T(), clientData.Name, clientResp["name"])
		assert.Equal(suite.T(), clientData.Email, clientResp["email"])
	}
}

func (suite *ClientIntegrationTestSuite) TestGetClientByID_NotFound() {
	// Use a non-existent UUID
	fakeUUID := "550e8400-e29b-41d4-a716-446655440000"

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/clients/%s", fakeUUID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (suite *ClientIntegrationTestSuite) TestGetClientByID_InvalidUUID() {
	// Use an invalid UUID format
	invalidUUID := "invalid-uuid"

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/clients/%s", invalidUUID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)

	var response dtos.ErrorResponseDTO
	json.NewDecoder(resp.Body).Decode(&response)

	assert.False(suite.T(), response.Success)
	assert.Contains(suite.T(), response.Message, "Invalid UUID format")
}

func (suite *ClientIntegrationTestSuite) TestGetClientByID_Unauthorized() {
	fakeUUID := "550e8400-e29b-41d4-a716-446655440000"

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/clients/%s", fakeUUID), nil)
	// No Authorization header

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusUnauthorized, resp.StatusCode)
}

func TestClientIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ClientIntegrationTestSuite))
}

func (suite *ClientIntegrationTestSuite) TestDeleteClient_Success() {
	// Create a client first
	clientData := dtos.ClientRequest{
		Name:          "PT. Test Company",
		Email:         "test@company.com",
		PhoneNumber:   "+628123456789",
		Address:       "Jl. Test No. 123, Jakarta",
		ContactPerson: "John Doe",
	}

	clientBody, _ := json.Marshal(clientData)
	createReq := httptest.NewRequest("POST", "/api/v1/clients", bytes.NewBuffer(clientBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	createResp, err := suite.app.Test(createReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, createResp.StatusCode)

	var createResponse dtos.SuccessResponse
	json.NewDecoder(createResp.Body).Decode(&createResponse)

	// Extract client UUID
	var clientUUID string
	if clientResp, ok := createResponse.Data.(map[string]interface{}); ok {
		clientUUID = clientResp["uuid"].(string)
	}

	assert.NotEmpty(suite.T(), clientUUID)

	// Now delete the client
	deleteReq := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/clients/%s/delete", clientUUID), nil)
	deleteReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	deleteResp, err := suite.app.Test(deleteReq)
	assert.NoError(suite.T(), err)
	// Print response body for debugging
	respBody, _ := io.ReadAll(deleteResp.Body)
	fmt.Printf("Delete Response Body: %s\n", string(respBody))
	// Reset body for further reading
	deleteResp.Body = io.NopCloser(bytes.NewBuffer(respBody))

	assert.Equal(suite.T(), fiber.StatusOK, deleteResp.StatusCode)

	var deleteResponse dtos.SuccessResponse
	json.NewDecoder(deleteResp.Body).Decode(&deleteResponse)

	assert.True(suite.T(), deleteResponse.Success)

}

func (suite *ClientIntegrationTestSuite) TestDeleteClient_NotFound() {
	// Use a non-existent UUID
	fakeUUID := "550e8400-e29b-41d4-a716-446655440000"

	deleteReq := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/clients/%s/delete", fakeUUID), nil)
	deleteReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	deleteResp, err := suite.app.Test(deleteReq)
	assert.NoError(suite.T(), err)
	// Print response body for debugging
	respBody, _ := io.ReadAll(deleteResp.Body)
	fmt.Printf("Delete Response Body: %s\n", string(respBody))
	// Reset body for further reading
	deleteResp.Body = io.NopCloser(bytes.NewBuffer(respBody))

	assert.Equal(suite.T(), fiber.StatusNotFound, deleteResp.StatusCode)
}

func (suite *ClientIntegrationTestSuite) TestUpdateClient_Success() {
	// Create a client first
	clientData := dtos.ClientRequest{
		Name:          "PT. Test Company",
		Email:         "test@company.com",
		PhoneNumber:   "+628123456789",
		Address:       "Jl. Test No. 123, Jakarta",
		ContactPerson: "John Doe",
	}

	clientBody, _ := json.Marshal(clientData)

	createReq := httptest.NewRequest("POST", "/api/v1/clients", bytes.NewBuffer(clientBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	createResp, err := suite.app.Test(createReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, createResp.StatusCode)

	var createResponse dtos.SuccessResponse
	json.NewDecoder(createResp.Body).Decode(&createResponse)

	// Extract client UUID
	var clientUUID string
	if clientResp, ok := createResponse.Data.(map[string]interface{}); ok {
		clientUUID = clientResp["uuid"].(string)
	}

	assert.NotEmpty(suite.T(), clientUUID)

	// Update the client
	updateData := dtos.ClientUpdateRequest{
		Name:          "PT. Updated Company",
		Email:         "updated@company.com",
		PhoneNumber:   "+628123456789",
		Address:       "Jl. Updated No. 456, Jakarta",
		ContactPerson: "Jane Doe",
	}

	updateBody, _ := json.Marshal(updateData)

	updateReq := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/clients/%s/update", clientUUID), bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.token))

	updateResp, err := suite.app.Test(updateReq)

	// Print response body for debugging
	respBody, _ := io.ReadAll(updateResp.Body)
	fmt.Printf("Update Response Body: %s\n", string(respBody))
	// Reset body for further reading
	updateResp.Body = io.NopCloser(bytes.NewBuffer(respBody))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, updateResp.StatusCode)

	var updateResponse dtos.SuccessResponse
	json.NewDecoder(updateResp.Body).Decode(&updateResponse)
	assert.True(suite.T(), updateResponse.Success)
	assert.Equal(suite.T(), "Client updated successfully", updateResponse.Message)

}
