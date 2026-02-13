package vulnerabilities

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// VulnerabilitiesServiceInterface defines the interface for vulnerabilities operations
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	VulnerabilitiesServiceInterface interface {
		// ListVulnerabilities returns a list of Vulnerabilities
		//
		// Returns security vulnerabilities affecting installed formulae, including CVE IDs with CVSS scores,
		// affected formula names, outdated devices, support status, and Homebrew core versions.
		// May return 403 Forbidden on Free tier plans.
		ListVulnerabilities(ctx context.Context) (*VulnerabilitiesResponse, *interfaces.Response, error)

		// ListVulnerabilitiesCSV returns a list of Vulnerabilities in CSV format
		//
		// Returns vulnerability data as CSV with columns: vulnerabilities, formula, outdated_devices, supported, homebrew_core_version.
		ListVulnerabilitiesCSV(ctx context.Context) ([]byte, *interfaces.Response, error)
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

// ListVulnerabilities retrieves all vulnerabilities in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/vulnerabilities.json
//
// Note: This endpoint may return 403 on Free tier plans
func (s *Service) ListVulnerabilities(ctx context.Context) (*VulnerabilitiesResponse, *interfaces.Response, error) {
	endpoint := EndpointVulnerabilitiesJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result VulnerabilitiesResponse
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListVulnerabilitiesCSV retrieves all vulnerabilities in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/vulnerabilities.csv
//
// Note: This endpoint may return 403 on Free tier plans
func (s *Service) ListVulnerabilitiesCSV(ctx context.Context) ([]byte, *interfaces.Response, error) {
	endpoint := EndpointVulnerabilitiesCSV

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
