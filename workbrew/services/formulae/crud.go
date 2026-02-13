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
		ListFormulae(ctx context.Context) (*FormulaeResponse, *interfaces.Response, error)

		// ListFormulaeCSV returns a list of Formulae in CSV format
		//
		// Returns formulae data as CSV with columns: name, devices, outdated, installed_on_request, installed_as_dependency,
		// vulnerabilities, deprecated, license, homebrew_core_version.
		ListFormulaeCSV(ctx context.Context) ([]byte, *interfaces.Response, error)
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
func (s *Service) ListFormulae(ctx context.Context) (*FormulaeResponse, *interfaces.Response, error) {
	endpoint := EndpointFormulaeJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result FormulaeResponse
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListFormulaeCSV retrieves all formulae in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/formulae.csv
func (s *Service) ListFormulaeCSV(ctx context.Context) ([]byte, *interfaces.Response, error) {
	endpoint := EndpointFormulaeCSV

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
