package mocks

import (
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

// DevicesMock handles mock responses for devices service
type DevicesMock struct{}

// RegisterMocks registers all HTTP mock responders for devices service
func (m *DevicesMock) RegisterMocks(baseURL string) {
	// GET /devices.json
	httpmock.RegisterResponder("GET", baseURL+"/devices.json", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_get_devices.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
		}

		resp := httpmock.NewBytesResponse(200, mockData)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})

	// GET /devices.csv
	httpmock.RegisterResponder("GET", baseURL+"/devices.csv", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_get_devices_csv.txt")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
		}

		resp := httpmock.NewBytesResponse(200, mockData)
		resp.Header.Set("Content-Type", "text/csv")
		return resp, nil
	})
}

// RegisterErrorMocks registers mock responders that return error responses
func (m *DevicesMock) RegisterErrorMocks(baseURL string) {
	// GET /devices.json - Return unauthorized error
	httpmock.RegisterResponder("GET", baseURL+"/devices.json", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid or missing API key"]}`), nil
	})

	// GET /devices.csv - Return forbidden error
	httpmock.RegisterResponder("GET", baseURL+"/devices.csv", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(403, `{"message":"Forbidden","errors":["Please upgrade your plan to get access to this feature."]}`), nil
	})
}

// CleanupMockState clears all mock state
func (m *DevicesMock) CleanupMockState() {
	// No state to clean for devices service
}
