package brewtaps

// BrewTap represents a single brew tap entry
type BrewTap struct {
	Tap                string   `json:"tap"`
	Devices            []string `json:"devices"`
	FormulaeInstalled  int      `json:"formulae_installed"`
	CasksInstalled     int      `json:"casks_installed"`
	AvailablePackages  string   `json:"available_packages"`
}

// BrewTapsResponse is the response from GET /brew_taps.json
type BrewTapsResponse []BrewTap
