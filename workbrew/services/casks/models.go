package casks

// Cask represents a single cask entry
type Cask struct {
	Name                 string   `json:"name"`
	DisplayName          *string  `json:"display_name"`
	Devices              []string `json:"devices"`
	Outdated             bool     `json:"outdated"`
	Deprecated           *string  `json:"deprecated"`
	HomebrewCaskVersion  *string  `json:"homebrew_cask_version"`
}

// CasksResponse is the response from GET /casks.json
type CasksResponse []Cask
