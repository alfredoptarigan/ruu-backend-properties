package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// CreateTestApp creates a new Fiber app for testing
func CreateTestApp() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
}

// MakeJSONRequest creates a JSON request for testing
func MakeJSONRequest(app *fiber.App, method, url string, body interface{}) (*httptest.ResponseRecorder, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		return nil, err
	}

	recorder := httptest.NewRecorder()
	resp.Body.Close()

	// Copy response headers
	for k, v := range resp.Header {
		recorder.Header()[k] = v
	}

	// Copy response status code
	recorder.WriteHeader(resp.StatusCode)

	// Copy response body
	if _, err := io.Copy(recorder, resp.Body); err != nil {
		return nil, err
	}

	return recorder, nil
}

// MakeMultipartRequest creates a multipart form request for testing
func MakeMultipartRequest(app *fiber.App, method, url string, fields map[string]string, files map[string][]byte) (*httptest.ResponseRecorder, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	// Add files
	for fieldName, fileContent := range files {
		part, err := writer.CreateFormFile(fieldName, "test.jpg")
		if err != nil {
			return nil, err
		}
		if _, err := part.Write(fileContent); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	if err != nil {
		return nil, err
	}

	recorder := httptest.NewRecorder()
	resp.Body.Close()

	// Copy response headers
	for k, v := range resp.Header {
		recorder.Header()[k] = v
	}

	// Copy response status code
	recorder.WriteHeader(resp.StatusCode)

	// Copy response body
	if _, err := io.Copy(recorder, resp.Body); err != nil {
		return nil, err
	}

	return recorder, nil
}

// AssertJSONResponse asserts JSON response structure
func AssertJSONResponse(t *testing.T, resp *httptest.ResponseRecorder, expectedStatus int, expectedSuccess bool) map[string]interface{} {
	assert.Equal(t, expectedStatus, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedSuccess, response["success"])

	return response
}
