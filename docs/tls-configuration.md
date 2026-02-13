# TLS/SSL Configuration

## What is TLS Configuration?

TLS (Transport Layer Security) configuration controls how the SDK establishes secure HTTPS connections to the Workbrew API. You can customize certificate validation, use mutual TLS, and set minimum TLS versions.

## Why Configure TLS?

TLS configuration helps you:

- **Use custom certificates** - Work with private CAs or self-signed certificates
- **Enable mutual TLS** - Use client certificates for enhanced authentication
- **Meet security requirements** - Enforce minimum TLS versions (TLS 1.2, 1.3)
- **Corporate environments** - Integrate with enterprise certificate infrastructures
- **Compliance** - Meet regulatory requirements for encryption

## When to Configure It

Configure TLS when:

- Using private or internal CAs
- Required to use client certificates
- Working behind corporate proxies with SSL inspection
- Meeting compliance requirements (PCI-DSS, HIPAA, etc.)
- Enforcing specific TLS versions
- Testing with self-signed certificates (development only)

## Basic Example

Here's how to configure basic TLS settings:

```go
package main

import (
    "crypto/tls"
    "log"
    "os"

    "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
)

func main() {
    // Create client with minimum TLS 1.2
    workbrewClient, err := client.NewClient(
        os.Getenv("WORKBREW_API_KEY"),
        os.Getenv("WORKBREW_WORKSPACE"),
        client.WithMinTLSVersion(tls.VersionTLS12),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Use client normally
    // All connections now use TLS 1.2 or higher
}
```

## Configuration Options

### Option 1: Minimum TLS Version

Enforce a minimum TLS version for all connections:

```go
import "crypto/tls"

// Require TLS 1.2 or higher (recommended)
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithMinTLSVersion(tls.VersionTLS12),
)

// Require TLS 1.3 (most secure)
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithMinTLSVersion(tls.VersionTLS13),
)
```

**When to use:**
- TLS 1.2: Industry standard, widely compatible
- TLS 1.3: Maximum security, modern systems

**Default:** System default (usually TLS 1.2+)

---

### Option 2: Custom Root Certificates

Add custom CA certificates for server validation:

```go
// From file
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithRootCertificates(
        "/path/to/ca-cert.pem",
        "/path/to/another-ca.pem",
    ),
)

// From string
caCertPEM := `-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJAKL0UG+mRK...
-----END CERTIFICATE-----`

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithRootCertificateFromString(caCertPEM),
)
```

**When to use:**
- Working with private/internal CAs
- Corporate environments with custom certificates
- Self-signed certificates (development/testing)

---

### Option 3: Client Certificates (Mutual TLS)

Use client certificates for authentication:

```go
// From files
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithClientCertificate(
        "/path/to/client-cert.pem",
        "/path/to/client-key.pem",
    ),
)

// From strings
clientCertPEM := `-----BEGIN CERTIFICATE-----...`
clientKeyPEM := `-----BEGIN PRIVATE KEY-----...`

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithClientCertificateFromString(clientCertPEM, clientKeyPEM),
)
```

**When to use:**
- Mutual TLS (mTLS) requirements
- Enhanced authentication beyond API keys
- Enterprise security policies

---

### Option 4: Custom TLS Configuration

Full control over TLS settings:

```go
import "crypto/tls"

tlsConfig := &tls.Config{
    MinVersion:         tls.VersionTLS12,
    MaxVersion:         tls.VersionTLS13,
    InsecureSkipVerify: false, // NEVER use true in production
    CipherSuites: []uint16{
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
    },
}

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithTLSClientConfig(tlsConfig),
)
```

**When to use:**
- Specific cipher suite requirements
- Custom security policies
- Advanced TLS configuration needs

---

### Option 5: Disable Certificate Verification (DEVELOPMENT ONLY)

⚠️ **WARNING**: Only for development/testing with self-signed certificates!

```go
// NEVER use this in production!
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithInsecureSkipVerify(),
)
```

**When to use:** Testing with self-signed certificates in development

**⚠️ Security Risk:** Disables certificate validation, vulnerable to MITM attacks

---

## Common Scenarios

### Scenario 1: Corporate Proxy with SSL Inspection

```go
// Load corporate CA certificate
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithRootCertificates("/etc/ssl/certs/corporate-ca.pem"),
    client.WithProxy("http://proxy.company.com:8080"),
)
```

### Scenario 2: Mutual TLS Authentication

```go
// Use client certificate + API key
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithClientCertificate(
        "/path/to/client.crt",
        "/path/to/client.key",
    ),
    client.WithMinTLSVersion(tls.VersionTLS12),
)
```

### Scenario 3: High Security Environment

```go
import "crypto/tls"

// Enforce TLS 1.3 with strong ciphers
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS13,
    CipherSuites: []uint16{
        tls.TLS_AES_256_GCM_SHA384,
        tls.TLS_CHACHA20_POLY1305_SHA256,
    },
}

workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithTLSClientConfig(tlsConfig),
)
```

### Scenario 4: Development with Self-Signed Certificates

```go
// Load self-signed CA for testing
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithRootCertificates("/tmp/self-signed-ca.pem"),
)

// Or temporarily skip verification (NOT for production)
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithInsecureSkipVerify(),
)
```

## Troubleshooting

### Certificate Verification Failed

**Error:** `x509: certificate signed by unknown authority`

**Solutions:**
```go
// Add the CA certificate
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithRootCertificates("/path/to/ca-cert.pem"),
)
```

### TLS Handshake Failure

**Error:** `tls: handshake failure`

**Solutions:**
```go
// Check minimum TLS version
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithMinTLSVersion(tls.VersionTLS12),
)

// Or use custom TLS config with compatible cipher suites
```

### Client Certificate Errors

**Error:** `tls: bad certificate`

**Solutions:**
- Verify certificate and key match
- Check certificate is not expired
- Ensure private key is in correct format

```go
// Reload certificates
workbrewClient, err := client.NewClient(
    apiKey,
    workspace,
    client.WithClientCertificate(
        "/path/to/valid-cert.pem",
        "/path/to/valid-key.pem",
    ),
)
```

## Security Best Practices

✅ **Do:**
- Use TLS 1.2 or higher
- Validate certificates in production
- Protect private keys with file permissions (chmod 600)
- Use strong cipher suites
- Rotate certificates before expiration

❌ **Don't:**
- Disable certificate verification in production
- Use self-signed certificates in production
- Commit private keys to version control
- Use outdated TLS versions (1.0, 1.1)
- Ignore certificate expiration warnings

## Testing TLS Configuration

```go
func TestTLSConfiguration(t *testing.T) {
    // Test with custom CA
    workbrewClient, err := client.NewClient(
        "test-key",
        "test-workspace",
        client.WithRootCertificates("/path/to/test-ca.pem"),
    )
    assert.NoError(t, err)

    // Test minimum TLS version
    workbrewClient, err = client.NewClient(
        "test-key",
        "test-workspace",
        client.WithMinTLSVersion(tls.VersionTLS12),
    )
    assert.NoError(t, err)
}
```

## Related Documentation

- [Proxy Support](proxy.md) - Configure proxies (often used with custom CAs)
- [Authentication](authentication.md) - API key configuration
- [Debugging](debugging.md) - Debug TLS handshake issues
- [Go TLS Package](https://pkg.go.dev/crypto/tls)
