package brewfiles

const (
	// EndpointBrewfilesJSON is the endpoint for brewfiles in JSON format
	EndpointBrewfilesJSON = "/brewfiles.json"

	// EndpointBrewfilesCSV is the endpoint for brewfiles in CSV format
	EndpointBrewfilesCSV = "/brewfiles.csv"

	// EndpointBrewfileLabelFormat is the format string for a specific brewfile by label
	EndpointBrewfileLabelFormat = "/brewfiles/%s.json"

	// EndpointBrewfileRunsJSONFormat is the format string for brewfile runs in JSON
	EndpointBrewfileRunsJSONFormat = "/brewfiles/%s/runs.json"

	// EndpointBrewfileRunsCSVFormat is the format string for brewfile runs in CSV
	EndpointBrewfileRunsCSVFormat = "/brewfiles/%s/runs.csv"
)
