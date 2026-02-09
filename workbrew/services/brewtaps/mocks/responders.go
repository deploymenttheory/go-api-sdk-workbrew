package mocks

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jarcoal/httpmock"
)

// BrewTapsMock handles mock HTTP responses for brew taps
type BrewTapsMock struct{}

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
func (m *BrewTapsMock) RegisterMocks(baseURL string) {
	// Mock GET /brew_taps.json
	httpmock.RegisterResponder("GET", baseURL+"/brew_taps.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("brewtaps_list.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			return httpmock.NewBytesResponse(200, mockData), nil
		},
	)

	// Mock GET /brew_taps.csv
	httpmock.RegisterResponder("GET", baseURL+"/brew_taps.csv",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("brewtaps_list.csv")
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
func (m *BrewTapsMock) RegisterErrorMocks(baseURL string) {
	// Mock unauthorized error for JSON
	httpmock.RegisterResponder("GET", baseURL+"/brew_taps.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_unauthorized.json")
			if err != nil {
				return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid API key"]}`), nil
			}
			return httpmock.NewBytesResponse(401, mockData), nil
		},
	)

	// Mock unauthorized error for CSV
	httpmock.RegisterResponder("GET", baseURL+"/brew_taps.csv",
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
func (m *BrewTapsMock) CleanupMockState() {
	// No cleanup needed for stateless mocks
}
