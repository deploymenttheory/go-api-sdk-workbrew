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
		ListDeviceGroups(ctx context.Context) (*DeviceGroupsResponse, error)

		// ListDeviceGroupsCSV returns a list of Device Groups in CSV format
		//
		// Returns device group data as CSV with columns: id, name, devices.
		ListDeviceGroupsCSV(ctx context.Context) ([]byte, error)
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
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/device_groups.json"
func (s *Service) ListDeviceGroups(ctx context.Context) (*DeviceGroupsResponse, error) {
	endpoint := EndpointDeviceGroupsJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result DeviceGroupsResponse
	err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListDeviceGroupsCSV retrieves all device groups in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/device_groups.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/device_groups.csv"
func (s *Service) ListDeviceGroupsCSV(ctx context.Context) ([]byte, error) {
	endpoint := EndpointDeviceGroupsCSV

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
