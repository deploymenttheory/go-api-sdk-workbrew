# Custom Headers

## What are Custom Headers?

Custom headers allow you to add additional HTTP headers to all requests (global headers) or specific requests. This is useful for adding metadata, tracking identifiers, or custom authentication schemes.

## Why Use Custom Headers?

Custom headers help you:

- **Track requests** - Add request IDs for debugging and tracing
- **Add metadata** - Include application version, user context, etc.
- **Custom authentication** - Add additional auth headers beyond API key
- **Compliance** - Include required headers for auditing
- **Integration** - Pass data to intermediate proxies or gateways

## When to Configure Them

Add custom headers when:

- Need to correlate requests across systems
- Adding application metadata for monitoring
- Working with API gateways that require specific headers
- Implementing custom authentication schemes
- Meeting compliance requirements for request tracking

## Basic Example

Here's how to add custom headers:

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"
)

func main() {
    // Add a global header to all requests
    workbrewClient, err := client.NewClient(
        os.Getenv("WORKBREW_API_KEY"),
        os.Getenv("WORKBREW_WORKSPACE"),
        client.WithGlobalHeader("X-Application-Name", "MyDeviceManager"),
        client.WithGlobalHeader("X-Application-Version", "1.0.0"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // All requests now include these headers
    devicesService := devices.NewService(workbrewClient)
    result, _, err := devicesService.ListDevices(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Found %d devices", len(result))
}
```

## Configuration Options

### Option 1: Single Global Header

Add one header that applies to all requests:

```go
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeader("X-Request-ID", requestID),
)
```

**When to use:** Adding a single tracking or metadata header

---

### Option 2: Multiple Global Headers

Add multiple headers at once:

```go
headers := map[string]string{
    "X-Application-Name":    "MyApp",
    "X-Application-Version": "1.0.0",
    "X-Environment":         "production",
    "X-Region":              "us-east-1",
}

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeaders(headers),
)
```

**When to use:** Adding multiple metadata or tracking headers

---

### Option 3: Chain Multiple Headers

Add headers one at a time with multiple options:

```go
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeader("X-App-Name", "MyApp"),
    client.WithGlobalHeader("X-App-Version", "1.0.0"),
    client.WithGlobalHeader("X-User-ID", userID),
)
```

**When to use:** Building headers conditionally or from different sources

---

### Option 4: Dynamic Headers

Generate headers dynamically:

```go
import "github.com/google/uuid"

// Generate unique request ID for each client
requestID := uuid.New().String()

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeader("X-Request-ID", requestID),
    client.WithGlobalHeader("X-Timestamp", time.Now().Format(time.RFC3339)),
)
```

**When to use:** Headers that change per client instance

---

### Option 5: Override Global Headers Per Request

Global headers can be overridden on a per-request basis:

```go
// Note: Per-request header override is handled by the service methods
// Global headers serve as defaults
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeader("X-Priority", "normal"),
)

// Individual service methods handle per-request headers internally
```

**When to use:** Different header values for specific requests

---

## Common Use Cases

### Use Case 1: Request Tracking

```go
import "github.com/google/uuid"

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeader("X-Request-ID", uuid.New().String()),
    client.WithGlobalHeader("X-Correlation-ID", correlationID),
)
```

### Use Case 2: Application Metadata

```go
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeaders(map[string]string{
        "X-Application-Name":    "SecurityScanner",
        "X-Application-Version": version,
        "X-Build-Number":        buildNumber,
        "X-Environment":         env,
    }),
)
```

### Use Case 3: User Context

```go
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeaders(map[string]string{
        "X-User-ID":      userID,
        "X-Organization": orgID,
        "X-Role":         role,
    }),
)
```

### Use Case 4: Custom Authentication

```go
// Additional auth header beyond API key
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeader("X-Custom-Token", customToken),
    client.WithGlobalHeader("X-Auth-Type", "dual"),
)
```

### Use Case 5: API Gateway Integration

```go
// Headers required by API gateway
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeaders(map[string]string{
        "X-Gateway-Key":  gatewayKey,
        "X-API-Version":  "v3",
        "X-Client-Type":  "sdk-go",
    }),
)
```

## Header Naming Conventions

### Standard Patterns

```go
// Use X- prefix for custom headers (traditional)
"X-Application-Name"
"X-Request-ID"
"X-User-Context"

