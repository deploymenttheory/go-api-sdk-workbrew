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

// VulnerabilitiesMock handles mock HTTP responses for vulnerabilities
type VulnerabilitiesMock struct{}

// RegisterMocks registers all success mock responses
func (m *VulnerabilitiesMock) RegisterMocks(baseURL string) {
	// Mock GET /vulnerabilities.json
	httpmock.RegisterResponder("GET", baseURL+"/vulnerabilities.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_get_vulnerabilities.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	// Mock GET /vulnerabilities.csv
	httpmock.RegisterResponder("GET", baseURL+"/vulnerabilities.csv",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_get_vulnerabilities.csv")
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
func (m *VulnerabilitiesMock) RegisterErrorMocks(baseURL string) {
	// Mock unauthorized error for JSON
	httpmock.RegisterResponder("GET", baseURL+"/vulnerabilities.json",
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
	httpmock.RegisterResponder("GET", baseURL+"/vulnerabilities.csv",
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

// RegisterForbiddenMocks registers 403 forbidden error mocks (Free tier)
func (m *VulnerabilitiesMock) RegisterForbiddenMocks(baseURL string) {
	httpmock.RegisterResponder("GET", baseURL+"/vulnerabilities.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_forbidden.json")
			if err != nil {
				return httpmock.NewStringResponse(403, `{"message":"Forbidden","errors":["Free tier"]}`), nil
			}
			resp := httpmock.NewBytesResponse(403, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// CleanupMockState cleans up any state from the mock
func (m *VulnerabilitiesMock) CleanupMockState() {
	// No cleanup needed for stateless mocks
}
