package brewfiles

// BrewfileDevice represents a device associated with a brewfile
type BrewfileDevice struct {
	SerialNumber string `json:"serial_number,omitempty"`
}

// Brewfile represents a single brewfile entry
type Brewfile struct {
	Label              string           `json:"label,omitempty"`
	Slug               string           `json:"slug,omitempty"`
	Content            string           `json:"content,omitempty"`
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
// Matches the actual API response schema
type BrewfileRun struct {
	Label      string `json:"label"`
	Device     string `json:"device"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Success    bool   `json:"success"`
	Output     string `json:"output"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
}

// BrewfileRunsResponse is the response from GET /brewfiles/{label}/runs.json
type BrewfileRunsResponse []BrewfileRun
