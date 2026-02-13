package devicegroups

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// DeviceGroupsServiceInterface defines the interface for device groups operations
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	DeviceGroupsServiceInterface interface {
		// ListDeviceGroups returns a list of Device Groups
		//
		// Returns device groups with their IDs, names, and assigned device serial numbers.
		ListDeviceGroups(ctx context.Context) (*DeviceGroupsResponse, *interfaces.Response, error)

		// ListDeviceGroupsCSV returns a list of Device Groups in CSV format
		//
		// Returns device group data as CSV with columns: id, name, devices.
		ListDeviceGroupsCSV(ctx context.Context) ([]byte, *interfaces.Response, error)
	}

	// Service handles communication with the device groups
	// related methods of the Workbrew API.
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements DeviceGroupsServiceInterface
var _ DeviceGroupsServiceInterface = (*Service)(nil)

// NewService creates a new device groups service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// ListDeviceGroups retrieves all device groups in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/device_groups.json
func (s *Service) ListDeviceGroups(ctx context.Context) (*DeviceGroupsResponse, *interfaces.Response, error) {
	endpoint := EndpointDeviceGroupsJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result DeviceGroupsResponse
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListDeviceGroupsCSV retrieves all device groups in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/device_groups.csv
func (s *Service) ListDeviceGroupsCSV(ctx context.Context) ([]byte, *interfaces.Response, error) {
	endpoint := EndpointDeviceGroupsCSV

	headers := map[string]string{
		"Accept": "text/csv",
	}

	queryParams := make(map[string]string)

	resp, csvData, err := s.client.GetBytes(ctx, endpoint, queryParams, headers)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
