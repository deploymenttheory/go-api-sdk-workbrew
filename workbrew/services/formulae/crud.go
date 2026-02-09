package formulae

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// FormulaeServiceInterface defines the interface for formulae operations
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	FormulaeServiceInterface interface {
		// ListFormulae returns a list of Formulae
		//
		// Returns installed Homebrew formulae with names, assigned devices, outdated status, installation type (on request/dependency),
		// known vulnerabilities, deprecation status, licenses, and Homebrew core versions.
		ListFormulae(ctx context.Context) (*FormulaeResponse, error)

		// ListFormulaeCSV returns a list of Formulae in CSV format
		//
		// Returns formulae data as CSV with columns: name, devices, outdated, installed_on_request, installed_as_dependency, 
		// vulnerabilities, deprecated, license, homebrew_core_version.
		ListFormulaeCSV(ctx context.Context) ([]byte, error)
	}

	// Service handles communication with the formulae
	// related methods of the Workbrew API.
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements FormulaeServiceInterface
var _ FormulaeServiceInterface = (*Service)(nil)

// NewService creates a new formulae service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// ListFormulae retrieves all formulae in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/formulae.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/formulae.json"
func (s *Service) ListFormulae(ctx context.Context) (*FormulaeResponse, error) {
	endpoint := EndpointFormulaeJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result FormulaeResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListFormulaeCSV retrieves all formulae in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/formulae.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/formulae.csv"
func (s *Service) ListFormulaeCSV(ctx context.Context) ([]byte, error) {
	endpoint := EndpointFormulaeCSV

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
