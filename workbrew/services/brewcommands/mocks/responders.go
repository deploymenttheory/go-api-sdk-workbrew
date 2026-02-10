package mocks

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jarcoal/httpmock"
)

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"message":"Not Found","errors":["Resource not found"]}`))
}

// loadMockResponse loads a mock response file
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

// BrewCommandsMock handles mock responses for brew commands service
type BrewCommandsMock struct{}

// RegisterMocks registers all HTTP mock responders for brew commands service
func (m *BrewCommandsMock) RegisterMocks(baseURL string) {
	// GET /brew_commands.json
	httpmock.RegisterResponder("GET", baseURL+"/brew_commands.json", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_get_brew_commands.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
		}

		resp := httpmock.NewBytesResponse(200, mockData)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})

	// POST /brew_commands.json
	httpmock.RegisterResponder("POST", baseURL+"/brew_commands.json", func(req *http.Request) (*http.Response, error) {
		// Read request body
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return httpmock.NewStringResponse(400, `{"message":"Bad Request","errors":["Invalid request body"]}`), nil
		}

		// Validate request
		var requestBody map[string]any
		if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"message":"Bad Request","errors":["Invalid JSON"]}`), nil
		}

		// Check for required arguments field
		if _, ok := requestBody["arguments"]; !ok {
			return httpmock.NewStringResponse(422, `{"message":"Validation Error","errors":["Arguments field is required"]}`), nil
		}

		// Load success response
		mockData, err := loadMockResponse("validate_create_brew_command.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
		}

		resp := httpmock.NewBytesResponse(201, mockData)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})
}

// RegisterErrorMocks registers mock responders that return error responses
func (m *BrewCommandsMock) RegisterErrorMocks(baseURL string) {
	// GET /brew_commands.json - Return error
	httpmock.RegisterResponder("GET", baseURL+"/brew_commands.json", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid or missing API key"]}`), nil
	})

	// POST /brew_commands.json - Return forbidden (Free tier)
	httpmock.RegisterResponder("POST", baseURL+"/brew_commands.json", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(403, `{"message":"An error occurred when trying to create Brew Command","errors":["Please upgrade your plan to get access to Brew Commands."]}`), nil
	})
}

// CleanupMockState clears all mock state
func (m *BrewCommandsMock) CleanupMockState() {
	// No state to clean for brew commands service
}
