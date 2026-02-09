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

// EventsMock handles mock HTTP responses for events
type EventsMock struct{}

// RegisterMocks registers all success mock responses
func (m *EventsMock) RegisterMocks(baseURL string) {
	// Mock GET /events.json
	httpmock.RegisterResponder("GET", baseURL+"/events.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_get_events.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	// Mock GET /events.csv
	httpmock.RegisterResponder("GET", baseURL+"/events.csv",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_get_events.csv")
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
func (m *EventsMock) RegisterErrorMocks(baseURL string) {
	// Mock unauthorized error
	httpmock.RegisterResponder("GET", baseURL+"/events.json",
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

	httpmock.RegisterResponder("GET", baseURL+"/events.csv",
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
func (m *EventsMock) CleanupMockState() {
	// No cleanup needed for stateless mocks
}
