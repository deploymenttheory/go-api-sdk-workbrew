package brewfiles

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// BrewfilesServiceInterface defines the interface for brewfiles operations
	BrewfilesServiceInterface interface {
		GetBrewfiles(ctx context.Context) (*BrewfilesResponse, error)
		GetBrewfilesCSV(ctx context.Context) ([]byte, error)
		CreateBrewfile(ctx context.Context, request *CreateBrewfileRequest) (*BrewfileMessageResponse, error)
		UpdateBrewfile(ctx context.Context, label string, request *UpdateBrewfileRequest) (*BrewfileMessageResponse, error)
		DeleteBrewfile(ctx context.Context, label string) (*BrewfileMessageResponse, error)
		GetBrewfileRuns(ctx context.Context, label string) (*BrewfileRunsResponse, error)
		GetBrewfileRunsCSV(ctx context.Context, label string) ([]byte, error)
	}

	// Service handles communication with the brewfiles
	// related methods of the Workbrew API.
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements BrewfilesServiceInterface
var _ BrewfilesServiceInterface = (*Service)(nil)

// NewService creates a new brewfiles service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// GetBrewfiles retrieves all brewfiles in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brewfiles.json"
func (s *Service) GetBrewfiles(ctx context.Context) (*BrewfilesResponse, error) {
	endpoint := EndpointBrewfilesJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewfilesResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBrewfilesCSV retrieves all brewfiles in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brewfiles.csv"
func (s *Service) GetBrewfilesCSV(ctx context.Context) ([]byte, error) {
	endpoint := EndpointBrewfilesCSV

	headers := map[string]string{
		"Accept": "text/csv",
	}

	queryParams := make(map[string]string)

	csvData, err := s.client.GetCSV(ctx, endpoint, queryParams, headers)
	if err != nil {
		return nil, err
	}

	return csvData, nil
}

// CreateBrewfile creates a new brewfile
// URL: POST https://console.workbrew.com/workspaces/{workspace_name}/brewfiles.json
//
// Response codes:
//   - 201: Brewfile created successfully
//   - 422: Validation error
//
// Example cURL:
//
//	curl -X POST \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  -H "Content-Type: application/json" \
//	  -d '{"label":"my-brewfile","content":"brew \"wget\""}' \
//	  "https://console.workbrew.com/workspaces/{workspace}/brewfiles.json"
func (s *Service) CreateBrewfile(ctx context.Context, request *CreateBrewfileRequest) (*BrewfileMessageResponse, error) {
	endpoint := EndpointBrewfilesJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result BrewfileMessageResponse
	err := s.client.Post(ctx, endpoint, request, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateBrewfile updates an existing brewfile
// URL: PUT https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}.json
//
// Response codes:
//   - 200: Brewfile updated successfully
//   - 422: Validation error
//
// Example cURL:
//
//	curl -X PUT \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  -H "Content-Type: application/json" \
//	  -d '{"content":"brew \"wget\"\nbrew \"htop\""}' \
//	  "https://console.workbrew.com/workspaces/{workspace}/brewfiles/my-brewfile.json"
func (s *Service) UpdateBrewfile(ctx context.Context, label string, request *UpdateBrewfileRequest) (*BrewfileMessageResponse, error) {
	if label == "" {
		return nil, fmt.Errorf("brewfile label is required")
	}

	endpoint := fmt.Sprintf(EndpointBrewfileLabelFormat, label)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result BrewfileMessageResponse
	err := s.client.Put(ctx, endpoint, request, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteBrewfile deletes a brewfile
// URL: DELETE https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}.json
//
// Response codes:
//   - 200: Brewfile deleted successfully
//
// Example cURL:
//
//	curl -X DELETE \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brewfiles/my-brewfile.json"
func (s *Service) DeleteBrewfile(ctx context.Context, label string) (*BrewfileMessageResponse, error) {
	if label == "" {
		return nil, fmt.Errorf("brewfile label is required")
	}

	endpoint := fmt.Sprintf(EndpointBrewfileLabelFormat, label)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewfileMessageResponse
	err := s.client.Delete(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBrewfileRuns retrieves all runs for a specific brewfile in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}/runs.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brewfiles/my-brewfile/runs.json"
func (s *Service) GetBrewfileRuns(ctx context.Context, label string) (*BrewfileRunsResponse, error) {
	if label == "" {
		return nil, fmt.Errorf("brewfile label is required")
	}

	endpoint := fmt.Sprintf(EndpointBrewfileRunsJSONFormat, label)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewfileRunsResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBrewfileRunsCSV retrieves all runs for a specific brewfile in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}/runs.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brewfiles/my-brewfile/runs.csv"
func (s *Service) GetBrewfileRunsCSV(ctx context.Context, label string) ([]byte, error) {
	if label == "" {
		return nil, fmt.Errorf("brewfile label is required")
	}

	endpoint := fmt.Sprintf(EndpointBrewfileRunsCSVFormat, label)

	headers := map[string]string{
		"Accept": "text/csv",
	}

	queryParams := make(map[string]string)

	csvData, err := s.client.GetCSV(ctx, endpoint, queryParams, headers)
	if err != nil {
		return nil, err
	}

	return csvData, nil
}
