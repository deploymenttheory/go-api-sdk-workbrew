package events

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// EventsServiceInterface defines the interface for events operations
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	EventsServiceInterface interface {
		// ListEvents returns a list of audit log events
		//
		// Returns audit log events with IDs, event types, timestamps, actor information, and target details.
		// Supports filtering by actor type (user, system, or all) via query options.
		ListEvents(ctx context.Context, opts *RequestQueryOptions) (*EventsResponse, *interfaces.Response, error)

		// ListEventsCSV returns audit log events as CSV
		//
		// Returns audit log event data as CSV with columns: id, event_type, occurred_at, actor_id, actor_type, target_id, target_type, target_identifier.
		// Supports filtering by actor type and optional download parameter via query options.
		ListEventsCSV(ctx context.Context, opts *RequestQueryOptions) ([]byte, *interfaces.Response, error)
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

// ListEvents retrieves all events in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/events.json
//
// Parameters:
//   - opts: Optional query parameters (filter by actor type: user, system, all)
func (s *Service) ListEvents(ctx context.Context, opts *RequestQueryOptions) (*EventsResponse, *interfaces.Response, error) {
	endpoint := EndpointEventsJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	queryParams := s.client.QueryBuilder().
		AddIfNotEmpty("filter", opts.Filter).
		Build()

	var result EventsResponse
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListEventsCSV retrieves all events in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/events.csv
//
// Parameters:
//   - opts: Optional query parameters (filter by actor type, download flag)
func (s *Service) ListEventsCSV(ctx context.Context, opts *RequestQueryOptions) ([]byte, *interfaces.Response, error) {
	endpoint := EndpointEventsCSV

	headers := map[string]string{
		"Accept": "text/csv",
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	qb := s.client.QueryBuilder().AddIfNotEmpty("filter", opts.Filter)
	if opts.Download {
		qb.AddString("download", "1")
	}
	queryParams := qb.Build()

	resp, csvData, err := s.client.GetCSV(ctx, endpoint, queryParams, headers)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
