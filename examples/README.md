# Workbrew API SDK Examples

This directory contains comprehensive examples for all Workbrew API SDK functions. Each example is a complete, runnable Go program demonstrating how to use a specific API operation.

## Prerequisites

Set the following environment variables before running any example:

```bash
export WORKBREW_API_KEY="your-api-key"
export WORKBREW_WORKSPACE="your-workspace"
```

Some examples require additional environment variables (noted in each example).

## Running Examples

Navigate to any example directory and run:

```bash
cd examples/workbrew/<service>/<function>
go run main.go
```

## Available Examples

### Analytics (2 examples)
- `analytics/ListAnalytics` - List analytics data
- `analytics/ListAnalyticsCSV` - List analytics data in CSV format

### Brew Commands (5 examples)
- `brewcommands/ListBrewCommands` - List all brew commands
- `brewcommands/ListBrewCommandsCSV` - List brew commands in CSV format
- `brewcommands/CreateBrewCommand` - Create a new brew command
- `brewcommands/ListBrewCommandRuns` - List runs for a specific brew command (requires `BREW_COMMAND_LABEL`)
- `brewcommands/ListBrewCommandRunsCSV` - List brew command runs in CSV format (requires `BREW_COMMAND_LABEL`)

### Brew Configurations (2 examples)
- `brewconfigurations/ListBrewConfigurations` - List brew configurations
- `brewconfigurations/ListBrewConfigurationsCSV` - List brew configurations in CSV format

### Brewfiles (7 examples)
- `brewfiles/ListBrewfiles` - List all brewfiles
- `brewfiles/ListBrewfilesCSV` - List brewfiles in CSV format
- `brewfiles/CreateBrewfile` - Create a new brewfile
- `brewfiles/UpdateBrewfile` - Update an existing brewfile (requires `BREWFILE_LABEL`)
- `brewfiles/DeleteBrewfile` - Delete a brewfile (requires `BREWFILE_LABEL`)
- `brewfiles/ListBrewfileRuns` - List runs for a specific brewfile (requires `BREWFILE_LABEL`)
- `brewfiles/ListBrewfileRunsCSV` - List brewfile runs in CSV format (requires `BREWFILE_LABEL`)

### Brew Taps (2 examples)
- `brewtaps/ListBrewTaps` - List all brew taps
- `brewtaps/ListBrewTapsCSV` - List brew taps in CSV format

### Casks (2 examples)
- `casks/ListCasks` - List all casks
- `casks/ListCasksCSV` - List casks in CSV format

### Device Groups (2 examples)
- `devicegroups/ListDeviceGroups` - List all device groups
- `devicegroups/ListDeviceGroupsCSV` - List device groups in CSV format

### Devices (2 examples)
- `devices/ListDevices` - List all devices
- `devices/ListDevicesCSV` - List devices in CSV format

### Events (2 examples)
- `events/ListEvents` - List audit log events
- `events/ListEventsCSV` - List events in CSV format

### Formulae (2 examples)
- `formulae/ListFormulae` - List all formulae
- `formulae/ListFormulaeCSV` - List formulae in CSV format

### Licenses (2 examples)
- `licenses/ListLicenses` - List all licenses
- `licenses/ListLicensesCSV` - List licenses in CSV format

### Vulnerabilities (2 examples)
- `vulnerabilities/ListVulnerabilities` - List all vulnerabilities
- `vulnerabilities/ListVulnerabilitiesCSV` - List vulnerabilities in CSV format

### Vulnerability Changes (2 examples)
- `vulnerabilitychanges/ListVulnerabilityChanges` - List vulnerability change events
- `vulnerabilitychanges/ListVulnerabilityChangesCSV` - List vulnerability changes in CSV format

## Total: 34 Examples

Each example includes:
- Proper error handling
- Logging setup with zap
- Environment variable validation
- Clear output formatting
- Comments explaining key concepts
