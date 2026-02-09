package mocks

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jarcoal/httpmock"
)

// EventsMock handles mock HTTP responses for events
type EventsMock struct{}

// loadMockResponse loads a mock response file from the mocks directory
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	data, err := os.ReadFile(mockPath)
	if err != nil {
		absPath, _ := filepath.Abs(mockPath)
		return nil, fmt.Errorf("failed to load mock file %s (tried: %s, %s): %w", filename, mockPath, absPath, err)
	}
	return data, nil
}

// RegisterMocks registers all success mock responses
func (m *EventsMock) RegisterMocks(baseURL string) {
	// Mock GET /events.json
	httpmock.RegisterResponder("GET", baseURL+"/events.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("events_list.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			return httpmock.NewBytesResponse(200, mockData), nil
		},
	)

	// Mock GET /events.csv
	httpmock.RegisterResponder("GET", baseURL+"/events.csv",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("events_list.csv")
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
			return httpmock.NewBytesResponse(401, mockData), nil
		},
	)

	httpmock.RegisterResponder("GET", baseURL+"/events.csv",
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
func (m *EventsMock) CleanupMockState() {
	// No cleanup needed for stateless mocks
}
