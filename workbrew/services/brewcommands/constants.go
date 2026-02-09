package brewcommands

// API endpoints for brew commands
const (
	EndpointBrewCommandsJSON          = "/brew_commands.json"
	EndpointBrewCommandsCSV           = "/brew_commands.csv"
	EndpointBrewCommandRunsJSONFormat = "/brew_commands/%s/runs.json" // {brew_command_label}
	EndpointBrewCommandRunsCSVFormat  = "/brew_commands/%s/runs.csv"  // {brew_command_label}
)

// Recurrence values per swagger spec
const (
	RecurrenceOnce    = "once"
	RecurrenceDaily   = "daily"
	RecurrenceWeekly  = "weekly"
	RecurrenceMonthly = "monthly"
)
