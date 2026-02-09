package devicegroups

// DeviceGroup represents a single device group entry
type DeviceGroup struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Devices []string `json:"devices"`
}

// DeviceGroupsResponse is the response from GET /device_groups.json
type DeviceGroupsResponse []DeviceGroup
