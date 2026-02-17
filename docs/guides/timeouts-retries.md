# Timeouts & Retries

## What are Timeouts and Retries?

Timeouts control how long the SDK waits for an API response before giving up. Retries automatically retry failed requests when they encounter transient errors like network issues or rate limits.

## Why Use Timeouts and Retries?

Proper timeout and retry configuration helps you:

- **Prevent hanging requests** - Avoid waiting indefinitely for responses
- **Handle transient failures** - Automatically recover from temporary network issues
- **Respect rate limits** - Retry with backoff when hitting API quotas
- **Improve reliability** - Make your application more resilient to intermittent failures
- **Control resource usage** - Free up resources from slow or failing requests

## When to Use It

Configure timeouts and retries when:

- Making API calls over unreliable networks
- Running long-lived services that need resilience
- Implementing critical workflows that must handle transient failures
- Dealing with rate-limited APIs
- Running in production environments where reliability is critical

## Basic Example

Here's how to configure timeouts and retries:

```go
package main

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"
)

func main() {
    // Create client with timeout and retry configuration
    workbrewClient, err := client.NewClient(
        os.Getenv("WORKBREW_API_KEY"),
        os.Getenv("WORKBREW_WORKSPACE"),
        client.WithTimeout(30*time.Second),  // 30 second timeout
        client.WithRetryCount(3),             // Retry up to 3 times
    )
    if err != nil {
        log.Fatal(err)
    }

    // Use the client - timeouts and retries are automatic
    devicesService := devices.NewService(workbrewClient)
    result, _, err := devicesService.ListDevices(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Found %d devices", len(result))
}
```

**What happens:**
- If the request takes longer than 30 seconds, it times out
- If the request fails with a retryable error, it automatically retries up to 3 times
- Retries use exponential backoff to avoid overwhelming the server

## Alternative Configuration Options

### Option 1: Custom Timeout

Set a timeout appropriate for your use case:

```go
// Short timeout for quick operations
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithTimeout(10*time.Second),
)

// Longer timeout for large file downloads
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithTimeout(5*time.Minute),
)
```

**When to use:**
- Short timeouts (5-15s): Simple lookups, file reports
- Medium timeouts (30-60s): File uploads, URL scans
- Long timeouts (2-5min): Large file downloads, bulk operations

**Default:** 120 seconds (2 minutes)

---

### Option 2: Retry Configuration

Configure retry behavior for different scenarios.

**Available retry options:**
- `WithRetryCount(n)` - How many times to retry (default: 3)
- `WithRetryWaitTime(d)` - Initial wait time before first retry (default: 2s)
- `WithRetryMaxWaitTime(d)` - Maximum wait time between retries (default: 10s)

The wait time doubles with each retry (exponential backoff) up to the maximum.

```go
import "time"

// Conservative: Few retries, quick backoff
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithRetryCount(2),                      // Retry twice
    client.WithRetryWaitTime(1*time.Second),       // Wait 1s initially
    client.WithRetryMaxWaitTime(5*time.Second),    // Max wait 5s
)

// Aggressive: More retries, longer backoff
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithRetryCount(5),                      // Retry 5 times
    client.WithRetryWaitTime(3*time.Second),       // Wait 3s initially
    client.WithRetryMaxWaitTime(30*time.Second),   // Max wait 30s
)
```

**When to use:**
- Conservative: Rate-limited APIs, quick failures preferred
- Aggressive: Unreliable networks, high importance operations

**Defaults:**
- Retry count: 3
- Wait time: 2 seconds
- Max wait time: 10 seconds

---

### Option 3: Context-Based Timeouts

Use context for per-request timeouts:

```go
func getDevicesWithTimeout(devicesService *devices.Service) error {
    // Create context with 5 second timeout for this specific request
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, _, err := devicesService.ListDevices(ctx)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return fmt.Errorf("request timed out after 5 seconds")
        }
        return err
    }

    log.Printf("Found %d devices", len(*result))
    return nil
}
```

**When to use:** When different operations need different timeouts, or when you want dynamic timeout control.

**Note:** Context timeout takes precedence over client timeout.

---

### Option 4: Disable Retries

Disable retries when you want to fail fast:

```go
// No retries - fail immediately on any error
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithRetryCount(0),
)
```

**When to use:**
- Testing/debugging
- Operations where retries don't make sense (non-idempotent operations)
- When you want to implement custom retry logic

---

## Retry Behavior

### What Gets Retried

The SDK automatically retries:
- ✅ Network errors (connection refused, timeout, etc.)
- ✅ 5xx server errors (500, 502, 503, 504)
- ✅ 429 rate limit errors
- ✅ Request timeout errors

### What Doesn't Get Retried

