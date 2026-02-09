package licenses

// License represents a single license entry
type License struct {
	Name         string `json:"name"`
	DeviceCount  int    `json:"device_count"`
	FormulaCount int    `json:"formula_count"`
}

// LicensesResponse is the response from GET /licenses.json
type LicensesResponse []License
