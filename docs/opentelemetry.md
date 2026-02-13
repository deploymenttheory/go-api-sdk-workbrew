# OpenTelemetry Tracing

## What is OpenTelemetry Tracing?

OpenTelemetry tracing provides distributed tracing capabilities for your Workbrew API calls. It automatically captures detailed information about each HTTP request, including timing, status codes, errors, and request/response metadata.

## Why Use OpenTelemetry?

OpenTelemetry tracing helps you:

- **Monitor performance** - Track how long API calls take and identify bottlenecks
- **Debug issues** - See the complete flow of requests across your application
- **Track errors** - Automatically capture and report API errors with full context
- **Improve observability** - Integrate with platforms like Jaeger, Zipkin, DataDog, Honeycomb, etc.
- **Understand dependencies** - Visualize how your application interacts with Workbrew

## When to Use It

Use OpenTelemetry tracing when:

- Running in production environments where observability is critical
- Debugging complex issues that span multiple services
- Monitoring API performance and identifying slow requests
- Tracking error rates and failure patterns
- Meeting compliance or SLA requirements for observability

## Basic Example

Here's a simple example showing how to enable OpenTelemetry tracing:

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
    // Step 1: Initialize OpenTelemetry exporter
    exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
    if err != nil {
        log.Fatal(err)
    }

    // Step 2: Create tracer provider
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
    )
    defer tp.Shutdown(context.Background())

    // Step 3: Set as global tracer provider
    otel.SetTracerProvider(tp)

    // Step 4: Create Workbrew client with tracing enabled
    workbrewClient, err := client.NewClient(
        os.Getenv("WORKBREW_API_KEY"),
        os.Getenv("WORKBREW_WORKSPACE"),
        client.WithTracing(nil), // nil uses global tracer provider
    )
    if err != nil {
        log.Fatal(err)
    }

    // Step 5: Use the client normally - tracing happens automatically!
    devicesService := devices.NewService(workbrewClient)
    result, _, err := devicesService.ListDevices(context.Background())
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }

    log.Printf("Found %d devices", len(result))

    // Traces are automatically exported - check your console output!
}
```

**What you get:**

- All HTTP requests are automatically traced
- Spans include method, URL, status code, timing
- Errors are automatically recorded
- Zero code changes needed in your business logic

## Alternative Configuration Options

### Option 1: Using Default Configuration

The simplest approach uses the global OpenTelemetry tracer provider:

```go
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithTracing(nil), // Uses otel.GetTracerProvider()
)
```

**When to use:** For most applications where you've already configured OpenTelemetry globally.

---

### Option 2: Custom Tracer Provider

Provide a specific tracer provider for more control:

```go
// Create your own tracer provider
myTracerProvider := trace.NewTracerProvider(
    trace.WithBatcher(myExporter),
    trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(0.1))), // Sample 10%
)

// Configure the client to use it
otelConfig := &client.OTelConfig{
    TracerProvider: myTracerProvider,
    ServiceName:    "my-workbrew-client",
}

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithTracing(otelConfig),
)
```

**When to use:** When you need different tracing configurations for different clients, or want to override the global tracer provider.

---

### Option 3: Custom Span Naming

Customize how spans are named for better organization in your tracing UI:

```go
otelConfig := &client.OTelConfig{
    SpanNameFormatter: func(operation string, req *http.Request) string {
        // Custom format: "Workbrew: GET /files/{hash}"
        return fmt.Sprintf("Workbrew: %s %s", req.Method, req.URL.Path)
    },
}

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithTracing(otelConfig),
)
```

**When to use:** When you want more descriptive span names in your tracing dashboard (e.g., Jaeger, Zipkin).

---

### Option 4: Custom Propagators

Control how trace context is propagated across service boundaries:

```go
import "go.opentelemetry.io/otel/propagation"

otelConfig := &client.OTelConfig{
    Propagators: propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ),
}

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithTracing(otelConfig),
)
```

**When to use:** When integrating with systems that use specific trace context formats (W3C Trace Context, B3, etc.).

---

## Integration with Popular Backends

### Jaeger

```go
import (
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/trace"
)

exporter, _ := jaeger.New(jaeger.WithCollectorEndpoint(
    jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
))

tp := trace.NewTracerProvider(trace.WithBatcher(exporter))
otel.SetTracerProvider(tp)

workbrewClient, _ := client.NewClient(apiKey, client.WithTracing(nil))
```

### OTLP (OpenTelemetry Protocol)

```go
import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

exporter, _ := otlptracegrpc.New(context.Background(),
    otlptracegrpc.WithEndpoint("otel-collector:4317"),
    otlptracegrpc.WithInsecure(),
)

tp := trace.NewTracerProvider(trace.WithBatcher(exporter))

workbrewClient, _ := client.NewClient(apiKey, client.WithTracing(nil))
```

## What Gets Traced

Each HTTP request creates a span with the following information:

| Attribute | Description | Example |
|-----------|-------------|---------|
| `http.method` | HTTP method | `GET`, `POST` |
| `http.url` | Full URL | `https://api.workbrew.com/devices` |
| `http.status_code` | Response status | `200`, `404`, `429` |
| `http.request_content_length` | Request size in bytes | `1024` |
| `http.response_content_length` | Response size in bytes | `4096` |
| Span duration | Request timing | `245ms` |
| Span status | Success or error | `Ok`, `Error` |

All attributes follow [OpenTelemetry semantic conventions](https://opentelemetry.io/docs/specs/semconv/http/) for HTTP clients.

## Disabling Tracing

To disable tracing, simply omit the `WithTracing()` option:

```go
// No tracing - client works normally without instrumentation
workbrewClient, err := client.NewClient(apiKey)
```

## Performance Considerations

- **Minimal overhead**: OpenTelemetry adds microseconds of latency per request
- **Async export**: Spans are batched and exported in the background
- **Sampling**: Use sampling to reduce overhead in high-traffic scenarios:
  ```go
  trace.WithSampler(trace.TraceIDRatioBased(0.1)) // Sample 10%
  ```
- **No-op when disabled**: Without `WithTracing()`, there's zero tracing overhead

## Complete Example

See [examples/workbrew/observability/tracing/main.go](../../examples/workbrew/observability/tracing/main.go) for a complete working example.

## Related Documentation

- [Structured Logging](logging.md) - Combine tracing with logging for complete observability
- [Error Handling](error-handling.md) - How errors are captured in traces
- [Context Support](context.md) - Using context for trace propagation
- [OpenTelemetry Go Documentation](https://opentelemetry.io/docs/languages/go/)
