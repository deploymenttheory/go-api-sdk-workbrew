# Go SDK for Workbrew API

[![Go Report Card](https://goreportcard.com/badge/github.com/deploymenttheory/go-api-sdk-workbrew)](https://goreportcard.com/report/github.com/deploymenttheory/go-api-sdk-workbrew)
[![GoDoc](https://pkg.go.dev/badge/github.com/deploymenttheory/go-api-sdk-workbrew)](https://pkg.go.dev/github.com/deploymenttheory/go-api-sdk-workbrew)
[![License](https://img.shields.io/github/license/deploymenttheory/go-api-sdk-workbrew)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/deploymenttheory/go-api-sdk-workbrew)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/deploymenttheory/go-api-sdk-workbrew)](https://github.com/deploymenttheory/go-api-sdk-workbrew/releases)
[![codecov](https://codecov.io/gh/deploymenttheory/go-api-sdk-workbrew/graph/badge.svg)](https://codecov.io/gh/deploymenttheory/go-api-sdk-workbrew)
[![Tests](https://github.com/deploymenttheory/go-api-sdk-workbrew/workflows/go%20%7C%20Unit%20Tests/badge.svg)](https://github.com/deploymenttheory/go-api-sdk-workbrew/actions)
![Status: Stable](https://img.shields.io/badge/status-stable-brightgreen)

A community Go client library for the [Workbrew API](https://console.workbrew.com/documentation/api) v0.

## Quick Start

Get started quickly with the SDK using the **[Quick Start Guide](docs/guides/quick-start.md)**, which includes:
- Installation instructions
- Your first API call
- Common operations (devices, brewfiles, analytics)
- Error handling patterns
- Response metadata access
- CSV export examples
- Links to configuration guides for production use

## Examples

The [examples directory](examples/workbrew/) contains complete working examples for all SDK features:

Each example includes a complete `main.go` with comments explaining the code. Browse by service:
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

## SDK Services

### Device Management

- **Devices**: List and export device inventory with MDM information
- **Device Groups**: Manage device groups and memberships
- **Analytics**: Retrieve device analytics and metrics

### Brew Package Management

- **Brewfiles**: Full CRUD operations for Brewfile management
- **Brew Commands**: Create and track brew command executions
- **Brew Configurations**: View Homebrew configuration settings
- **Brew Taps**: List available Homebrew taps
- **Formulae**: Browse available Homebrew formulae
- **Casks**: Browse available Homebrew casks

### Security & Compliance

- **Vulnerabilities**: Track package vulnerabilities across your fleet
- **Vulnerability Changes**: Monitor vulnerability status changes over time
- **Events**: Audit log for all system and user activities
- **Licenses**: View and manage software licenses

## HTTP Client Configuration

The SDK includes a powerful HTTP client with production-ready configuration options:

- **[Authentication](docs/guides/authentication.md)** - Secure API key and workspace management
- **[Timeouts & Retries](docs/guides/timeouts-retries.md)** - Configurable timeouts and automatic retry logic
- **[TLS/SSL Configuration](docs/guides/tls-configuration.md)** - Custom certificates, mutual TLS, and security settings
- **[Proxy Support](docs/guides/proxy.md)** - HTTP/HTTPS/SOCKS5 proxy configuration
- **[Custom Headers](docs/guides/custom-headers.md)** - Global and per-request header management
- **[Structured Logging](docs/guides/logging.md)** - Integration with zap for production logging
- **[OpenTelemetry Tracing](docs/guides/opentelemetry.md)** - Distributed tracing and observability
- **[Debug Mode](docs/guides/debugging.md)** - Detailed request/response inspection

## Configuration Options

The SDK client supports extensive configuration through functional options. Below is the complete list of available configuration options grouped by category.

### Basic Configuration

```go
client.WithAPIVersion("v0")              // Set API version
client.WithBaseURL("https://...")        // Custom base URL
client.WithTimeout(30*time.Second)       // Request timeout
client.WithRetryCount(3)                 // Number of retry attempts
```

### TLS/Security

```go
client.WithMinTLSVersion(tls.VersionTLS12)                    // Minimum TLS version
client.WithTLSClientConfig(tlsConfig)                         // Custom TLS configuration
client.WithRootCertificates("/path/to/ca.pem")                // Custom CA certificates
client.WithRootCertificateFromString(caPEM)                   // CA certificate from string
client.WithClientCertificate("/path/cert.pem", "/path/key.pem") // Client certificate (mTLS)
client.WithClientCertificateFromString(certPEM, keyPEM)       // Client cert from string
client.WithInsecureSkipVerify()                               // Skip cert verification (dev only!)
```

### Network

```go
client.WithProxy("http://proxy:8080")    // HTTP/HTTPS/SOCKS5 proxy
client.WithTransport(customTransport)    // Custom HTTP transport
```

### Headers

```go
client.WithUserAgent("MyApp/1.0")                      // Set User-Agent header
client.WithCustomAgent("MyApp")                        // Custom agent identifier
client.WithGlobalHeader("X-Custom-Header", "value")    // Add single global header
client.WithGlobalHeaders(map[string]string{...})       // Add multiple global headers
```

### Observability

```go
client.WithLogger(zapLogger)            // Structured logging with zap
client.WithTracing(otelConfig)          // OpenTelemetry distributed tracing
client.WithDebug()                      // Enable debug mode (dev only!)
```

### Example: Production Configuration

```go
import (
    "crypto/tls"
    "time"
    "go.uber.org/zap"
    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
)

logger, _ := zap.NewProduction()

apiClient, err := client.NewClient(
    "your-api-key",
    "your-workspace",
    client.WithTimeout(30*time.Second),
    client.WithRetryCount(3),
    client.WithLogger(logger),
    client.WithMinTLSVersion(tls.VersionTLS12),
    client.WithGlobalHeader("X-Application-Name", "MyDeviceManager"),
)
```

See the [configuration guides](docs/guides/) for detailed documentation on each option.

## Documentation

- [Workbrew API Documentation](https://console.workbrew.com/documentation/api)
- [GoDoc](https://pkg.go.dev/github.com/deploymenttheory/go-api-sdk-workbrew)
- [Examples Directory](./examples/workbrew/)

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/deploymenttheory/go-api-sdk-workbrew/issues)
- **Documentation**: [API Docs](https://console.workbrew.com/documentation/api)

## Disclaimer

This is an unofficial SDK and is not affiliated with or endorsed by Workbrew.