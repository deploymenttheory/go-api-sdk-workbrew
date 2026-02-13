# Quick Start Guide

Get up and running with the Workbrew Go SDK in minutes.

## Prerequisites

- Go 1.25 or higher
- A Workbrew API key (contact Workbrew for access)

## Installation

```bash
go get github.com/deploymenttheory/go-api-sdk-workbrew
```

## Your First API Call

Here's a complete example that lists devices:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"
)

func main() {
    // Step 1: Create the client with your API key and workspace
    apiClient, err := client.NewClient(
        os.Getenv("WORKBREW_API_KEY"),
        os.Getenv("WORKBREW_WORKSPACE"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Step 2: Create a service (devices, formulae, casks, etc.)
    devicesService := devices.NewService(apiClient)

    // Step 3: Make an API call
    result, resp, err := devicesService.ListDevices(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    // Step 4: Use the results
    fmt.Printf("Status Code: %d\n", resp.StatusCode)
    fmt.Printf("Total Devices: %d\n", len(result))
    if len(result) > 0 {
        fmt.Printf("First Device ID: %s\n", result[0].ID)
        fmt.Printf("First Device Name: %s\n", result[0].Name)
    }
}
```

**Run it:**

```bash
export WORKBREW_API_KEY="your-api-key-here"
go run main.go
```

**Output:**

```text
Status Code: 200
Total Devices: 10
First Device ID: device-123
First Device Name: MacBook-Pro
```

## Common Operations

### List Formulae

```go
import "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/formulae"

formulaeService := formulae.NewService(apiClient)

result, _, err := formulaeService.ListFormulae(context.Background())
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Formulae: %d\n", len(result))
for _, formula := range result {
    fmt.Printf("Formula: %s (version: %s)\n", formula.Name, formula.Version)
}
```

### List Casks

```go
import "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/casks"

casksService := casks.NewService(apiClient)

result, _, err := casksService.ListCasks(context.Background())
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Casks: %d\n", len(result))
for _, cask := range result {
    fmt.Printf("Cask: %s (version: %s)\n", cask.Name, cask.Version)
}
```

### Get Device Groups

```go
import "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devicegroups"

deviceGroupsService := devicegroups.NewService(apiClient)

result, _, err := deviceGroupsService.ListDeviceGroups(context.Background())
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Device Groups: %d\n", len(result))
for _, group := range result {
    fmt.Printf("Group: %s (%d devices)\n", group.Name, group.DeviceCount)
}
```

### Check Vulnerabilities

```go
import "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/vulnerabilities"

vulnService := vulnerabilities.NewService(apiClient)

result, _, err := vulnService.ListVulnerabilities(context.Background())
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Vulnerabilities: %d\n", len(result))
for _, vuln := range result {
    fmt.Printf("CVE: %s (Severity: %s)\n", vuln.CVE, vuln.Severity)
}
```

## Error Handling

Always check for errors and handle common cases:

```go
result, resp, err := devicesService.ListDevices(ctx)

if err != nil {
    // Check for specific error types
    if client.IsNotFound(err) {
        fmt.Println("Resource not found")
        return
    }

    if client.IsUnauthorized(err) {
        fmt.Println("Invalid API key or authentication failed")
        return
    }

    if client.IsForbidden(err) {
        fmt.Println("Access forbidden - check your permissions")
        return
    }

    if client.IsRateLimited(err) {
        fmt.Println("Rate limit exceeded - wait before retrying")
        return
    }

    // Other errors
    log.Fatal(err)
}

// Success - use the result
fmt.Printf("Found %d devices\n", len(result))
```

## Response Metadata

Every API call returns response metadata:

```go
result, resp, err := devicesService.ListDevices(ctx)

// Access response metadata
fmt.Printf("Status Code: %d\n", resp.StatusCode)
fmt.Printf("Request Duration: %v\n", resp.Duration)
fmt.Printf("Response Size: %d bytes\n", resp.Size)
fmt.Printf("Received At: %v\n", resp.ReceivedAt)

// Check if response was successful
if client.IsResponseSuccess(resp) {
    fmt.Println("Request successful")
}
```

## Next Steps

### Production Configuration

For production use, configure the client with appropriate settings:

```go
import (
    "time"
    "go.uber.org/zap"
)

logger, _ := zap.NewProduction()

apiClient, err := client.NewClient(
    os.Getenv("WORKBREW_API_KEY"),
    os.Getenv("WORKBREW_WORKSPACE"),
    client.WithTimeout(30*time.Second),
    client.WithRetryCount(3),
    client.WithLogger(logger),
)
```

**Learn more:**

- **[Authentication](authentication.md)** - Secure API key management
- **[Timeouts & Retries](timeouts-retries.md)** - Configure resilience
- **[Structured Logging](logging.md)** - Production logging with zap

### Advanced Features

Enhance your integration with advanced client features:

**Observability:**

- **[OpenTelemetry Tracing](opentelemetry.md)** - Distributed tracing for monitoring
- **[Debug Mode](debugging.md)** - Detailed request/response inspection

**Network Configuration:**

- **[TLS/SSL Configuration](tls-configuration.md)** - Custom certificates and mutual TLS
- **[Proxy Support](proxy.md)** - Route traffic through proxies
- **[Custom Headers](custom-headers.md)** - Add tracking or metadata headers

### API Coverage

Explore all available services:

**Device Management:**
- Devices - List, view, and manage devices
- Device Groups - Organize devices into groups
- Brew Configurations - Manage Homebrew configurations

**Brew Package Management:**
- Formulae - Manage Homebrew formulae
- Casks - Manage Homebrew casks
- Brew Files - Manage Brewfiles
- Brew Taps - Manage custom taps
- Brew Commands - Execute brew commands

**Licensing & Security:**
- Licenses - Manage software licenses
- Vulnerabilities - View and track CVEs
- Vulnerability Changes - Monitor vulnerability updates

**Analytics & Events:**
- Analytics - View usage analytics
- Events - Track system events and activities

### Examples

Check out the [examples directory](../../examples/workbrew/) for complete working examples:

- Device management
- Formulae and Cask operations
- Brewfile management
- Vulnerability tracking
- License management
- Event monitoring

## Troubleshooting

### "Invalid API Key" Error

```go
// Verify your API key is set correctly
apiKey := os.Getenv("WORKBREW_API_KEY")
if apiKey == "" {
    log.Fatal("WORKBREW_API_KEY environment variable not set")
}

// Check for authentication errors
if err != nil && client.IsUnauthorized(err) {
    log.Fatal("Invalid API key - check your credentials")
}
```

### "Rate Limit Exceeded" Error

```go
// Handle rate limit errors
if client.IsRateLimited(err) {
    log.Println("Rate limit exceeded - waiting 60 seconds")
    time.Sleep(60 * time.Second)
    // Retry request
}
```

### "Forbidden" Error

```go
// Handle forbidden errors
if client.IsForbidden(err) {
    log.Println("Access forbidden - check workspace ID and permissions")
    
    // Verify workspace is set
    apiClient.SetWorkspace("correct-workspace-id")
}
```

## Getting Help

- **[Full Documentation](../../README.md)** - Complete SDK documentation
- **[Workbrew Documentation](https://docs.workbrew.com)** - Workbrew API documentation
- **[GitHub Issues](https://github.com/deploymenttheory/go-api-sdk-workbrew/issues)** - Report bugs or request features
- **[GoDoc](https://pkg.go.dev/github.com/deploymenttheory/go-api-sdk-workbrew)** - Package documentation

## Complete Example

Here's a complete example with error handling, logging, and workspace management:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"
    "go.uber.org/zap"
)

func main() {
    // Initialize logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    // Create client with production settings
    apiClient, err := client.NewClient(
        os.Getenv("WORKBREW_API_KEY"),
        os.Getenv("WORKBREW_WORKSPACE"),
        client.WithTimeout(30*time.Second),
        client.WithRetryCount(3),
        client.WithLogger(logger),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create service
    devicesService := devices.NewService(apiClient)

    // Make API call with context
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    result, resp, err := devicesService.ListDevices(ctx)
    if err != nil {
        logger.Error("Failed to list devices",
            zap.Error(err),
            zap.Int("status_code", resp.StatusCode),
        )
        log.Fatal(err)
    }

    // Log response metadata
    logger.Info("API call successful",
        zap.Int("status_code", resp.StatusCode),
        zap.Duration("duration", resp.Duration),
        zap.Int64("size", resp.Size),
    )

    // Use results
    fmt.Printf("Total Devices: %d\n", len(result))
    for _, device := range result {
        fmt.Printf("Device: %s (ID: %s)\n", device.Name, device.ID)
    }
}
```

---

**Ready to build?** Start with this quick start and explore the configuration guides to customize the SDK for your needs!
