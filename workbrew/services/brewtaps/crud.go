package brewtaps

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// BrewTapsServiceInterface defines the interface for brew taps operations
	BrewTapsServiceInterface interface {
		GetBrewTaps(ctx context.Context) (*BrewTapsResponse, error)
		GetBrewTapsCSV(ctx context.Context) ([]byte, error)
	}

	// Service handles communication with the brew taps
	// related methods of the Workbrew API.
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements BrewTapsServiceInterface
var _ BrewTapsServiceInterface = (*Service)(nil)

// NewService creates a new brew taps service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// GetBrewTaps retrieves all brew taps in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_taps.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brew_taps.json"
func (s *Service) GetBrewTaps(ctx context.Context) (*BrewTapsResponse, error) {
	endpoint := EndpointBrewTapsJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewTapsResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBrewTapsCSV retrieves all brew taps in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_taps.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brew_taps.csv"
func (s *Service) GetBrewTapsCSV(ctx context.Context) ([]byte, error) {
	endpoint := EndpointBrewTapsCSV

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
