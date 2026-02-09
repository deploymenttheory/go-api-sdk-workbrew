package events

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// EventsServiceInterface defines the interface for events operations
	EventsServiceInterface interface {
		GetEvents(ctx context.Context, opts *RequestQueryOptions) (*EventsResponse, error)
		GetEventsCSV(ctx context.Context, opts *RequestQueryOptions) ([]byte, error)
	}

	// Service handles communication with the events
	// related methods of the Workbrew API.
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements EventsServiceInterface
var _ EventsServiceInterface = (*Service)(nil)

// NewService creates a new events service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// GetEvents retrieves all events in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/events.json
//
// Parameters:
//   - opts: Optional query parameters (filter by actor type: user, system, all)
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/events.json?filter=user"
func (s *Service) GetEvents(ctx context.Context, opts *RequestQueryOptions) (*EventsResponse, error) {
	endpoint := EndpointEventsJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	// Nil-safe options handling
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	// Build query parameters
	queryParams := s.client.QueryBuilder().
		AddIfNotEmpty("filter", opts.Filter).
		Build()

	var result EventsResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetEventsCSV retrieves all events in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/events.csv
//
// Parameters:
//   - opts: Optional query parameters (filter by actor type, download flag)
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/events.csv?filter=user&download=1"
func (s *Service) GetEventsCSV(ctx context.Context, opts *RequestQueryOptions) ([]byte, error) {
	endpoint := EndpointEventsCSV

	headers := map[string]string{
		"Accept": "text/csv",
	}

	// Nil-safe options handling
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	// Build query parameters
	qb := s.client.QueryBuilder().AddIfNotEmpty("filter", opts.Filter)
	if opts.Download {
		qb.AddString("download", "1")
	}
	queryParams := qb.Build()

	csvData, err := s.client.GetCSV(ctx, endpoint, queryParams, headers)
	if err != nil {
		return nil, err
	}

	return csvData, nil
}
