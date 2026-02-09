package devices

import (
	"time"
)

// Device represents a device in the workspace
// Matches the schema from swagger specification
type Device struct {
	SerialNumber        string      `json:"serial_number"`
	Groups              []string    `json:"groups"`
	MDMUserOrDeviceName *string     `json:"mdm_user_or_device_name"` // nullable
	LastSeenAt          TimeOrNever `json:"last_seen_at"`            // date-time or "Never"
	CommandLastRunAt    TimeOrNever `json:"command_last_run_at"`     // date-time or "Never"
	DeviceType          string      `json:"device_type"`
	OSVersion           string      `json:"os_version"`
	HomebrewPrefix      string      `json:"homebrew_prefix"`
	HomebrewVersion     string      `json:"homebrew_version"`
	WorkbrewVersion     string      `json:"workbrew_version"`
	FormulaeCount       int         `json:"formulae_count"`
	CasksCount          int         `json:"casks_count"`
}

// DevicesResponse represents the response from the devices.json endpoint
type DevicesResponse []Device

// TimeOrNever handles the "oneOf" type from swagger: date-time or "Never"
type TimeOrNever struct {
	Time  *time.Time
	Never bool
}

// UnmarshalJSON implements custom unmarshaling for TimeOrNever
func (t *TimeOrNever) UnmarshalJSON(data []byte) error {
	str := string(data)

	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	if str == "Never" {
		t.Never = true
		t.Time = nil
		return nil
	}

	parsedTime, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	t.Time = &parsedTime
	t.Never = false
	return nil
}

// MarshalJSON implements custom marshaling for TimeOrNever
func (t TimeOrNever) MarshalJSON() ([]byte, error) {
	if t.Never || t.Time == nil {
		return []byte(`"Never"`), nil
	}
	return []byte(`"` + t.Time.Format(time.RFC3339) + `"`), nil
}

// String returns a string representation of TimeOrNever
func (t TimeOrNever) String() string {
	if t.Never || t.Time == nil {
		return "Never"
	}
	return t.Time.Format(time.RFC3339)
}

// TimeOrStatus handles the "oneOf" type from swagger: date-time or status strings
// Supports: "Never", "Not Started", "Not Finished", and RFC3339 date-time
type TimeOrStatus struct {
	Time   *time.Time
	Status string // "Never", "Not Started", "Not Finished", or empty if Time is set
}

// UnmarshalJSON implements custom unmarshaling for TimeOrStatus
func (t *TimeOrStatus) UnmarshalJSON(data []byte) error {
	str := string(data)

	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	switch str {
	case "Never", "Not Started", "Not Finished":
		t.Status = str
		t.Time = nil
		return nil
	}

	parsedTime, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	t.Time = &parsedTime
	t.Status = ""
	return nil
}

// MarshalJSON implements custom marshaling for TimeOrStatus
func (t TimeOrStatus) MarshalJSON() ([]byte, error) {
	if t.Status != "" {
		return []byte(`"` + t.Status + `"`), nil
	}
	if t.Time == nil {
		return []byte(`"Never"`), nil
	}
	return []byte(`"` + t.Time.Format(time.RFC3339) + `"`), nil
}

// String returns a string representation of TimeOrStatus
func (t TimeOrStatus) String() string {
	if t.Status != "" {
		return t.Status
	}
	if t.Time == nil {
		return "Never"
	}
	return t.Time.Format(time.RFC3339)
}

// IsNever returns true if the status is "Never"
func (t TimeOrStatus) IsNever() bool {
	return t.Status == "Never"
}

// IsNotStarted returns true if the status is "Not Started"
func (t TimeOrStatus) IsNotStarted() bool {
	return t.Status == "Not Started"
}

// IsNotFinished returns true if the status is "Not Finished"
func (t TimeOrStatus) IsNotFinished() bool {
	return t.Status == "Not Finished"
}

// HasTime returns true if a valid time is set
func (t TimeOrStatus) HasTime() bool {
	return t.Time != nil
}
