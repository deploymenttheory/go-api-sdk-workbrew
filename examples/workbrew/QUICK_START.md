# Workbrew SDK Quick Start Guide

This guide covers different ways to set up and use the Workbrew Go SDK, from simple to advanced configurations.

## Table of Contents

1. [Environment Variables Setup](#1-environment-variables-setup)
2. [Basic Client Setup](#2-basic-client-setup)
3. [Client with Custom Options](#3-client-with-custom-options)
4. [Working with Services](#4-working-with-services)
5. [Error Handling](#5-error-handling)
6. [CSV Export](#6-csv-export)
7. [Creating and Updating Resources](#7-creating-and-updating-resources)

---

## 1. Environment Variables Setup

The simplest way to get started is using environment variables:

```bash
export WORKBREW_API_KEY="your-api-key-here"
export WORKBREW_WORKSPACE="your-workspace-name"
```

## 2. Basic Client Setup

### Minimal Example

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
    apiKey := os.Getenv("WORKBREW_API_KEY")
    workspace := os.Getenv("WORKBREW_WORKSPACE")

    if apiKey == "" || workspace == "" {
        log.Fatal("WORKBREW_API_KEY and WORKBREW_WORKSPACE must be set")
    }

    // Create a simple logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    // Create HTTP client
    httpClient, err := client.NewClient(apiKey, workspace,
        client.WithLogger(logger),
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create a service
    devicesService := devices.NewService(httpClient)

    // List devices
    ctx := context.Background()
    devicesList, err := devicesService.ListDevices(ctx)
    if err != nil {
        log.Fatalf("Failed to list devices: %v", err)
    }

    fmt.Printf("Found %d devices\n", len(*devicesList))
}
```

## 3. Client with Custom Options

### Production Configuration

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
    "go.uber.org/zap"
)

func main() {
    // Create production logger with sampling
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    // Create client with custom options
    httpClient, err := client.NewClient(
        "your-api-key",
        "your-workspace",
        client.WithBaseURL("https://console.workbrew.com"),
        client.WithTimeout(30*time.Second),
        client.WithRetryCount(3),
        client.WithLogger(logger),
        client.WithDebug(false), // Set to true for development
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Client is ready to use
    _ = httpClient
}
```

### Development Configuration

```go
func main() {
    // Create development logger with debug output
    logger, _ := zap.NewDevelopment()
    defer logger.Sync()

    httpClient, err := client.NewClient(
        "your-api-key",
        "your-workspace",
        client.WithLogger(logger),
        client.WithDebug(true), // Enable debug logging
        client.WithTimeout(60*time.Second), // Longer timeout for debugging
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Client is ready to use
    _ = httpClient
}
```

## 4. Working with Services

### List Operations

```go
import (
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/formulae"
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/casks"
)

func listResources(httpClient interfaces.HTTPClient) {
    ctx := context.Background()

    // List devices
    devicesService := devices.NewService(httpClient)
    devicesList, err := devicesService.ListDevices(ctx)
    if err != nil {
        log.Printf("Error: %v", err)
    }
    fmt.Printf("Devices: %d\n", len(*devicesList))

    // List formulae
    formulaeService := formulae.NewService(httpClient)
    formulaeList, err := formulaeService.ListFormulae(ctx)
    if err != nil {
        log.Printf("Error: %v", err)
    }
    fmt.Printf("Formulae: %d\n", len(*formulaeList))

    // List casks
    casksService := casks.NewService(httpClient)
    casksList, err := casksService.ListCasks(ctx)
    if err != nil {
        log.Printf("Error: %v", err)
    }
    fmt.Printf("Casks: %d\n", len(*casksList))
}
```

### Filtered Queries

```go
import (
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/events"
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/vulnerabilitychanges"
)

func filteredQueries(httpClient interfaces.HTTPClient) {
    ctx := context.Background()

    // List events filtered by actor type
    eventsService := events.NewService(httpClient)
    
    opts := &events.RequestQueryOptions{
        Filter: "user", // "user", "system", or "" for all
    }
    
    eventsList, err := eventsService.ListEvents(ctx, opts)
    if err != nil {
        log.Printf("Error: %v", err)
    }
    fmt.Printf("User events: %d\n", len(*eventsList))

    // List vulnerability changes with filters
    vulnService := vulnerabilitychanges.NewService(httpClient)
    
    vulnOpts := &vulnerabilitychanges.RequestQueryOptions{
        Status: "detected", // "detected", "fixed", or ""
        Query:  "curl",     // Search term
    }
    
    changes, err := vulnService.ListVulnerabilityChanges(ctx, vulnOpts)
    if err != nil {
        log.Printf("Error: %v", err)
    }
    fmt.Printf("Vulnerability changes: %d\n", len(*changes))
}
```

## 5. Error Handling

### Comprehensive Error Handling

```go
import "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"

func handleErrors(devicesService devices.DevicesServiceInterface) {
    ctx := context.Background()
    
    devicesList, err := devicesService.ListDevices(ctx)
    if err != nil {
        // Check if it's an API error
        if apiErr, ok := err.(*client.APIError); ok {
            fmt.Printf("API Error: %s\n", apiErr.Message)
            fmt.Printf("Status Code: %d\n", apiErr.StatusCode)
            fmt.Printf("Endpoint: %s %s\n", apiErr.Method, apiErr.Endpoint)
            
            // Handle specific status codes
            switch apiErr.StatusCode {
            case 401:
                log.Fatal("Authentication failed - check your API key")
            case 403:
                log.Fatal("Access forbidden - this feature may require a paid tier")
            case 404:
                log.Fatal("Resource not found")
            case 422:
                log.Printf("Validation errors: %v", apiErr.Errors)
            case 500:
                log.Fatal("Server error - try again later")
            default:
                log.Fatalf("API error: %s", apiErr.Message)
            }
        } else {
            // Network or other error
            log.Fatalf("Error: %v", err)
        }
        return
    }

    // Success
    fmt.Printf("Retrieved %d devices\n", len(*devicesList))
}
```

## 6. CSV Export

### Exporting to CSV

```go
import (
    "os"
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"
)

func exportToCSV(httpClient interfaces.HTTPClient) error {
    ctx := context.Background()
    
    devicesService := devices.NewService(httpClient)
    
    // Get CSV data
    csvData, err := devicesService.ListDevicesCSV(ctx)
    if err != nil {
        return fmt.Errorf("failed to get CSV: %w", err)
    }

    // Write to file
    err = os.WriteFile("devices.csv", csvData, 0644)
    if err != nil {
        return fmt.Errorf("failed to write file: %w", err)
    }

    fmt.Printf("Exported %d bytes to devices.csv\n", len(csvData))
    return nil
}
```

## 7. Creating and Updating Resources

### Creating a Brewfile

```go
import "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewfiles"

func createBrewfile(httpClient interfaces.HTTPClient) error {
    ctx := context.Background()
    
    brewfilesService := brewfiles.NewService(httpClient)

    // Create a brewfile
    deviceSerial := "TC6R2DHVHG"
    request := &brewfiles.CreateBrewfileRequest{
        Label:               "my-brewfile",
        Content:             "brew \"wget\"\nbrew \"htop\"",
        DeviceSerialNumbers: &deviceSerial, // Pointer to string
    }

    response, err := brewfilesService.CreateBrewfile(ctx, request)
    if err != nil {
        return fmt.Errorf("failed to create brewfile: %w", err)
    }

    fmt.Printf("Success: %s\n", response.Message)
    return nil
}
```

### Updating a Brewfile

```go
func updateBrewfile(httpClient interfaces.HTTPClient, label string) error {
    ctx := context.Background()
    
    brewfilesService := brewfiles.NewService(httpClient)

    request := &brewfiles.UpdateBrewfileRequest{
        Content: "brew \"wget\"\nbrew \"htop\"\nbrew \"curl\"",
    }

    response, err := brewfilesService.UpdateBrewfile(ctx, label, request)
    if err != nil {
        return fmt.Errorf("failed to update brewfile: %w", err)
    }

    fmt.Printf("Success: %s\n", response.Message)
    return nil
}
```

### Creating a Brew Command

```go
import "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewcommands"

func createBrewCommand(httpClient interfaces.HTTPClient) error {
    ctx := context.Background()
    
    brewCommandsService := brewcommands.NewService(httpClient)

    recurrence := "once"
    request := &brewcommands.CreateBrewCommandRequest{
        Arguments:  "install wget",
        Recurrence: &recurrence,
        // DeviceIDs: nil means run on all devices
    }

    response, err := brewCommandsService.CreateBrewCommand(ctx, request)
    if err != nil {
        return fmt.Errorf("failed to create command: %w", err)
    }

    fmt.Printf("Success: %s\n", response.Message)
    return nil
}
```

## Working with Pointer Types

Many optional fields use pointer types. Here are common patterns:

```go
// Creating a pointer inline
deviceSerial := "ABC123"
request.DeviceSerialNumbers = &deviceSerial

// Checking nil pointers when reading
if device.MDMUserOrDeviceName != nil {
    fmt.Printf("MDM Name: %s\n", *device.MDMUserOrDeviceName)
}

// Using helper for optional strings
recurrence := "daily"
request.Recurrence = &recurrence
```

## Next Steps

- Browse the [complete examples directory](.) for all 34 working examples
- Read the [main README](../../README.md) for API reference
- Check the [Workbrew API documentation](https://console.workbrew.com/documentation/api) for endpoint details
