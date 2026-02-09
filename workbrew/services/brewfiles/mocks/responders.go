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

// BrewfilesMock handles mock HTTP responses for brewfiles
type BrewfilesMock struct{}

// RegisterMocks registers all success mock responses
func (m *BrewfilesMock) RegisterMocks(baseURL string) {
	// Mock GET /brewfiles.json
	httpmock.RegisterResponder("GET", baseURL+"/brewfiles.json",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_get_brewfiles.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	// Mock GET /brewfiles.csv
	httpmock.RegisterResponder("GET", baseURL+"/brewfiles.csv",
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_get_brewfiles.csv")
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
			mockData, err := loadMockResponse("validate_create_brewfile.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			resp := httpmock.NewBytesResponse(201, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	// Mock PUT /brewfiles/{label}.json - match any label
	httpmock.RegisterResponder("PUT", `=~^`+baseURL+`/brewfiles/[^/]+\.json$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_update_brewfile.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	// Mock DELETE /brewfiles/{label}.json - match any label
	httpmock.RegisterResponder("DELETE", `=~^`+baseURL+`/brewfiles/[^/]+\.json$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_delete_brewfile.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	// Mock GET /brewfiles/{label}/runs.json - match any label
	httpmock.RegisterResponder("GET", `=~^`+baseURL+`/brewfiles/[^/]+/runs\.json$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_get_brewfile_runs.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"message":"Internal Server Error","errors":["Failed to load mock data"]}`), nil
			}
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)

	// Mock GET /brewfiles/{label}/runs.csv - match any label
	httpmock.RegisterResponder("GET", `=~^`+baseURL+`/brewfiles/[^/]+/runs\.csv$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("validate_get_brewfile_runs.csv")
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
	// Mock unauthorized errors for all GET endpoints
	httpmock.RegisterResponder("GET", baseURL+"/brewfiles.json",
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

	httpmock.RegisterResponder("GET", baseURL+"/brewfiles.csv",
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

	// Mock unauthorized for POST
	httpmock.RegisterResponder("POST", baseURL+"/brewfiles.json",
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

	// Mock error responses for parameterized endpoints
	httpmock.RegisterResponder("PUT", `=~^`+baseURL+`/brewfiles/[^/]+\.json$`,
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

	httpmock.RegisterResponder("DELETE", `=~^`+baseURL+`/brewfiles/[^/]+\.json$`,
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

	httpmock.RegisterResponder("GET", `=~^`+baseURL+`/brewfiles/[^/]+/runs\.json$`,
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

	httpmock.RegisterResponder("GET", `=~^`+baseURL+`/brewfiles/[^/]+/runs\.csv$`,
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
func (m *BrewfilesMock) RegisterForbiddenMocks(baseURL string) {
	httpmock.RegisterResponder("POST", baseURL+"/brewfiles.json",
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

// RegisterValidationMocks registers 422 validation error mocks
func (m *BrewfilesMock) RegisterValidationMocks(baseURL string) {
	httpmock.RegisterResponder("PUT", `=~^`+baseURL+`/brewfiles/[^/]+\.json$`,
		func(req *http.Request) (*http.Response, error) {
			mockData, err := loadMockResponse("error_validation.json")
			if err != nil {
				return httpmock.NewStringResponse(422, `{"message":"Validation error","errors":["Invalid content"]}`), nil
			}
			resp := httpmock.NewBytesResponse(422, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// CleanupMockState cleans up any state from the mock
func (m *BrewfilesMock) CleanupMockState() {
	// No cleanup needed for stateless mocks
}
