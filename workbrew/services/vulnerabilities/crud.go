package vulnerabilities

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// VulnerabilitiesServiceInterface defines the interface for vulnerabilities operations
	VulnerabilitiesServiceInterface interface {
		GetVulnerabilities(ctx context.Context) (*VulnerabilitiesResponse, error)
		GetVulnerabilitiesCSV(ctx context.Context) ([]byte, error)
	}

	// Service handles communication with the vulnerabilities
	// related methods of the Workbrew API.
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements VulnerabilitiesServiceInterface
var _ VulnerabilitiesServiceInterface = (*Service)(nil)

// NewService creates a new vulnerabilities service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// GetVulnerabilities retrieves all vulnerabilities in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/vulnerabilities.json
//
// Note: This endpoint may return 403 on Free tier plans
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/vulnerabilities.json"
func (s *Service) GetVulnerabilities(ctx context.Context) (*VulnerabilitiesResponse, error) {
	endpoint := EndpointVulnerabilitiesJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result VulnerabilitiesResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetVulnerabilitiesCSV retrieves all vulnerabilities in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/vulnerabilities.csv
//
// Note: This endpoint may return 403 on Free tier plans
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/vulnerabilities.csv"
func (s *Service) GetVulnerabilitiesCSV(ctx context.Context) ([]byte, error) {
	endpoint := EndpointVulnerabilitiesCSV

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
