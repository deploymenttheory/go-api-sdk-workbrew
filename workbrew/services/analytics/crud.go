package analytics

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// AnalyticsServiceInterface defines the interface for analytics operations
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	AnalyticsServiceInterface interface {
		// ListAnalytics returns a list of analytics data showing command usage statistics per device
		//
		// Returns analytics records with device, command, last run timestamp, and count information
		ListAnalytics(ctx context.Context) (*AnalyticsResponse, error)

		// ListAnalyticsCSV returns a list of analytics data in CSV format
		//
		// Returns the same analytics data as ListAnalytics but formatted as CSV
		ListAnalyticsCSV(ctx context.Context) ([]byte, error)
	}

	// Service handles communication with the analytics
	// related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements AnalyticsServiceInterface
var _ AnalyticsServiceInterface = (*Service)(nil)

// NewService creates a new analytics service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// ListAnalytics retrieves all analytics in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/analytics.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/analytics.json"
func (s *Service) ListAnalytics(ctx context.Context) (*AnalyticsResponse, error) {
	endpoint := EndpointAnalyticsJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result AnalyticsResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListAnalyticsCSV retrieves all analytics in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/analytics.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/analytics.csv"
func (s *Service) ListAnalyticsCSV(ctx context.Context) ([]byte, error) {
	endpoint := EndpointAnalyticsCSV

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
