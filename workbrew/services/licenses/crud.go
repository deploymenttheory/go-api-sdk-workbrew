package licenses

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// LicensesServiceInterface defines the interface for licenses operations
	LicensesServiceInterface interface {
		GetLicenses(ctx context.Context) (*LicensesResponse, error)
		GetLicensesCSV(ctx context.Context) ([]byte, error)
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

// GetLicenses retrieves all licenses in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/licenses.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/licenses.json"
func (s *Service) GetLicenses(ctx context.Context) (*LicensesResponse, error) {
	endpoint := EndpointLicensesJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result LicensesResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetLicensesCSV retrieves all licenses in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/licenses.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/licenses.csv"
func (s *Service) GetLicensesCSV(ctx context.Context) ([]byte, error) {
	endpoint := EndpointLicensesCSV

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
