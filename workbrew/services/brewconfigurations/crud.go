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
		ListBrewConfigurations(ctx context.Context) (*BrewConfigurationsResponse, *interfaces.Response, error)

		// ListBrewConfigurationsCSV returns a list of Brew Configurations in CSV format
		//
		// Returns brew configuration data as CSV with columns: key, value, last_updated_by_user, device_group.
		ListBrewConfigurationsCSV(ctx context.Context) ([]byte, *interfaces.Response, error)
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
func (s *Service) ListBrewConfigurations(ctx context.Context) (*BrewConfigurationsResponse, *interfaces.Response, error) {
	endpoint := EndpointBrewConfigurationsJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewConfigurationsResponse
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListBrewConfigurationsCSV retrieves all brew configurations in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_configurations.csv
func (s *Service) ListBrewConfigurationsCSV(ctx context.Context) ([]byte, *interfaces.Response, error) {
	endpoint := EndpointBrewConfigurationsCSV

	headers := map[string]string{
		"Accept": "text/csv",
	}

	queryParams := make(map[string]string)

	resp, csvData, err := s.client.GetBytes(ctx, endpoint, queryParams, headers)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
