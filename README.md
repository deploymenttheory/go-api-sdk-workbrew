# Workbrew Go SDK

Community Go SDK for the [Workbrew API](https://console.workbrew.com/api-docs) v0.

## Features

- **Complete API Coverage**: All 36 endpoints across 13 services
- **Type-Safe**: Compile-time interface checking for all services
- **Context Support**: All methods accept `context.Context` for cancellation and timeouts
- **Comprehensive Error Handling**: Detailed error types for all HTTP status codes
- **CSV Export Support**: Download data in CSV format where available
- **Query Builder**: Fluent interface for building complex query parameters
- **Production Ready**: Follows proven architecture patterns from go-api-sdk-apple

## Installation

```bash
go get github.com/deploymenttheory/go-api-sdk-workbrew
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
)

func main() {
    // Create a new client
    client, err := workbrew.NewClient(
        "your-api-key",
        "your-workspace-name",
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Get all devices
    devices, err := client.Devices.GetDevices(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for _, device := range *devices {
        fmt.Printf("Device: %s\n", device.SerialNumber)
    }
}
```

## Environment Variables

You can also create a client from environment variables:

```bash
export WORKBREW_API_KEY="your-api-key"
export WORKBREW_WORKSPACE="your-workspace"
export WORKBREW_BASE_URL="https://console.workbrew.com"  # optional
export WORKBREW_API_VERSION="v0"  # optional
```

```go
client, err := workbrew.NewClientFromEnv()
```

## Available Services

All 13 Workbrew API services are supported:

### 1. Analytics
```go
// Get analytics data
analytics, err := client.Analytics.GetAnalytics(ctx)
csvData, err := client.Analytics.GetAnalyticsCSV(ctx)
```

### 2. Brew Commands
```go
// List brew commands
commands, err := client.BrewCommands.GetBrewCommands(ctx)

// Create a new brew command
request := &brewcommands.CreateBrewCommandRequest{
    Arguments: "install wget",
}
response, err := client.BrewCommands.CreateBrewCommand(ctx, request)

// Get command runs
runs, err := client.BrewCommands.GetBrewCommandRuns(ctx, "command-label")
```

### 3. Brew Configurations
```go
configs, err := client.BrewConfigurations.GetBrewConfigurations(ctx)
```

### 4. Brew Taps
```go
taps, err := client.BrewTaps.GetBrewTaps(ctx)
```

### 5. Brewfiles (Full CRUD)
```go
// List brewfiles
brewfiles, err := client.Brewfiles.GetBrewfiles(ctx)

// Create a brewfile
request := &brewfiles.CreateBrewfileRequest{
    Label:   "my-brewfile",
    Content: "brew \"wget\"",
}
created, err := client.Brewfiles.CreateBrewfile(ctx, request)

// Update a brewfile
updateReq := &brewfiles.UpdateBrewfileRequest{
    Content: "brew \"wget\"\nbrew \"htop\"",
}
updated, err := client.Brewfiles.UpdateBrewfile(ctx, "my-brewfile", updateReq)

// Delete a brewfile
deleted, err := client.Brewfiles.DeleteBrewfile(ctx, "my-brewfile")

// Get brewfile runs
runs, err := client.Brewfiles.GetBrewfileRuns(ctx, "my-brewfile")
```

### 6. Casks
```go
casks, err := client.Casks.GetCasks(ctx)
```

### 7. Device Groups
```go
groups, err := client.DeviceGroups.GetDeviceGroups(ctx)
```

### 8. Devices
```go
devices, err := client.Devices.GetDevices(ctx)
csvData, err := client.Devices.GetDevicesCSV(ctx)
```

### 9. Events
```go
// Get all events
events, err := client.Events.GetEvents(ctx, "")

// Filter by actor type
userEvents, err := client.Events.GetEvents(ctx, "user")
systemEvents, err := client.Events.GetEvents(ctx, "system")

// Export to CSV with download flag
csvData, err := client.Events.GetEventsCSV(ctx, "user", true)
```

### 10. Formulae
```go
formulae, err := client.Formulae.GetFormulae(ctx)
```

### 11. Licenses
```go
licenses, err := client.Licenses.GetLicenses(ctx)
```

### 12. Vulnerabilities
```go
vulns, err := client.Vulnerabilities.GetVulnerabilities(ctx)

// Note: May return 403 on Free tier plans
```

### 13. Vulnerability Changes
```go
// Get all vulnerability changes
changes, err := client.VulnerabilityChanges.GetVulnerabilityChanges(ctx, "", "")

// Filter by status
detected, err := client.VulnerabilityChanges.GetVulnerabilityChanges(ctx, "detected", "")
fixed, err := client.VulnerabilityChanges.GetVulnerabilityChanges(ctx, "fixed", "")

// Search for specific vulnerabilities
results, err := client.VulnerabilityChanges.GetVulnerabilityChanges(ctx, "detected", "curl")

// Export to CSV
csvData, err := client.VulnerabilityChanges.GetVulnerabilityChangesCSV(ctx, "detected", "curl", true)
```

## Error Handling

The SDK provides detailed error types and helper functions:

```go
devices, err := client.Devices.GetDevices(ctx)
if err != nil {
    // Check for specific error types
    if workbrewclient.IsUnauthorized(err) {
        log.Fatal("Invalid API key")
    }

    if workbrewclient.IsForbidden(err) {
        log.Fatal("Access forbidden")
    }

    if workbrewclient.IsFreeTierError(err) {
        log.Fatal("Feature not available on free tier")
    }

    if workbrewclient.IsNotFound(err) {
        log.Fatal("Resource not found")
    }

    if workbrewclient.IsValidationError(err) {
        log.Fatal("Validation error")
    }

    log.Fatal(err)
}
```

## Client Options

Customize client behavior with functional options:

```go
client, err := workbrew.NewClient(
    "your-api-key",
    "your-workspace",
    client.WithBaseURL("https://custom-url.com"),
    client.WithTimeout(30),
    client.WithRetryCount(3),
    client.WithAPIVersion("v0"),
    client.WithDebug(),
)
```

## Query Builder

For endpoints with complex query parameters:

```go

queryParams := client.QueryBuilder().
    AddString("filter", "user").
    AddInt("limit", 100).
    AddBool("include_deleted", false).
    Build()

// Use with custom API calls if needed
```

## API Documentation

For complete API documentation, see:
- [Workbrew API Documentation](https://console.workbrew.com/api-docs)
- [SDK Implementation Status](./SDK_IMPLEMENTATION_COMPLETE.md)

## Architecture

This SDK follows the same architecture pattern as the [go-api-sdk-apple](https://github.com/deploymenttheory/go-api-sdk-apple) project:

- **Service Interface Pattern**: Each service defines an interface with compile-time checking
- **Consistent Error Handling**: All HTTP status codes properly handled with helper functions
- **Type-Safe Models**: Request/response types for all endpoints
- **Context Support**: All operations support cancellation and timeouts

## Requirements

- Go 1.21 or higher
- Valid Workbrew API key
- Workbrew workspace name

## Dependencies

- [resty.dev/v3](https://github.com/go-resty/resty) - HTTP client
- [go.uber.org/zap](https://github.com/uber-go/zap) - Structured logging

## Contributing

This is an internal SDK following the deployment theology architecture standards. For issues or feature requests, please contact the development team.

## License

Copyright © 2024 Deployment Theology. All rights reserved.

## Related Projects

- [go-api-sdk-apple](https://github.com/deploymenttheory/go-api-sdk-apple) - Apple Business Manager/School Manager Go SDK
