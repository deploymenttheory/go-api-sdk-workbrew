package casks

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// CasksServiceInterface defines the interface for casks operations
	CasksServiceInterface interface {
		GetCasks(ctx context.Context) (*CasksResponse, error)
		GetCasksCSV(ctx context.Context) ([]byte, error)
	}

	// Service handles communication with the casks
	// related methods of the Workbrew API.
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements CasksServiceInterface
var _ CasksServiceInterface = (*Service)(nil)

// NewService creates a new casks service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// GetCasks retrieves all casks in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/casks.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/casks.json"
func (s *Service) GetCasks(ctx context.Context) (*CasksResponse, error) {
	endpoint := EndpointCasksJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result CasksResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCasksCSV retrieves all casks in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/casks.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/casks.csv"
func (s *Service) GetCasksCSV(ctx context.Context) ([]byte, error) {
	endpoint := EndpointCasksCSV

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
