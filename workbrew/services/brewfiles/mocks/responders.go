package mocks

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jarcoal/httpmock"
)

// BrewfilesMock handles mock HTTP responses for brewfiles
type BrewfilesMock struct{}

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
func (m *BrewfilesMock) RegisterMocks(baseURL string) {
	// Mock GET /brewfiles.json
	httpmock.RegisterResponder("GET", baseURL+"/brewfiles.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("brewfiles_list.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			return httpmock.NewBytesResponse(200, mockData), nil
		},
	)

	// Mock GET /brewfiles.csv
	httpmock.RegisterResponder("GET", baseURL+"/brewfiles.csv",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("brewfiles_list.csv")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "text/csv")
			return resp, nil
		},
	)

	// Mock POST /brewfiles.json
	httpmock.RegisterResponder("POST", baseURL+"/brewfiles.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("brewfile_created.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			return httpmock.NewBytesResponse(201, mockData), nil
		},
	)

	// Mock PUT /brewfiles/{label}.json - match any label
	httpmock.RegisterResponder("PUT", `=~^`+baseURL+`/brewfiles/[^/]+\.json$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("brewfile_updated.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			return httpmock.NewBytesResponse(200, mockData), nil
		},
	)

	// Mock DELETE /brewfiles/{label}.json - match any label
	httpmock.RegisterResponder("DELETE", `=~^`+baseURL+`/brewfiles/[^/]+\.json$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("brewfile_deleted.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			return httpmock.NewBytesResponse(200, mockData), nil
		},
	)

	// Mock GET /brewfiles/{label}/runs.json - match any label
	httpmock.RegisterResponder("GET", `=~^`+baseURL+`/brewfiles/[^/]+/runs\.json$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("brewfile_runs_list.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			return httpmock.NewBytesResponse(200, mockData), nil
		},
	)

	// Mock GET /brewfiles/{label}/runs.csv - match any label
	httpmock.RegisterResponder("GET", `=~^`+baseURL+`/brewfiles/[^/]+/runs\.csv$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("brewfile_runs_list.csv")
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
func (m *BrewfilesMock) RegisterErrorMocks(baseURL string) {
	// Mock unauthorized errors for all endpoints
	endpoints := []string{
		baseURL + "/brewfiles.json",
		baseURL + "/brewfiles.csv",
	}

	for _, endpoint := range endpoints {
		httpmock.RegisterResponder("GET", endpoint,
			func(req *http.Request) (*http.Response, error) {
				mockData, err := loadMockResponse("error_unauthorized.json")
				if err != nil {
					return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid API key"]}`), nil
				}
				return httpmock.NewBytesResponse(401, mockData), nil
			},
		)
	}

	// Mock unauthorized for POST
	httpmock.RegisterResponder("POST", baseURL+"/brewfiles.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_unauthorized.json")
			if err != nil {
				return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid API key"]}`), nil
			}
			return httpmock.NewBytesResponse(401, mockData), nil
		},
	)

	// Mock error responses for parameterized endpoints
	httpmock.RegisterResponder("PUT", `=~^`+baseURL+`/brewfiles/[^/]+\.json$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_unauthorized.json")
			if err != nil {
				return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid API key"]}`), nil
			}
			return httpmock.NewBytesResponse(401, mockData), nil
		},
	)

	httpmock.RegisterResponder("DELETE", `=~^`+baseURL+`/brewfiles/[^/]+\.json$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_unauthorized.json")
			if err != nil {
				return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid API key"]}`), nil
			}
			return httpmock.NewBytesResponse(401, mockData), nil
		},
	)

	httpmock.RegisterResponder("GET", `=~^`+baseURL+`/brewfiles/[^/]+/runs\.json$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_unauthorized.json")
			if err != nil {
				return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid API key"]}`), nil
			}
			return httpmock.NewBytesResponse(401, mockData), nil
		},
	)

	httpmock.RegisterResponder("GET", `=~^`+baseURL+`/brewfiles/[^/]+/runs\.csv$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_unauthorized.json")
			if err != nil {
				return httpmock.NewStringResponse(401, `{"message":"Unauthorized","errors":["Invalid API key"]}`), nil
			}
			return httpmock.NewBytesResponse(401, mockData), nil
		},
	)
}

// RegisterForbiddenMocks registers 403 forbidden error mocks (Free tier)
func (m *BrewfilesMock) RegisterForbiddenMocks(baseURL string) {
	httpmock.RegisterResponder("POST", baseURL+"/brewfiles.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_forbidden.json")
			if err != nil {
				return httpmock.NewStringResponse(403, `{"message":"Forbidden","errors":["Free tier"]}`), nil
			}
			return httpmock.NewBytesResponse(403, mockData), nil
		},
	)
}

// RegisterValidationMocks registers 422 validation error mocks
func (m *BrewfilesMock) RegisterValidationMocks(baseURL string) {
	httpmock.RegisterResponder("PUT", `=~^`+baseURL+`/brewfiles/[^/]+\.json$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_validation.json")
			if err != nil {
				return httpmock.NewStringResponse(422, `{"message":"Validation error","errors":["Invalid content"]}`), nil
			}
			return httpmock.NewBytesResponse(422, mockData), nil
		},
	)
}

// CleanupMockState cleans up any state from the mock
func (m *BrewfilesMock) CleanupMockState() {
	// No cleanup needed for stateless mocks
}
