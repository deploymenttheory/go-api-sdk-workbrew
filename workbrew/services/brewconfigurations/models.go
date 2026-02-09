package brewconfigurations

// BrewConfiguration represents a single brew configuration entry
type BrewConfiguration struct {
	Key               string `json:"key"`
	Value             string `json:"value"`
	LastUpdatedByUser string `json:"last_updated_by_user"`
	DeviceGroup       string `json:"device_group"`
}

// BrewConfigurationsResponse is the response from GET /brew_configurations.json
type BrewConfigurationsResponse []BrewConfiguration
