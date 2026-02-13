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

    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
)

func main() {
    // Step 1: Create the client with your API key and workspace
    workbrewClient, err := workbrew.NewClient(
        os.Getenv("WORKBREW_API_KEY"),
        os.Getenv("WORKBREW_WORKSPACE"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Step 2: Make an API call using the service
    result, resp, err := workbrewClient.Devices.ListDevices(
        context.Background(),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Step 3: Use the results
    fmt.Printf("Status Code: %d\n", resp.StatusCode)
    fmt.Printf("Total Devices: %d\n", len(*result))
    if len(*result) > 0 {
        fmt.Printf("First Device Serial: %s\n", (*result)[0].SerialNumber)
        if (*result)[0].MDMUserOrDeviceName != nil {
            fmt.Printf("First Device Name: %s\n", *(*result)[0].MDMUserOrDeviceName)
        }
    }
}
```

**Run it:**

```bash
export WORKBREW_API_KEY="your-api-key-here"
export WORKBREW_WORKSPACE="your-workspace-id"
go run main.go
```

**Output:**

```text
Status Code: 200
Total Devices: 10
First Device Serial: TC6R2DHVHG
First Device Name: MacBook-Pro
```

## Common Operations

### List Formulae

```go
result, _, err := workbrewClient.Formulae.ListFormulae(
    context.Background(),
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Formulae: %d\n", len(*result))
for _, formula := range *result {
    version := "N/A"
    if formula.HomebrewCoreVersion != nil {
        version = *formula.HomebrewCoreVersion
    }
    fmt.Printf("Formula: %s (version: %s)\n", formula.Name, version)
}
```

### List Casks

```go
result, _, err := workbrewClient.Casks.ListCasks(
    context.Background(),
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Casks: %d\n", len(*result))
for _, cask := range *result {
    version := "N/A"
    if cask.HomebrewCaskVersion != nil {
        version = *cask.HomebrewCaskVersion
    }
    fmt.Printf("Cask: %s (version: %s)\n", cask.Name, version)
}
```

### Get Device Groups

```go
result, _, err := workbrewClient.DeviceGroups.ListDeviceGroups(
    context.Background(),
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Device Groups: %d\n", len(*result))
for _, group := range *result {
    deviceCount := 0
    if group.DeviceCount != nil {
        deviceCount = *group.DeviceCount
    }
    fmt.Printf("Group: %s (%d devices)\n", group.Name, deviceCount)
}
```

### Check Vulnerabilities

```go
result, _, err := workbrewClient.Vulnerabilities.ListVulnerabilities(
    context.Background(),
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Vulnerabilities: %d\n", len(*result))
for _, vuln := range *result {
    fmt.Printf("CVE: %s (Severity: %s)\n", vuln.CVE, vuln.Severity)
}
```


## Response Metadata

Every API call returns response metadata:

```go
result, resp, err := workbrewClient.Devices.ListDevices(ctx)

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

## Complete Production Example

Here's a complete example bringing together all the concepts with error handling, logging, and configuration:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
    "go.uber.org/zap"
)

func main() {
    // Initialize logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    // Create client with production settings
    workbrewClient, err := workbrew.NewClient(
        os.Getenv("WORKBREW_API_KEY"),
        os.Getenv("WORKBREW_WORKSPACE"),
        client.WithTimeout(30*time.Second),
        client.WithRetryCount(3),
        client.WithLogger(logger),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Make API call with context
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    result, resp, err := workbrewClient.Devices.ListDevices(ctx)

    // Handle errors
    if err != nil {
        if client.IsUnauthorized(err) {
            logger.Warn("Invalid API key or authentication failed")
            return
        }
        logger.Error("API call failed", zap.Error(err))
        log.Fatal(err)
    }

    // Log response metadata
    logger.Info("API call successful",
        zap.Int("status_code", resp.StatusCode),
        zap.Duration("duration", resp.Duration),
        zap.Int64("size", resp.Size),
    )

    // Use results
    fmt.Printf("Total Devices: %d\n", len(*result))
    for _, device := range *result {
        name := "Unknown"
        if device.MDMUserOrDeviceName != nil {
            name = *device.MDMUserOrDeviceName
        }
        fmt.Printf("Device: %s (Serial: %s)\n", name, device.SerialNumber)
    }
}
```

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

## Next Steps

### Client Configuration Options

Learn different ways to create and configure a Workbrew client:

**[📚 Client Building Examples](../../examples/workbrew/_build_client/)**

Four comprehensive examples showing different client setup scenarios:

1. **[Basic Client](../../examples/workbrew/_build_client/new_client/)** - Simplest setup with minimal configuration
2. **[Environment-Based Client](../../examples/workbrew/_build_client/new_client_with_env_var/)** - 12-factor app compliant configuration
3. **[Production Client with Logger](../../examples/workbrew/_build_client/new_client_with_logger/)** - Custom timeouts, retries, and structured logging
4. **[OpenTelemetry Client](../../examples/workbrew/_build_client/new_client_with_open_telemetry/)** - Distributed tracing for observability

Each example includes complete working code, when to use each approach, and security best practices.

### Configuration Guides

**Essential:**
- **[Authentication](authentication.md)** - Secure API key management
- **[Timeouts & Retries](timeouts-retries.md)** - Configure resilience
- **[Structured Logging](logging.md)** - Production logging with zap

**Advanced:**
- **[OpenTelemetry Tracing](opentelemetry.md)** - Distributed tracing for monitoring
- **[Debug Mode](debugging.md)** - Detailed request/response inspection
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

### More Examples

Check out the [examples directory](../../examples/workbrew/) for complete working examples:
- Device management
- Formulae and Cask operations
- Brewfile management
- Vulnerability tracking
- License management
- Event monitoring

## Getting Help

- **[Full Documentation](../../README.md)** - Complete SDK documentation
- **[Workbrew Documentation](https://docs.workbrew.com)** - Workbrew API documentation
- **[GitHub Issues](https://github.com/deploymenttheory/go-api-sdk-workbrew/issues)** - Report bugs or request features
- **[GoDoc](https://pkg.go.dev/github.com/deploymenttheory/go-api-sdk-workbrew)** - Package documentation

---

**Ready to build?** Start with this quick start and explore the configuration guides to customize the SDK for your needs!
