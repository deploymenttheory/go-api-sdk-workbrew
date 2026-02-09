package brewconfigurations

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// BrewConfigurationsServiceInterface defines the interface for brew configurations operations
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	BrewConfigurationsServiceInterface interface {
		// ListBrewConfigurations returns a list of Brew Configurations
		//
		// Returns Homebrew environment variable configurations with their keys, values, last updated user, and assigned device groups.
		ListBrewConfigurations(ctx context.Context) (*BrewConfigurationsResponse, error)

		// ListBrewConfigurationsCSV returns a list of Brew Configurations in CSV format
		//
		// Returns brew configuration data as CSV with columns: key, value, last_updated_by_user, device_group.
		ListBrewConfigurationsCSV(ctx context.Context) ([]byte, error)
	}

	// Service handles communication with the brew configurations
	// related methods of the Workbrew API.
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements BrewConfigurationsServiceInterface
var _ BrewConfigurationsServiceInterface = (*Service)(nil)

// NewService creates a new brew configurations service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// ListBrewConfigurations retrieves all brew configurations in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_configurations.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brew_configurations.json"
func (s *Service) ListBrewConfigurations(ctx context.Context) (*BrewConfigurationsResponse, error) {
	endpoint := EndpointBrewConfigurationsJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewConfigurationsResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListBrewConfigurationsCSV retrieves all brew configurations in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_configurations.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brew_configurations.csv"
func (s *Service) ListBrewConfigurationsCSV(ctx context.Context) ([]byte, error) {
	endpoint := EndpointBrewConfigurationsCSV

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
