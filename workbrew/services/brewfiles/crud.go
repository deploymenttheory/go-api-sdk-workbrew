package brewfiles

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// BrewfilesServiceInterface defines the interface for brewfiles operations
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	BrewfilesServiceInterface interface {
		// ListBrewfiles returns a list of brewfiles with their status and device assignments
		//
		// Returns brewfiles with last updated user, start/finish timestamps, assigned devices, and run count
		ListBrewfiles(ctx context.Context) (*BrewfilesResponse, *interfaces.Response, error)

		// ListBrewfilesCSV returns a list of brewfiles in CSV format
		//
		// Returns the same brewfiles data as ListBrewfiles but formatted as CSV
		ListBrewfilesCSV(ctx context.Context) ([]byte, *interfaces.Response, error)

		// CreateBrewfile creates a new brewfile with specified label, content, and device/group assignment
		//
		// Requires label and content fields. Can assign to specific devices via device_serial_numbers or to a device group via device_group_id
		CreateBrewfile(ctx context.Context, request *CreateBrewfileRequest) (*BrewfileMessageResponse, *interfaces.Response, error)

		// UpdateBrewfile updates an existing brewfile's content and device assignments
		//
		// Updates the brewfile identified by label. Can update content, device_serial_numbers, or device_group_id
		UpdateBrewfile(ctx context.Context, label string, request *UpdateBrewfileRequest) (*BrewfileMessageResponse, *interfaces.Response, error)

		// DeleteBrewfile deletes a brewfile by its label
		//
		// Permanently removes the brewfile identified by the specified label
		DeleteBrewfile(ctx context.Context, label string) (*BrewfileMessageResponse, *interfaces.Response, error)

		// ListBrewfileRuns returns a list of brewfile runs for a specific brewfile
		//
		// Returns run history including label, device, timestamps, success status, and output for the specified brewfile label
		ListBrewfileRuns(ctx context.Context, label string) (*BrewfileRunsResponse, *interfaces.Response, error)

		// ListBrewfileRunsCSV returns a list of brewfile runs in CSV format
		//
		// Returns the same run data as ListBrewfileRuns but formatted as CSV
		ListBrewfileRunsCSV(ctx context.Context, label string) ([]byte, *interfaces.Response, error)
	}

	// Service handles communication with the brewfiles
	// related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements BrewfilesServiceInterface
var _ BrewfilesServiceInterface = (*Service)(nil)

// NewService creates a new brewfiles service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// ListBrewfiles retrieves all brewfiles in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles.json
func (s *Service) ListBrewfiles(ctx context.Context) (*BrewfilesResponse, *interfaces.Response, error) {
	endpoint := EndpointBrewfilesJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewfilesResponse
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListBrewfilesCSV retrieves all brewfiles in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles.csv
func (s *Service) ListBrewfilesCSV(ctx context.Context) ([]byte, *interfaces.Response, error) {
	endpoint := EndpointBrewfilesCSV

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

// CreateBrewfile creates a new brewfile
// URL: POST https://console.workbrew.com/workspaces/{workspace_name}/brewfiles.json
func (s *Service) CreateBrewfile(ctx context.Context, request *CreateBrewfileRequest) (*BrewfileMessageResponse, *interfaces.Response, error) {
	endpoint := EndpointBrewfilesJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result BrewfileMessageResponse
	resp, err := s.client.Post(ctx, endpoint, request, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// UpdateBrewfile updates an existing brewfile
// URL: PUT https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}.json
//
// Response codes:
//   - 200: Brewfile updated successfully
//   - 422: Validation error
func (s *Service) UpdateBrewfile(ctx context.Context, label string, request *UpdateBrewfileRequest) (*BrewfileMessageResponse, *interfaces.Response, error) {
	if label == "" {
		return nil, nil, fmt.Errorf("brewfile label is required")
	}

	endpoint := fmt.Sprintf(EndpointBrewfileLabelFormat, label)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result BrewfileMessageResponse
	resp, err := s.client.Put(ctx, endpoint, request, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// DeleteBrewfile deletes a brewfile
// URL: DELETE https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}.json
//
// Response codes:
//   - 200: Brewfile deleted successfully
func (s *Service) DeleteBrewfile(ctx context.Context, label string) (*BrewfileMessageResponse, *interfaces.Response, error) {
	if label == "" {
		return nil, nil, fmt.Errorf("brewfile label is required")
	}

	endpoint := fmt.Sprintf(EndpointBrewfileLabelFormat, label)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewfileMessageResponse
	resp, err := s.client.Delete(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListBrewfileRuns retrieves all runs for a specific brewfile in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}/runs.json
func (s *Service) ListBrewfileRuns(ctx context.Context, label string) (*BrewfileRunsResponse, *interfaces.Response, error) {
	if label == "" {
		return nil, nil, fmt.Errorf("brewfile label is required")
	}

	endpoint := fmt.Sprintf(EndpointBrewfileRunsJSONFormat, label)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewfileRunsResponse
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListBrewfileRunsCSV retrieves all runs for a specific brewfile in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}/runs.csv
func (s *Service) ListBrewfileRunsCSV(ctx context.Context, label string) ([]byte, *interfaces.Response, error) {
	if label == "" {
		return nil, nil, fmt.Errorf("brewfile label is required")
	}

	endpoint := fmt.Sprintf(EndpointBrewfileRunsCSVFormat, label)

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
