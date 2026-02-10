# Go API SDK for Workbrew

[![Go Reference](https://pkg.go.dev/badge/github.com/deploymenttheory/go-api-sdk-workbrew.svg)](https://pkg.go.dev/github.com/deploymenttheory/go-api-sdk-workbrew)
[![Go Report Card](https://goreportcard.com/badge/github.com/deploymenttheory/go-api-sdk-workbrew)](https://goreportcard.com/report/github.com/deploymenttheory/go-api-sdk-workbrew)
[![License](https://img.shields.io/github/license/deploymenttheory/go-api-sdk-workbrew)](https://github.com/deploymenttheory/go-api-sdk-workbrew/blob/main/LICENSE)

Community Go SDK for the [Workbrew API](https://console.workbrew.com/documentation/api) v0.

## Features

- **Complete API Coverage**: All 36 endpoints across 13 services
- **Type-Safe**: Compile-time interface checking for all services
- **Context Support**: All methods accept `context.Context` for cancellation and timeouts
- **Comprehensive Error Handling**: Detailed error types for all HTTP status codes
- **CSV Export Support**: Download data in CSV format where available
- **Production Ready**: Follows proven architecture patterns from go-api-sdk-apple
- **34 Complete Examples**: Working examples for every CRUD operation

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
    "os"

    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"
    "go.uber.org/zap"
)

func main() {
    // Method 1: Direct client creation
    apiKey := os.Getenv("WORKBREW_API_KEY")
    workspace := os.Getenv("WORKBREW_WORKSPACE")

    logger, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }
    defer logger.Sync()

    httpClient, err := client.NewClient(apiKey, workspace,
        client.WithLogger(logger),
        client.WithBaseURL("https://console.workbrew.com"),
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create services
    devicesService := devices.NewService(httpClient)

    // List all devices
    ctx := context.Background()
    devicesList, err := devicesService.ListDevices(ctx)
    if err != nil {
        log.Fatalf("Failed to list devices: %v", err)
    }

    fmt.Printf("Found %d devices\n", len(*devicesList))
    for _, device := range *devicesList {
        fmt.Printf("Device: %s (%s)\n", device.SerialNumber, device.DeviceType)
    }
}
```

## Environment Variables

Set these environment variables for easier client configuration:

```bash
export WORKBREW_API_KEY="your-api-key"
export WORKBREW_WORKSPACE="your-workspace"
export WORKBREW_BASE_URL="https://console.workbrew.com"  # optional
```

## Supported Services

All 13 Workbrew API services are fully supported with complete CRUD operations where applicable:

| Service | Endpoints | CRUD Support |
|---------|-----------|--------------|
| **Analytics** | 2 | List, List CSV |
| **Brew Commands** | 5 | List, Create, List Runs, CSV exports |
| **Brew Configurations** | 2 | List, List CSV |
| **Brew Taps** | 2 | List, List CSV |
| **Brewfiles** | 7 | Full CRUD + Runs |
| **Casks** | 2 | List, List CSV |
| **Device Groups** | 2 | List, List CSV |
| **Devices** | 2 | List, List CSV |
| **Events** | 2 | List (filterable), List CSV |
| **Formulae** | 2 | List, List CSV |
| **Licenses** | 2 | List, List CSV |
| **Vulnerabilities** | 2 | List, List CSV |
| **Vulnerability Changes** | 2 | List (filterable), List CSV |

📖 **[Browse Examples →](./examples/workbrew/)**

The examples directory contains a working example for every CRUD operation, organized by service:
- `analytics/` - 2 examples
- `brewcommands/` - 5 examples
- `brewconfigurations/` - 2 examples
- `brewtaps/` - 2 examples
- `brewfiles/` - 7 examples (full CRUD)
- `casks/` - 2 examples
- `devicegroups/` - 2 examples
- `devices/` - 2 examples
- `events/` - 2 examples
- `formulae/` - 2 examples
- `licenses/` - 2 examples
- `vulnerabilities/` - 2 examples
- `vulnerabilitychanges/` - 2 examples

Each example is a standalone executable demonstrating best practices.

## Common Usage Patterns

### Listing Resources

```go
// Most services follow this pattern
devicesService := devices.NewService(httpClient)
devicesList, err := devicesService.ListDevices(ctx)

formulaeService := formulae.NewService(httpClient)
formulaeList, err := formulaeService.ListFormulae(ctx)
```

### Creating Resources

```go
// Brewfiles support full CRUD
brewfilesService := brewfiles.NewService(httpClient)

deviceSerial := "TC6R2DHVHG"
request := &brewfiles.CreateBrewfileRequest{
    Label:               "my-brewfile",
    Content:             "brew \"wget\"\nbrew \"htop\"",
    DeviceSerialNumbers: &deviceSerial,
}

response, err := brewfilesService.CreateBrewfile(ctx, request)
```

### Filtering and Querying

```go
// Events support filtering by actor type
eventsService := events.NewService(httpClient)

opts := &events.RequestQueryOptions{
    Filter: "user", // "user", "system", or "" for all
}

eventsList, err := eventsService.ListEvents(ctx, opts)
```

### CSV Export

```go
// Most list endpoints have CSV variants
devicesService := devices.NewService(httpClient)
csvData, err := devicesService.ListDevicesCSV(ctx)

// Write to file or process as needed
os.WriteFile("devices.csv", csvData, 0644)
```

## Error Handling

The SDK provides centralized error handling with detailed error information:

```go
import "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"

devices, err := devicesService.ListDevices(ctx)
if err != nil {
    // All API errors are *client.APIError
    if apiErr, ok := err.(*client.APIError); ok {
        fmt.Printf("API Error: %s (Status: %d)\n", apiErr.Message, apiErr.StatusCode)
        fmt.Printf("Endpoint: %s %s\n", apiErr.Method, apiErr.Endpoint)
        
        // Check specific status codes
        switch apiErr.StatusCode {
        case 401:
            log.Fatal("Invalid API key")
        case 403:
            log.Fatal("Access forbidden (may require paid tier)")
        case 404:
            log.Fatal("Resource not found")
        case 422:
            log.Printf("Validation errors: %v", apiErr.Errors)
        }
    }
    
    log.Fatalf("Error: %v", err)
}
```

## Client Options

Customize client behavior with functional options:

```go
import "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"

httpClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithBaseURL("https://custom-url.com"),
    client.WithTimeout(30 * time.Second),
    client.WithRetryCount(3),
    client.WithLogger(customLogger),
    client.WithDebug(true),
)
```

## Special Types

### TimeOrStatus

Brew commands use a special `devices.TimeOrStatus` type that handles multiple status values:

```go
import "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"

// Check status
if command.StartedAt.IsNotStarted() {
    fmt.Println("Command has not started")
}

if command.FinishedAt.HasTime() {
    fmt.Printf("Finished at: %s\n", command.FinishedAt.String())
}
```

### Pointer Types

Optional fields use pointer types for proper null handling:

```go
// Creating requests with optional fields
deviceSerial := "ABC123"
request := &brewfiles.CreateBrewfileRequest{
    Label:               "my-brewfile",
    Content:             "brew \"wget\"",
    DeviceSerialNumbers: &deviceSerial, // pointer to string
}

// Reading responses with optional fields
if device.MDMUserOrDeviceName != nil {
    fmt.Printf("MDM Name: %s\n", *device.MDMUserOrDeviceName)
}
```

## API Documentation

For complete API documentation, see:
- [Workbrew API Documentation](https://console.workbrew.com/documentation/api)
- [SDK Examples Directory](./examples/workbrew/)

## Architecture

This SDK follows the same architecture pattern as the [go-api-sdk-apple](https://github.com/deploymenttheory/go-api-sdk-apple) project:

- **Service Interface Pattern**: Each service defines an interface with compile-time checking
- **Consistent Error Handling**: All HTTP status codes properly handled with centralized error types
- **Type-Safe Models**: Request/response types for all endpoints
- **Context Support**: All operations support cancellation and timeouts
- **Structured Logging**: Comprehensive request/response logging with zap

## Requirements

- Go 1.21 or higher
- Valid Workbrew API key
- Workbrew workspace name

## Dependencies

- [go.uber.org/zap](https://github.com/uber-go/zap) - Structured logging
- Standard library only for HTTP client

## Contributing

This is a community SDK following deployment Theory architecture standards. For issues or feature requests, please open an issue on GitHub.

## License

Copyright © 2026 Deployment Theory. All rights reserved.