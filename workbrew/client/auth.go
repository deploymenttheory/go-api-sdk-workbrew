package client

import (
	"fmt"

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

// SetupAuthentication configures the HTTP client with bearer token authentication.
// It sets the Authorization header and API version header for all requests.
//
// Parameters:
//   - client: The resty HTTP client to configure
//   - authConfig: Authentication configuration containing API key and version
//   - logger: Logger instance for logging authentication setup
//
// Returns:
//   - error: Any error encountered during authentication setup
//
// The function:
//   - Validates the authentication configuration
//   - Sets Bearer token authentication scheme
//   - Adds the API version header (X-Workbrew-API-Version)
//   - Logs successful authentication configuration
func SetupAuthentication(client *resty.Client, authConfig *AuthConfig, logger *zap.Logger) error {
	if err := authConfig.Validate(); err != nil {
		logger.Error("Authentication validation failed", zap.Error(err))
		return fmt.Errorf("authentication validation failed: %w", err)
	}

	client.SetAuthScheme("Bearer")
	client.SetAuthToken(authConfig.APIKey)

	apiVersion := authConfig.APIVersion
	if apiVersion == "" {
		apiVersion = DefaultAPIVersion
	}
	client.SetHeader(APIVersionHeader, apiVersion)

	logger.Info("Authentication configured successfully",
		zap.String("api_version", apiVersion))

	return nil
}