// Or modern convention without X-
"Application-Name"
"Request-ID"
"User-Context"

// Use kebab-case for readability
"X-Application-Name"  // Good
"X-APPLICATION-NAME"  // Less readable
"X_Application_Name"  // Don't use underscores
```

### Reserved Headers

Some headers are automatically set by the SDK:
- `User-Agent` - Set via `WithUserAgent()` or `WithCustomAgent()`
- `Authorization` - Set via API key in `NewClient()`
- `X-Workbrew-API-Version` - Set via API version configuration
- `Accept-Encoding` - Set automatically for compression
- `Content-Type` - Set automatically based on request body

**Note:** Workspace is part of the base URL (`/workspaces/{workspace}`), not a header.

## Troubleshooting

### Headers Not Appearing in Requests

**Problem:** Custom headers don't appear in requests

**Solutions:**
```go
// Verify headers are set
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeader("X-Test", "value"),
    client.WithDebug(), // Enable debug mode to see headers
)

// Check logs for header values
```

### Header Values Being Overwritten

**Problem:** Global headers are replaced by default values

**Note:** Per-request headers take precedence over global headers

**Solution:** Ensure you're setting the correct header names:
```go
// Global header
client.WithGlobalHeader("X-Custom-Header", "global-value")

// If a service method sets the same header, it will override global
```

### Special Characters in Header Values

**Problem:** Header values with special characters cause errors

**Solution:** URL encode special characters:
```go
import "net/url"

encodedValue := url.QueryEscape("value with spaces")
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeader("X-Custom", encodedValue),
)
```

## Security Considerations

✅ **Do:**
- Validate header values before setting
- Use headers for non-sensitive metadata
- Document custom headers in API integration guides
- Use standard header naming conventions
- Keep header values concise

❌ **Don't:**
- Put sensitive data in headers (passwords, tokens, PII)
- Use headers for large data payloads
- Include credentials unless encrypted in transit
- Log header values that contain secrets
- Use non-standard or ambiguous header names

## Testing with Custom Headers

```go
func TestCustomHeaders(t *testing.T) {
    // Test single header
    workbrewClient, err := client.NewClient(
        "test-key",
        "test-workspace",
        client.WithGlobalHeader("X-Test", "value"),
    )
    assert.NoError(t, err)

    // Test multiple headers
    headers := map[string]string{
        "X-Test-1": "value1",
        "X-Test-2": "value2",
    }
    workbrewClient, err = client.NewClient(
        "test-key",
        client.WithGlobalHeaders(headers),
    )
    assert.NoError(t, err)
}
```

### Inspecting Headers in Tests

```go
// Use debug mode to see actual headers sent
workbrewClient, err := client.NewClient(
    "test-key",
    client.WithGlobalHeader("X-Test", "value"),
    client.WithDebug(),
)

// Or inspect via HTTP mock
```

## Examples by Language/Framework

### With OpenTelemetry

```go
import "go.opentelemetry.io/otel/trace"

// Add trace context to headers
spanCtx := trace.SpanContextFromContext(ctx)
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithGlobalHeader("X-Trace-ID", spanCtx.TraceID().String()),
    client.WithGlobalHeader("X-Span-ID", spanCtx.SpanID().String()),
)
```

### With Request ID Propagation

```go
// Propagate request ID from incoming HTTP request
func handleRequest(w http.ResponseWriter, r *http.Request) {
    requestID := r.Header.Get("X-Request-ID")
    if requestID == "" {
        requestID = uuid.New().String()
    }

    workbrewClient, _ := client.NewClient(
        apiKey,
        client.WithGlobalHeader("X-Request-ID", requestID),
    )

    // Use client...
}
```

## Related Documentation

- [Authentication](authentication.md) - Configure API key (also a header)
- [Debugging](debugging.md) - View headers in debug output
- [Logging](logging.md) - Log header values (be careful with sensitive data)
- [OpenTelemetry](opentelemetry.md) - Integrate trace context in headers
