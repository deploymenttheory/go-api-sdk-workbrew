package brewconfigurations

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// BrewConfigurationsServiceInterface defines the interface for brew configurations operations
	BrewConfigurationsServiceInterface interface {
		// GetBrewConfigurations retrieves all brew configurations in JSON format
		GetBrewConfigurations(ctx context.Context) (*BrewConfigurationsResponse, error)

		// GetBrewConfigurationsCSV retrieves all brew configurations in CSV format
		GetBrewConfigurationsCSV(ctx context.Context) ([]byte, error)
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

// GetBrewConfigurations retrieves all brew configurations in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_configurations.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brew_configurations.json"
func (s *Service) GetBrewConfigurations(ctx context.Context) (*BrewConfigurationsResponse, error) {
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

// GetBrewConfigurationsCSV retrieves all brew configurations in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_configurations.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brew_configurations.csv"
func (s *Service) GetBrewConfigurationsCSV(ctx context.Context) ([]byte, error) {
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
