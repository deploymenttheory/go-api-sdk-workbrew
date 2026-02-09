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

// AnalyticsMock handles mock HTTP responses for analytics
type AnalyticsMock struct{}

// RegisterMocks registers all success mock responses
func (m *AnalyticsMock) RegisterMocks(baseURL string) {
	// Mock GET /analytics.json
	httpmock.RegisterResponder("GET", baseURL+"/analytics.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_get_analytics.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	// Mock GET /analytics.csv
	httpmock.RegisterResponder("GET", baseURL+"/analytics.csv",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_get_analytics.csv")
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
			resp := httpmock.NewBytesResponse(401, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	// Mock unauthorized error for CSV
	httpmock.RegisterResponder("GET", baseURL+"/analytics.csv",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_unauthorized.json")
			if err != nil {
				return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid API key"]}`), nil
			}
			resp := httpmock.NewBytesResponse(401, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// CleanupMockState cleans up any state from the mock
func (m *AnalyticsMock) CleanupMockState() {
	// No stateful data to cleanup for analytics
}
