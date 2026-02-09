package brewfiles

// BrewfileDevice represents a device associated with a brewfile
type BrewfileDevice struct {
	SerialNumber string `json:"serial_number,omitempty"`
}

// Brewfile represents a single brewfile entry
type Brewfile struct {
	Label              string           `json:"label"`
	Slug               string           `json:"slug"`
	Content            string           `json:"content"`
	LastUpdatedByUser  string           `json:"last_updated_by_user"`
	StartedAt          string           `json:"started_at"`
	FinishedAt         string           `json:"finished_at"`
	Devices            []BrewfileDevice `json:"devices"`
	RunCount           int              `json:"run_count"`
}

// BrewfilesResponse is the response from GET /brewfiles.json
type BrewfilesResponse []Brewfile

// CreateBrewfileRequest represents the request body for creating a brewfile
type CreateBrewfileRequest struct {
	Label               string  `json:"label"`
	Content             string  `json:"content"`
	DeviceSerialNumbers *string `json:"device_serial_numbers,omitempty"`
	DeviceGroupID       *string `json:"device_group_id,omitempty"`
}

// UpdateBrewfileRequest represents the request body for updating a brewfile
type UpdateBrewfileRequest struct {
	Content             string  `json:"content"`
	DeviceSerialNumbers *string `json:"device_serial_numbers,omitempty"`
	DeviceGroupID       *string `json:"device_group_id,omitempty"`
}

// BrewfileMessageResponse represents a simple message response
type BrewfileMessageResponse struct {
	Message string `json:"message"`
}

// BrewfileRun represents a single brewfile run
type BrewfileRun struct {
	Device     string `json:"device,omitempty"`
	StartedAt  string `json:"started_at,omitempty"`
	FinishedAt string `json:"finished_at,omitempty"`
	Status     string `json:"status,omitempty"`
}

// BrewfileRunsResponse is the response from GET /brewfiles/{label}/runs.json
type BrewfileRunsResponse []BrewfileRun
