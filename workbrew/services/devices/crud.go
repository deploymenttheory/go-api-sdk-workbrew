package devices

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
)

type (
	// DevicesServiceInterface defines the interface for devices operations
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	DevicesServiceInterface interface {
		// ListDevices returns a list of devices
		//
		// Returns devices with serial numbers, group assignments, MDM names, last seen timestamps, device types,
		// OS versions, Homebrew/Workbrew versions, and installed package counts.
		ListDevices(ctx context.Context) (*DevicesResponse, *interfaces.Response, error)

		// ListDevicesCSV returns a list of devices in CSV format
		//
		// Returns device data as CSV with columns: serial_number, groups, mdm_user_or_device_name, last_seen_at,
		// command_last_run_at, device_type, os_version, homebrew_prefix, homebrew_version, workbrew_version, formulae_count, casks_count.
		ListDevicesCSV(ctx context.Context) ([]byte, *interfaces.Response, error)
	}

	// Service handles communication with the devices
	// related methods of the Workbrew API.
	Service struct {
		client interfaces.HTTPClient
	}
)

// Ensure Service implements DevicesServiceInterface
var _ DevicesServiceInterface = (*Service)(nil)

// NewService creates a new devices service
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// ListDevices retrieves all devices in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/devices.json
func (s *Service) ListDevices(ctx context.Context) (*DevicesResponse, *interfaces.Response, error) {
	endpoint := EndpointDevicesJSON

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	queryParams := make(map[string]string)

	var result DevicesResponse
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListDevicesCSV retrieves all devices in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/devices.csv
func (s *Service) ListDevicesCSV(ctx context.Context) ([]byte, *interfaces.Response, error) {
	endpoint := EndpointDevicesCSV

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
