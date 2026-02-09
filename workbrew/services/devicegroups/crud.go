package devicegroups

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// DeviceGroupsServiceInterface defines the interface for device groups operations
	DeviceGroupsServiceInterface interface {
		GetDeviceGroups(ctx context.Context) (*DeviceGroupsResponse, error)
		GetDeviceGroupsCSV(ctx context.Context) ([]byte, error)
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

// GetDeviceGroups retrieves all device groups in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/device_groups.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/device_groups.json"
func (s *Service) GetDeviceGroups(ctx context.Context) (*DeviceGroupsResponse, error) {
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

// GetDeviceGroupsCSV retrieves all device groups in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/device_groups.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/device_groups.csv"
func (s *Service) GetDeviceGroupsCSV(ctx context.Context) ([]byte, error) {
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
