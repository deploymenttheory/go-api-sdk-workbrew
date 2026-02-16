package client

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"resty.dev/v3"
)

// AuthConfig holds authentication configuration for the Workbrew API.
// It contains the API key and optional API version for authenticating requests.
type AuthConfig struct {
	// APIKey is the bearer token for authentication
	APIKey string

	// APIVersion is the optional API version (defaults to v0)
	APIVersion string
}

// AuthManager handles thread-safe API key management for the Workbrew API.
// It provides concurrent-safe access to API keys and supports runtime key rotation.
type AuthManager struct {
	authConfig *AuthConfig
	logger     *zap.Logger
	mu         sync.RWMutex
}

// NewAuthManager creates a new auth manager with the provided configuration.
//
// Parameters:
//   - authConfig: Authentication configuration containing API key and version
//   - logger: Logger instance for logging authentication operations
//
// Returns:
//   - *AuthManager: A new thread-safe auth manager instance
func NewAuthManager(authConfig *AuthConfig, logger *zap.Logger) *AuthManager {
	return &AuthManager{
		authConfig: authConfig,
		logger:     logger,
	}
}

// Validate checks if the authentication configuration is valid.
// It ensures that the API key is not empty before allowing API requests.
//
// Returns:
//   - error: ValidationError if API key is missing, nil if valid
func (a *AuthConfig) Validate() error {
	if a.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	return nil
}

// GetAPIKey returns the current API key in a thread-safe manner.
//
// Returns:
//   - string: The current API key
//   - error: Error if API key is not set
func (am *AuthManager) GetAPIKey() (string, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	if am.authConfig.APIKey == "" {
		return "", fmt.Errorf("API key is not set")
	}

	return am.authConfig.APIKey, nil
}

// UpdateAPIKey updates the API key in a thread-safe manner.
// This allows for runtime API key rotation without recreating the client.
//
// Parameters:
//   - newAPIKey: The new API key to use for authentication
//
// Returns:
//   - error: Error if the new API key is empty
//
// Example:
//
//	err := authManager.UpdateAPIKey("new-api-key-12345")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (am *AuthManager) UpdateAPIKey(newAPIKey string) error {
	if newAPIKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	oldKey := am.authConfig.APIKey
	am.authConfig.APIKey = newAPIKey

	am.logger.Info("API key updated successfully",
		zap.Bool("had_previous_key", oldKey != ""))

	return nil
}

// ValidateAPIKey validates that an API key is currently set.
//
// Returns:
//   - error: Error if API key is not set or invalid
func (am *AuthManager) ValidateAPIKey() error {
	am.mu.RLock()
	defer am.mu.RUnlock()

	return am.authConfig.Validate()
}

// GetAPIVersion returns the API version in a thread-safe manner.
//
// Returns:
//   - string: The configured API version or default if not set
func (am *AuthManager) GetAPIVersion() string {
	am.mu.RLock()
	defer am.mu.RUnlock()

	if am.authConfig.APIVersion == "" {
		return DefaultAPIVersion
	}
	return am.authConfig.APIVersion
}

// SetupAuthentication configures the HTTP client with bearer token authentication and middleware validation.
// It sets the Authorization header, API version header, and adds request middleware for token validation.
//
// Parameters:
//   - client: The resty HTTP client to configure
//   - authConfig: Authentication configuration containing API key and version
//   - logger: Logger instance for logging authentication setup
//
// Returns:
//   - *AuthManager: Thread-safe auth manager for runtime key management
//   - error: Any error encountered during authentication setup
//
// The function:
//   - Validates the authentication configuration
//   - Creates a thread-safe AuthManager
//   - Sets Bearer token authentication scheme
//   - Adds the API version header (X-Workbrew-API-Version)
//   - Configures request middleware to validate API key before each request
//   - Logs successful authentication configuration
func SetupAuthentication(client *resty.Client, authConfig *AuthConfig, logger *zap.Logger) (*AuthManager, error) {
	if err := authConfig.Validate(); err != nil {
		logger.Error("Authentication validation failed", zap.Error(err))
		return nil, fmt.Errorf("authentication validation failed: %w", err)
	}

	// Create auth manager for thread-safe key management
	authManager := NewAuthManager(authConfig, logger)

	// Set initial authentication
	client.SetAuthScheme("Bearer")
	client.SetAuthToken(authConfig.APIKey)

	apiVersion := authConfig.APIVersion
	if apiVersion == "" {
		apiVersion = DefaultAPIVersion
	}
	client.SetHeader(APIVersionHeader, apiVersion)

	// Add request middleware to validate API key before each request
	// This ensures the key is always present and allows for runtime key rotation
	client.AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
		apiKey, err := authManager.GetAPIKey()
		if err != nil {
			logger.Error("Failed to get valid API key for request", zap.Error(err))
			return fmt.Errorf("failed to get valid API key: %w", err)
		}
		// Update the request with current API key (supports runtime rotation)
		req.SetAuthScheme("Bearer")
		req.SetAuthToken(apiKey)
		return nil
	})

	logger.Info("Authentication configured successfully with middleware validation",
		zap.String("api_version", apiVersion))

	return authManager, nil
}
