package vulnerabilities

// VulnerabilityDetail represents a single vulnerability
type VulnerabilityDetail struct {
	CleanID   string   `json:"clean_id"`
	CVSSScore *float64 `json:"cvss_score"`
}

// Vulnerability represents a vulnerability entry with associated formula and devices
type Vulnerability struct {
	Vulnerabilities       []VulnerabilityDetail `json:"vulnerabilities"`
	Formula               string                `json:"formula"`
	OutdatedDevices       []string              `json:"outdated_devices"`
	Supported             bool                  `json:"supported"`
	HomebrewCoreVersion   string                `json:"homebrew_core_version"`
}

// VulnerabilitiesResponse is the response from GET /vulnerabilities.json
type VulnerabilitiesResponse []Vulnerability
