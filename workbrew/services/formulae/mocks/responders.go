package mocks

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jarcoal/httpmock"
)

// FormulaeMock handles mock HTTP responses for formulae
type FormulaeMock struct{}

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
func (m *FormulaeMock) RegisterMocks(baseURL string) {
	// Mock GET /formulae.json
	httpmock.RegisterResponder("GET", "https://console.workbrew.com/workspaces/test-workspace/formulae.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_get_formulae.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			return httpmock.NewBytesResponse(200, mockData), nil
		},
	)

	// Mock GET /formulae.csv
	httpmock.RegisterResponder("GET", "https://console.workbrew.com/workspaces/test-workspace/formulae.csv",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_get_formulae.csv")
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
func (m *FormulaeMock) RegisterErrorMocks(baseURL string) {
	// Mock unauthorized error for JSON
	httpmock.RegisterResponder("GET", "https://console.workbrew.com/workspaces/test-workspace/formulae.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_unauthorized.json")
			if err != nil {
				return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid API key"]}`), nil
			}
			return httpmock.NewBytesResponse(401, mockData), nil
		},
	)

	// Mock unauthorized error for CSV
	httpmock.RegisterResponder("GET", "https://console.workbrew.com/workspaces/test-workspace/formulae.csv",
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
func (m *FormulaeMock) CleanupMockState() {
	// No cleanup needed for stateless mocks
}
