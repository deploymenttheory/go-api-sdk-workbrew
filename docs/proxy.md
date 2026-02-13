# Proxy Support

## What is Proxy Support?

Proxy support allows the SDK to route all HTTP traffic through an intermediate proxy server. This is essential for corporate environments, privacy requirements, or network architectures that mandate proxy usage.

## Why Use a Proxy?

Proxy configuration helps you:

- **Corporate requirements** - Route traffic through corporate proxies
- **Access control** - Comply with network security policies
- **Privacy** - Mask client IP addresses
- **Logging & monitoring** - Centralize traffic inspection
- **Regional access** - Route through specific geographic locations

## When to Configure It

Configure a proxy when:

- Working in corporate environments with mandatory proxies
- Behind firewalls that require proxy for external access
- Need to route traffic through specific geographic locations
- Required for compliance or security policies
- Testing proxy behavior in development

## Basic Example

Here's how to configure a proxy:

```go
package main

import (
    "log"
    "os"

    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
)

func main() {
    // Configure client with HTTP proxy
    workbrewClient, err := client.NewClient(
        os.Getenv("WORKBREW_API_KEY"),
        os.Getenv("WORKBREW_WORKSPACE"),
        client.WithProxy("http://proxy.company.com:8080"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // All requests now route through the proxy
}
```

## Configuration Options

### Option 1: HTTP Proxy

Configure a standard HTTP proxy:

```go
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy("http://proxy.example.com:8080"),
)
```

**When to use:** Most common proxy type, standard corporate proxies

---

### Option 2: HTTPS Proxy

Use an HTTPS proxy for encrypted proxy connections:

```go
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy("https://secure-proxy.example.com:8443"),
)
```

**When to use:** When proxy connection itself needs to be encrypted

---

### Option 3: SOCKS5 Proxy

Configure a SOCKS5 proxy:

```go
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy("socks5://socks-proxy.example.com:1080"),
)
```

**When to use:** SOCKS proxies, SSH tunnels, advanced routing

---

### Option 4: Proxy with Authentication

Use a proxy that requires authentication:

```go
// URL-encoded credentials
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy("http://username:password@proxy.example.com:8080"),
)
```

**When to use:** Proxies that require username/password authentication

**Note:** Credentials in URL are visible in logs. Consider using environment variables.

---

### Option 5: Environment Variable Proxy

Use system environment variables for proxy configuration:

```bash
# Set environment variables
export HTTP_PROXY="http://proxy.example.com:8080"
export HTTPS_PROXY="http://proxy.example.com:8080"
export NO_PROXY="localhost,127.0.0.1"
```

```go
// Client automatically uses HTTP_PROXY/HTTPS_PROXY environment variables
// if no proxy is explicitly configured
workbrewClient, err := client.NewClient(apiKey)
```

**When to use:** System-wide proxy configuration, containerized environments

---

### Option 6: Disable Proxy

Explicitly disable proxy even if environment variables are set:

```go
import "net/http"

// Create transport without proxy
transport := &http.Transport{
    Proxy: nil,
}

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithTransport(transport),
)
```

**When to use:** Override system proxy settings for specific client

---

## Common Scenarios

### Scenario 1: Corporate Proxy

```go
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy("http://proxy.company.com:8080"),
)
```

### Scenario 2: Authenticated Corporate Proxy

```go
// Read credentials from environment
proxyURL := fmt.Sprintf("http://%s:%s@proxy.company.com:8080",
    os.Getenv("PROXY_USER"),
    os.Getenv("PROXY_PASS"))

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy(proxyURL),
)
```

### Scenario 3: Proxy with Custom CA Certificate

```go
// Corporate proxy with SSL inspection
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy("http://proxy.company.com:8080"),
    client.WithRootCertificates("/etc/ssl/certs/corporate-ca.pem"),
)
```

### Scenario 4: SSH Tunnel (SOCKS5)

```bash
# Create SSH tunnel
ssh -D 1080 -N user@jump-host.example.com
```

```go
// Use SOCKS5 proxy
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy("socks5://127.0.0.1:1080"),
)
```

### Scenario 5: Development with Local Proxy

```go
// Use local debugging proxy (e.g., mitmproxy, Burp Suite)
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy("http://127.0.0.1:8080"),
    client.WithInsecureSkipVerify(), // For SSL inspection in dev only!
)
```

## Troubleshooting

### Proxy Connection Failed

**Error:** `proxyconnect tcp: dial tcp: connection refused`

**Solutions:**
- Verify proxy URL is correct
- Check proxy server is running
- Ensure proxy port is accessible
- Test proxy with curl: `curl -x http://proxy:8080 https://api.workbrew.com`

### Proxy Authentication Failed

**Error:** `407 Proxy Authentication Required`

**Solutions:**
```go
// Ensure credentials are correct
proxyURL := fmt.Sprintf("http://%s:%s@proxy.example.com:8080",
    url.QueryEscape(username),
    url.QueryEscape(password))

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy(proxyURL),
)
```

### Certificate Errors with Proxy

**Error:** `x509: certificate signed by unknown authority`

**Solution:** Add proxy's CA certificate
```go
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy("http://proxy.company.com:8080"),
    client.WithRootCertificates("/path/to/proxy-ca.pem"),
)
```

### Proxy Timeout

**Error:** Request times out when using proxy

**Solutions:**
```go
// Increase timeout for proxy connections
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy("http://proxy.example.com:8080"),
    client.WithTimeout(60*time.Second),
)
```

## Security Best Practices

✅ **Do:**
- Use HTTPS proxies when possible
- Store proxy credentials in environment variables
- Validate proxy certificates in production
- Use authenticated proxies
- Monitor proxy connection logs

❌ **Don't:**
- Hardcode proxy credentials in source code
- Disable certificate verification in production
- Use HTTP proxies for sensitive data without encryption
- Ignore proxy authentication failures
- Commit proxy credentials to version control

## Testing Proxy Configuration

```go
func TestProxyConfiguration(t *testing.T) {
    // Test with proxy
    workbrewClient, err := client.NewClient(
        "test-key",
        client.WithProxy("http://proxy.test:8080"),
    )
    assert.NoError(t, err)
}
```

### Testing with Local Proxy

```bash
# Start mitmproxy for local testing
mitmproxy -p 8080

# Or Charles Proxy, Burp Suite, etc.
```

```go
// Configure client to use local proxy
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithProxy("http://127.0.0.1:8080"),
    client.WithInsecureSkipVerify(), // Only for testing!
)
```

## Environment Variable Reference

```bash
# HTTP proxy for non-HTTPS requests
export HTTP_PROXY="http://proxy.example.com:8080"

# HTTPS proxy for HTTPS requests  
export HTTPS_PROXY="http://proxy.example.com:8080"

# Hosts to bypass proxy (comma-separated)
export NO_PROXY="localhost,127.0.0.1,.internal.example.com"

# Proxy authentication
export PROXY_USER="username"
export PROXY_PASS="password"
```

## Related Documentation

- [TLS Configuration](tls-configuration.md) - Configure certificates for proxy SSL inspection
- [Authentication](authentication.md) - API key configuration
- [Timeouts & Retries](timeouts-retries.md) - Adjust timeouts for proxy connections
- [Debugging](debugging.md) - Debug proxy connection issues
