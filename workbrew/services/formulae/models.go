package formulae

// Formula represents a single formula entry
type Formula struct {
	Name                    string    `json:"name"`
	Devices                 []string  `json:"devices"`
	Outdated                bool      `json:"outdated"`
	InstalledOnRequest      bool      `json:"installed_on_request"`
	InstalledAsDependency   bool      `json:"installed_as_dependency"`
	Vulnerabilities         []string  `json:"vulnerabilities"`
	Deprecated              *string   `json:"deprecated"`
	License                 *[]string `json:"license"`
	HomebrewCoreVersion     *string   `json:"homebrew_core_version"`
}

// FormulaeResponse is the response from GET /formulae.json
type FormulaeResponse []Formula