The SDK does NOT retry:
- ❌ 4xx client errors (400, 401, 403, 404) - these won't succeed on retry
- ❌ Successful responses (2xx)
- ❌ Context cancellation
- ❌ Invalid request configuration

### Exponential Backoff

Retries use exponential backoff with jitter to prevent overwhelming servers:

**With default settings (2s initial, 10s max):**
```
Retry 1: Wait ~2s   (base wait time)
Retry 2: Wait ~4s   (2x backoff)
Retry 3: Wait ~8s   (4x backoff, approaching max)
Retry 4: Wait ~10s  (capped at max wait time)
```

**With custom settings:**
```go
// Example: 5s initial, 30s max
client.WithRetryWaitTime(5*time.Second)
client.WithRetryMaxWaitTime(30*time.Second)

// Results in:
Retry 1: Wait ~5s   (base wait time)
Retry 2: Wait ~10s  (2x backoff)
Retry 3: Wait ~20s  (4x backoff)
Retry 4: Wait ~30s  (capped at max wait time)
```

**Note:** Actual wait times include random jitter (±25%) to avoid thundering herd problems.

## Common Patterns

### Pattern 1: Production-Ready Configuration

```go
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithTimeout(30*time.Second),
    client.WithRetryCount(3),
    client.WithRetryWaitTime(2*time.Second),
    client.WithRetryMaxWaitTime(20*time.Second),
)
```

### Pattern 2: High-Availability Configuration

```go
// More aggressive retries for critical operations
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithTimeout(60*time.Second),
    client.WithRetryCount(5),
    client.WithRetryWaitTime(5*time.Second),
    client.WithRetryMaxWaitTime(60*time.Second),
)
```

### Pattern 3: Fast-Fail Configuration

```go
// Fail quickly for non-critical operations
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithTimeout(10*time.Second),
    client.WithRetryCount(1),
    client.WithRetryWaitTime(1*time.Second),
    client.WithRetryMaxWaitTime(5*time.Second),
)
```

## Handling Timeout Errors

```go
import "context"

result, _, err := devicesService.ListDevices(ctx)
if err != nil {
    // Check if error is due to timeout
    if ctx.Err() == context.DeadlineExceeded {
        log.Println("Request timed out - consider increasing timeout")
        return
    }

    // Check if error is transient (should retry)
    if client.IsServerError(err) || client.IsRateLimited(err) {
        log.Println("Transient error - retries were exhausted")
        return
    }

    // Other error
    log.Printf("Request failed: %v", err)
}
```

## Troubleshooting

### Request Always Times Out

**Symptoms:** Consistent timeout errors even with retries

**Solutions:**
1. Increase timeout: `client.WithTimeout(5*time.Minute)`
2. Check network connectivity
3. Verify Workbrew API is accessible
4. Check for proxy/firewall issues

### Too Many Retries

**Symptoms:** Requests taking very long to fail

**Solutions:**
1. Reduce retry count: `client.WithRetryCount(1)`
2. Reduce wait times:
   ```go
   client.WithRetryWaitTime(1*time.Second)      // Reduce initial wait
   client.WithRetryMaxWaitTime(5*time.Second)   // Reduce max wait
   ```
3. Check if the error is retryable (4xx errors shouldn't be retried)

### Rate Limit Errors Persist

**Symptoms:** Still getting 429 errors after retries

**Solutions:**
1. Increase retry wait times to respect rate limits:
   ```go
   client.WithRetryWaitTime(5*time.Second)      // Longer initial wait
   client.WithRetryMaxWaitTime(60*time.Second)  // Longer max wait
   ```
2. Increase retry count: `client.WithRetryCount(5)`
3. Implement application-level rate limiting
4. Contact Workbrew support about API limits

## Testing

### Simulating Timeouts

```go
func TestTimeout(t *testing.T) {
    // Create client with very short timeout
    workbrewClient, _ := client.NewClient(
        "test-api-key",
        "test-workspace",
        client.WithTimeout(1*time.Millisecond),
    )

    // This will timeout
    devicesService := devices.NewService(workbrewClient)
    ctx := context.Background()
    _, _, err := devicesService.ListDevices(ctx)

    assert.Error(t, err)
}
```

### Testing Retry Logic

```go
func TestRetries(t *testing.T) {
    // Create client with retries disabled for predictable testing
    workbrewClient, _ := client.NewClient(
        "test-api-key",
        "test-workspace",
        client.WithRetryCount(0),
    )

    // Test that errors are returned immediately
    // ... your test code
}
```

## Related Documentation

- [Authentication](authentication.md) - Configure API access
- [Logging](logging.md) - Log timeout and retry events
- [Debugging](debugging.md) - Debug timeout issues
