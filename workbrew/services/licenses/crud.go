package licenses

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// LicensesServiceInterface defines the interface for licenses operations
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	LicensesServiceInterface interface {
		// ListLicenses returns a list of Licenses
		//
		// Returns software licenses found across installed formulae, with license names and counts of affected devices and formulae.
		ListLicenses(ctx context.Context) (*LicensesResponse, *interfaces.Response, error)

		// ListLicensesCSV returns a list of Licenses in CSV format
		//
		// Returns license data as CSV with columns: name, device_count, formula_count.
		ListLicensesCSV(ctx context.Context) ([]byte, *interfaces.Response, error)
	}

	// Service handles communication with the licenses
	// related methods of the Workbrew API.
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements LicensesServiceInterface
var _ LicensesServiceInterface = (*Service)(nil)

// NewService creates a new licenses service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// ListLicenses retrieves all licenses in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/licenses.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/licenses.json"
func (s *Service) ListLicenses(ctx context.Context) (*LicensesResponse, *interfaces.Response, error) {
	endpoint := EndpointLicensesJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result LicensesResponse
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListLicensesCSV retrieves all licenses in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/licenses.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/licenses.csv"
func (s *Service) ListLicensesCSV(ctx context.Context) ([]byte, *interfaces.Response, error) {
	endpoint := EndpointLicensesCSV

	headers := map[string]string{
		"Accept": "text/csv",
	}

	queryParams := make(map[string]string)

	resp, csvData, err := s.client.GetCSV(ctx, endpoint, queryParams, headers)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
