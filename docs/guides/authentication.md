# Authentication

## What is API Key Authentication?

The Workbrew SDK uses Bearer token authentication to securely access the Workbrew API. Your API key is sent as a Bearer token with every request to identify and authorize your application.

## Why Use Proper Authentication?

Proper authentication handling helps you:

- **Secure your credentials** - Avoid hardcoding API keys in source code
- **Prevent unauthorized access** - Ensure only valid API keys are used
- **Support multiple environments** - Use different keys for dev, staging, and production
- **Audit usage** - Track which keys are making requests
- **Simplify key rotation** - Create new clients with new keys when needed

## When to Use It

Always use proper authentication when:

- Accessing the Workbrew API from any application
- Deploying to production environments
- Sharing code in version control systems
- Running automated tests or CI/CD pipelines
- Managing multiple Workbrew workspaces or API tiers

## Basic Example

Here's the recommended way to authenticate with the Workbrew SDK:

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
    // Step 1: Get API key and workspace from environment variables (recommended)
    apiKey := os.Getenv("WORKBREW_API_KEY")
    if apiKey == "" {
        log.Fatal("WORKBREW_API_KEY environment variable is required")
    }
    workspace := os.Getenv("WORKBREW_WORKSPACE")
    if workspace == "" {
        log.Fatal("WORKBREW_WORKSPACE environment variable is required")
    }

    // Step 2: Create client with API key and workspace
    workbrewClient, err := client.NewClient(apiKey, workspace)
    if err != nil {
        log.Fatal(err)
    }

    // Step 3: Use the client - authentication is automatic
    devicesService := devices.NewService(workbrewClient)
    result, _, err := devicesService.ListDevices(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Found %d devices", len(result))
}
```

**Run the example:**

```bash
export WORKBREW_API_KEY="your-api-key-here"
export WORKBREW_WORKSPACE="your-workspace-id"
go run main.go
```

## Alternative Configuration Options

### Option 1: Environment Variables (Recommended)

Store API keys in environment variables for security and flexibility:

```go
// Production: Read from environment
apiKey := os.Getenv("WORKBREW_API_KEY")
workspace := os.Getenv("WORKBREW_WORKSPACE")

workbrewClient, err := client.NewClient(apiKey, workspace)
```

**When to use:** Always in production. This is the most secure approach.

**Setup:**
```bash
# Linux/macOS
export WORKBREW_API_KEY="your-api-key"
export WORKBREW_WORKSPACE="your-workspace-id"

# Windows PowerShell
$env:WORKBREW_API_KEY="your-api-key"
$env:WORKBREW_WORKSPACE="your-workspace-id"

# Docker
docker run -e WORKBREW_API_KEY="your-api-key" -e WORKBREW_WORKSPACE="your-workspace-id" myapp

# Kubernetes Secret
kubectl create secret generic workbrew-credentials --from-literal=api-key="your-api-key" --from-literal=workspace="your-workspace-id"
```

---

### Option 2: Configuration Files

Load API keys from configuration files (not committed to version control):

```go
package main

import (
    "encoding/json"
    "os"
)

type Config struct {
    APIKey    string `json:"api_key"`
    Workspace string `json:"workspace"`
}

func loadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}

func main() {
    // Load from config file
    config, err := loadConfig("config.json")
    if err != nil {
        log.Fatal(err)
    }

    workbrewClient, err := client.NewClient(config.APIKey, config.Workspace)
    // ... use client
}
```

**config.json:**
```json
{
  "api_key": "your-api-key-here",
  "workspace": "your-workspace-id"
}
```

**When to use:** Development environments where you need per-developer configuration.

**.gitignore:**
```
config.json
*.local.json
```

---

### Option 3: Secret Management Services

Use dedicated secret management services for enterprise deployments:

**AWS Secrets Manager:**
```go
import (
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/secretsmanager"
)

func getAPIKeyFromAWS() (string, error) {
    sess := session.Must(session.NewSession())
    svc := secretsmanager.New(sess)

    result, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
        SecretId: aws.String("workbrew/api-key"),
    })
    if err != nil {
        return "", err
    }

    return *result.SecretString, nil
}

func main() {
    apiKey, err := getAPIKeyFromAWS()
    if err != nil {
        log.Fatal(err)
    }

    workspace := os.Getenv("WORKBREW_WORKSPACE")
    workbrewClient, err := client.NewClient(apiKey, workspace)
    // ... use client
}
```

**HashiCorp Vault:**
```go
import "github.com/hashicorp/vault/api"

