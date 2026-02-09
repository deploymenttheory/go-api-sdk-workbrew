package brewtaps

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// BrewTapsServiceInterface defines the interface for brew taps operations
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	BrewTapsServiceInterface interface {
		// ListBrewTaps returns a list of Taps
		//
		// Returns Homebrew taps with their names, assigned devices, installed formulae/casks counts, and available packages.
		ListBrewTaps(ctx context.Context) (*BrewTapsResponse, error)

		// ListBrewTapsCSV returns a list of Taps in CSV format
		//
		// Returns tap data as CSV with columns: tap, devices, formulae_installed, casks_installed, available_packages.
		ListBrewTapsCSV(ctx context.Context) ([]byte, error)
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

// ListBrewTaps retrieves all brew taps in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_taps.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brew_taps.json"
func (s *Service) ListBrewTaps(ctx context.Context) (*BrewTapsResponse, error) {
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

// ListBrewTapsCSV retrieves all brew taps in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_taps.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brew_taps.csv"
func (s *Service) ListBrewTapsCSV(ctx context.Context) ([]byte, error) {
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
