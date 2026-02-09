package mocks

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jarcoal/httpmock"
)

// AnalyticsMock handles mock HTTP responses for analytics
type AnalyticsMock struct{}

// loadMockResponse loads a mock response file from the mocks directory
func loadMockResponse(filename string) ([]byte, error) {
	// Get the path relative to the current working directory during test execution
	// Tests run from the service package directory, so mocks/ is a subdirectory
	mockPath := filepath.Join("mocks", filename)
	data, err := os.ReadFile(mockPath)
	if err != nil {
		// Try absolute path resolution as fallback
		absPath, _ := filepath.Abs(mockPath)
		return nil, fmt.Errorf("failed to load mock file %s (tried: %s, %s): %w", filename, mockPath, absPath, err)
	}
	return data, nil
}

// RegisterMocks registers all success mock responses
func (m *AnalyticsMock) RegisterMocks(baseURL string) {
	// Mock GET /analytics.json
	httpmock.RegisterResponder("GET", baseURL+"/analytics.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("analytics_list.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			return httpmock.NewBytesResponse(200, mockData), nil
		},
	)

	// Mock GET /analytics.csv
	httpmock.RegisterResponder("GET", baseURL+"/analytics.csv",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("analytics_list.csv")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "text/csv")
			return resp, nil
		},
	)
}

// RegisterErrorMocks registers error response mocks
func (m *AnalyticsMock) RegisterErrorMocks(baseURL string) {
	// Mock unauthorized error for JSON
	httpmock.RegisterResponder("GET", baseURL+"/analytics.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_unauthorized.json")
			if err != nil {
				return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid API key"]}`), nil
			}
			return httpmock.NewBytesResponse(401, mockData), nil
		},
	)

	// Mock unauthorized error for CSV
	httpmock.RegisterResponder("GET", baseURL+"/analytics.csv",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_unauthorized.json")
			if err != nil {
				return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid API key"]}`), nil
			}
			return httpmock.NewBytesResponse(401, mockData), nil
		},
	)
}

// CleanupMockState cleans up any state from the mock
func (m *AnalyticsMock) CleanupMockState() {
	// No stateful data to cleanup for analytics
}