func getAPIKeyFromVault() (string, error) {
    vaultClient, err := api.NewClient(api.DefaultConfig())
    if err != nil {
        return "", err
    }

    secret, err := vaultClient.Logical().Read("secret/data/workbrew")
    if err != nil {
        return "", err
    }

    apiKey := secret.Data["data"].(map[string]any)["api_key"].(string)
    return apiKey, nil
}
```

**When to use:** Production environments with compliance requirements or centralized secret management.

---

### Option 4: Multiple API Keys

Use different API keys for different purposes or rate limits:

```go
package main

type WorkbrewService struct {
    publicClient  *client.Client
    premiumClient *client.Client
}

func NewWorkbrewService() (*WorkbrewService, error) {
    workspace := os.Getenv("WORKBREW_WORKSPACE")

    // Public API (free tier)
    publicClient, err := client.NewClient(
        os.Getenv("WORKBREW_PUBLIC_API_KEY"),
        workspace,
    )
    if err != nil {
        return nil, err
    }

    // Premium API (higher limits)
    premiumClient, err := client.NewClient(
        os.Getenv("WORKBREW_PREMIUM_API_KEY"),
        workspace,
    )
    if err != nil {
        return nil, err
    }

    return &WorkbrewService{
        publicClient:  publicClient,
        premiumClient: premiumClient,
    }, nil
}

func (s *WorkbrewService) ListDevices(ctx context.Context, usePremium bool) {
    // Choose client based on requirements
    var c *client.Client
    if usePremium {
        c = s.premiumClient
    } else {
        c = s.publicClient
    }

    devicesService := devices.NewService(c)
    // ... use service
}
```

**When to use:** When you have multiple API keys with different rate limits or permissions.

---

## Security Best Practices

### ✅ Do:

- Store API keys in environment variables
- Use secret management services in production
- Rotate API keys regularly (recreate client with new key)
- Use different keys for different environments
- Revoke compromised keys immediately
- Add `*.env` files to `.gitignore`

### ❌ Don't:

- Hardcode API keys in source code
- Commit API keys to version control
- Share API keys in plaintext (email, chat, etc.)
- Use production keys in development
- Log API keys in application logs
- Store API keys in client-side code

### Key Rotation

When you need to rotate an API key, create a new client:

```go
// Old client with compromised key
oldClient, _ := client.NewClient(oldKey, workspace)

// Create new client with rotated key
newClient, _ := client.NewClient(newKey, workspace)

// Switch to using newClient
// The API key is immutable after client creation
```

## Troubleshooting

### Authentication Failed (401 Unauthorized)

**Symptoms:** `WrongCredentialsError` or `AuthenticationRequiredError`

**Solutions:**
```go
import "github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"

// Check for authentication errors
if err != nil && client.IsUnauthorized(err) {
    log.Println("Invalid API key - check your credentials")
}

// Or check specific error codes
if err != nil && client.IsWrongCredentials(err) {
    log.Println("API key is incorrect")
}
```

**Common causes:**
- Invalid or expired API key
- API key not set in environment variable
- Typo in the API key
- Using wrong key for the environment

### Rate Limiting (429 Too Many Requests)

**Symptoms:** `QuotaExceededError` or `TooManyRequestsError`

**Solution:**
```go
if err != nil && client.IsQuotaExceeded(err) {
    log.Println("Rate limit exceeded - wait before retrying")
    time.Sleep(60 * time.Second)
}
```

**Note:** Configure retry logic to handle rate limits - see [Timeouts & Retries](timeouts-retries.md) for more details.

## Testing with Authentication

### Unit Tests

Mock the client to avoid real API calls:

```go
func TestMyFunction(t *testing.T) {
    // Use a mock/test API key and workspace
    testClient, _ := client.NewClient("test-api-key", "test-workspace")

    // Configure mock HTTP responses
    // ... your test code
}
```

### Acceptance Tests

Use a dedicated test API key:

```bash
# Set test API key
export WORKBREW_TEST_API_KEY="your-test-key"

# Run acceptance tests
go test -tags=acceptance ./...
```

## Related Documentation

- [Timeouts & Retries](timeouts-retries.md) - Configure retry logic for auth errors
- [Logging](logging.md) - Log authentication events
- [Debugging](debugging.md) - Debug authentication issues
