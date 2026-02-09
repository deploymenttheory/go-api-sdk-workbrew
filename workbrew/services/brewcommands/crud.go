package brewcommands

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// BrewCommandsServiceInterface defines the interface for brew commands operations
	BrewCommandsServiceInterface interface {
		GetBrewCommands(ctx context.Context) (*BrewCommandsResponse, error)
		GetBrewCommandsCSV(ctx context.Context) ([]byte, error)
		CreateBrewCommand(ctx context.Context, request *CreateBrewCommandRequest) (*CreateBrewCommandResponse, error)
		GetBrewCommandRuns(ctx context.Context, brewCommandLabel string) (*BrewCommandRunsResponse, error)
		GetBrewCommandRunsCSV(ctx context.Context, brewCommandLabel string) ([]byte, error)
	}

	// Service handles communication with the brew commands
	// related methods of the Workbrew API.
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

// GetBrewCommands retrieves all brew commands in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/brew_commands.json"
func (s *Service) GetBrewCommands(ctx context.Context) (*BrewCommandsResponse, error) {
	endpoint := EndpointBrewCommandsJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewCommandsResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBrewCommandsCSV retrieves all brew commands in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands.csv
func (s *Service) GetBrewCommandsCSV(ctx context.Context) ([]byte, error) {
	endpoint := EndpointBrewCommandsCSV

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

// CreateBrewCommand creates a new brew command
// URL: POST https://console.workbrew.com/workspaces/{workspace_name}/brew_commands.json
//
// Response codes:
//   - 201: Brew Command created successfully
//   - 403: On a Free tier plan (requires upgrade)
//   - 422: Validation error (e.g., "Arguments cannot include `&&`")
//
// Example cURL:
//
//	curl -X POST \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  -H "Content-Type: application/json" \
//	  -d '{"arguments":"install wget"}' \
//	  "https://console.workbrew.com/workspaces/{workspace}/brew_commands.json"
func (s *Service) CreateBrewCommand(ctx context.Context, request *CreateBrewCommandRequest) (*CreateBrewCommandResponse, error) {
	endpoint := EndpointBrewCommandsJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result CreateBrewCommandResponse
	err := s.client.Post(ctx, endpoint, request, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBrewCommandRuns retrieves all runs for a specific brew command in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands/{brew_command_label}/runs.json
func (s *Service) GetBrewCommandRuns(ctx context.Context, brewCommandLabel string) (*BrewCommandRunsResponse, error) {
	if brewCommandLabel == "" {
		return nil, fmt.Errorf("brew command label is required")
	}

	endpoint := fmt.Sprintf(EndpointBrewCommandRunsJSONFormat, brewCommandLabel)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result BrewCommandRunsResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBrewCommandRunsCSV retrieves all runs for a specific brew command in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands/{brew_command_label}/runs.csv
func (s *Service) GetBrewCommandRunsCSV(ctx context.Context, brewCommandLabel string) ([]byte, error) {
	if brewCommandLabel == "" {
		return nil, fmt.Errorf("brew command label is required")
	}

	endpoint := fmt.Sprintf(EndpointBrewCommandRunsCSVFormat, brewCommandLabel)

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
