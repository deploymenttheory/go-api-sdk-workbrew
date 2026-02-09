package client

import (
	"fmt"

	"go.uber.org/zap"
	"resty.dev/v3"
)

// AuthConfig holds authentication configuration for the Workbrew API
type AuthConfig struct {
	// APIKey is the bearer token for authentication
	APIKey string

	// APIVersion is the optional API version (defaults to v0)
	APIVersion string
}

// Validate checks if the auth configuration is valid
func (a *AuthConfig) Validate() error {
	if a.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	return nil
}

// SetupAuthentication configures the resty client with bearer token authentication
func SetupAuthentication(client *resty.Client, authConfig *AuthConfig, logger *zap.Logger) error {
	if err := authConfig.Validate(); err != nil {
		logger.Error("Authentication validation failed", zap.Error(err))
		return fmt.Errorf("authentication validation failed: %w", err)
	}

	// Set the bearer token
	client.SetAuthScheme("Bearer")
	client.SetAuthToken(authConfig.APIKey)

	// Set API version header if specified
	apiVersion := authConfig.APIVersion
	if apiVersion == "" {
		apiVersion = DefaultAPIVersion
	}
	client.SetHeader(APIVersionHeader, apiVersion)

	logger.Info("Authentication configured successfully",
		zap.String("api_version", apiVersion))

	return nil
}
