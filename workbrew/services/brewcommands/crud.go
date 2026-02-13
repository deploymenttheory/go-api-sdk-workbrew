package brewcommands

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// BrewCommandsServiceInterface defines the interface for brew commands operations
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	BrewCommandsServiceInterface interface {
		// ListBrewCommands returns a list of brew commands with their configuration and status
		//
		// Returns brew commands with command, label, last updated user, start/finish timestamps, devices, and run count
		ListBrewCommands(ctx context.Context) (*BrewCommandsResponse, *interfaces.Response, error)

		// ListBrewCommandsCSV returns a list of brew commands in CSV format
		//
		// Returns the same brew commands data as ListBrewCommands but formatted as CSV
		ListBrewCommandsCSV(ctx context.Context) ([]byte, *interfaces.Response, error)

		// CreateBrewCommand creates a new brew command with specified arguments and optional device/timing configuration
		//
		// Requires arguments field. Optional fields include device_ids, run_after_datetime, and recurrence (once, daily, weekly, monthly)
		CreateBrewCommand(ctx context.Context, request *CreateBrewCommandRequest) (*CreateBrewCommandResponse, *interfaces.Response, error)

		// ListBrewCommandRuns returns a list of brew command runs for a specific brew command
		//
		// Returns run history including command, label, device, timestamps, success status, and output for the specified brew command label
		ListBrewCommandRuns(ctx context.Context, brewCommandLabel string) (*BrewCommandRunsResponse, *interfaces.Response, error)

		// ListBrewCommandRunsCSV returns a list of brew command runs in CSV format
		//
		// Returns the same run data as ListBrewCommandRuns but formatted as CSV
		ListBrewCommandRunsCSV(ctx context.Context, brewCommandLabel string) ([]byte, *interfaces.Response, error)
	}

	// Service handles communication with the brew commands
	// related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements BrewCommandsServiceInterface
var _ BrewCommandsServiceInterface = (*Service)(nil)

// NewService creates a new brew commands service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// ListBrewCommands retrieves all brew commands in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands.json
func (s *Service) ListBrewCommands(ctx context.Context) (*BrewCommandsResponse, *interfaces.Response, error) {
	endpoint := EndpointBrewCommandsJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewCommandsResponse
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListBrewCommandsCSV retrieves all brew commands in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands.csv
func (s *Service) ListBrewCommandsCSV(ctx context.Context) ([]byte, *interfaces.Response, error) {
	endpoint := EndpointBrewCommandsCSV

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

// CreateBrewCommand creates a new brew command
// URL: POST https://console.workbrew.com/workspaces/{workspace_name}/brew_commands.json
//
// Response codes:
//   - 201: Brew Command created successfully
//   - 403: On a Free tier plan (requires upgrade)
//   - 422: Validation error (e.g., "Arguments cannot include `&&`")
func (s *Service) CreateBrewCommand(ctx context.Context, request *CreateBrewCommandRequest) (*CreateBrewCommandResponse, *interfaces.Response, error) {
	endpoint := EndpointBrewCommandsJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result CreateBrewCommandResponse
	resp, err := s.client.Post(ctx, endpoint, request, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListBrewCommandRuns retrieves all runs for a specific brew command in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands/{brew_command_label}/runs.json
func (s *Service) ListBrewCommandRuns(ctx context.Context, brewCommandLabel string) (*BrewCommandRunsResponse, *interfaces.Response, error) {
	if brewCommandLabel == "" {
		return nil, nil, fmt.Errorf("brew command label is required")
	}

	endpoint := fmt.Sprintf(EndpointBrewCommandRunsJSONFormat, brewCommandLabel)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewCommandRunsResponse
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListBrewCommandRunsCSV retrieves all runs for a specific brew command in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands/{brew_command_label}/runs.csv
func (s *Service) ListBrewCommandRunsCSV(ctx context.Context, brewCommandLabel string) ([]byte, *interfaces.Response, error) {
	if brewCommandLabel == "" {
		return nil, nil, fmt.Errorf("brew command label is required")
	}

	endpoint := fmt.Sprintf(EndpointBrewCommandRunsCSVFormat, brewCommandLabel)

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
