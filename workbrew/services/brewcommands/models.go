package brewcommands

import (
	"time"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"
)

// BrewCommand represents a brew command in the system
// Matches the schema from swagger specification
type BrewCommand struct {
	Command            string                `json:"command"`
	Label              string                `json:"label"`
	LastUpdatedByUser  string                `json:"last_updated_by_user"`
	StartedAt          devices.TimeOrStatus  `json:"started_at"`  // date-time or "Not Started"
	FinishedAt         devices.TimeOrStatus  `json:"finished_at"` // date-time or "Not Finished"
	Devices            []string              `json:"devices"`
	RunCount           int                   `json:"run_count"`
}

// BrewCommandsResponse represents the response from the brew_commands.json endpoint
type BrewCommandsResponse []BrewCommand

// CreateBrewCommandRequest represents the request body for creating a brew command
// Per swagger spec
type CreateBrewCommandRequest struct {
	Arguments        string  `json:"arguments"`                   // Required: brew arguments (e.g., "install wget")
	DeviceIDs        *string `json:"device_ids,omitempty"`        // Optional: comma-separated UUIDs
	RunAfterDatetime *string `json:"run_after_datetime,omitempty"` // Optional: date_time format (e.g., "2025-01-10T10:09")
	Recurrence       *string `json:"recurrence,omitempty"`        // Optional: "once", "daily", "weekly", "monthly"
}

// CreateBrewCommandResponse represents the successful response from creating a brew command
// Status code: 201
type CreateBrewCommandResponse struct {
	Message string `json:"message"`
}

// BrewCommandRun represents a single execution of a brew command
// Matches the schema from swagger specification
type BrewCommandRun struct {
	Command    string               `json:"command"`
	Label      string               `json:"label"`
	Device     string               `json:"device"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
	Success    bool                 `json:"success"`
	Output     string               `json:"output"`
	StartedAt  devices.TimeOrStatus `json:"started_at"`  // date-time or "Not Started"
	FinishedAt devices.TimeOrStatus `json:"finished_at"` // date-time or "Not Finished"
}

// BrewCommandRunsResponse represents the response from the runs.json endpoint
type BrewCommandRunsResponse []BrewCommandRun
