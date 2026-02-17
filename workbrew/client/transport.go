package client

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// Transport represents the HTTP transport layer for Workbrew API.
// It provides methods for making HTTP requests to the Workbrew API with built-in
// authentication, retry logic, and request/response logging.
// This is an internal component - users should use workbrew.NewClient() instead.
type Transport struct {
	client        *resty.Client
	logger        *zap.Logger
	authConfig    *AuthConfig
	BaseURL       string
	globalHeaders map[string]string
	userAgent     string
}

// NewTransport creates a new Workbrew API transport with the provided API key and workspace.
// This is an internal function - users should use workbrew.NewClient() instead.
//
// Parameters:
//   - apiKey: Your Workbrew API key (required)
//   - workspaceName: The name of the workspace to use (required)
//   - options: Optional transport configuration options
//
// Returns:
//   - *Transport: Configured API transport instance
//   - error: Any error encountered during transport creation
//
// The transport is configured with:
//   - Default timeout of 120 seconds
//   - Automatic retry on transient failures (up to 3 retries)
//   - Gzip compression support
//   - Bearer token authentication
//   - Production-ready logger (use WithLogger to customize)
func NewTransport(apiKey string, workspaceName string, options ...ClientOption) (*Transport, error) {
	// Create default logger
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Create auth config
	authConfig := &AuthConfig{
		APIKey:     apiKey,
		APIVersion: DefaultAPIVersion,
	}

	// Format: "go-api-sdk-workbrew/1.0.0; gzip"
	// The "gzip" keyword helps with content serving optimization
	userAgent := fmt.Sprintf("%s/%s; gzip", UserAgentBase, Version)

	restyClient := resty.New()
	restyClient.SetTimeout(DefaultTimeout * time.Second)
	restyClient.SetRetryCount(MaxRetries)
	restyClient.SetRetryWaitTime(RetryWaitTime * time.Second)
	restyClient.SetRetryMaxWaitTime(RetryMaxWaitTime * time.Second)
	restyClient.SetHeader("User-Agent", userAgent)
	restyClient.SetHeader("Accept-Encoding", "gzip")

	transport := &Transport{
		client:        restyClient,
		logger:        logger,
		authConfig:    authConfig,
		BaseURL:       DefaultBaseURL,
		globalHeaders: make(map[string]string),
		userAgent:     userAgent,
	}

	// Apply any additional options
	for _, option := range options {
		if err := option(transport); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	if err := SetupAuthentication(restyClient, authConfig, logger); err != nil {
		return nil, fmt.Errorf("failed to setup authentication: %w", err)
	}

	baseURLWithWorkspace := fmt.Sprintf("%s/workspaces/%s", transport.BaseURL, workspaceName)
	restyClient.SetBaseURL(baseURLWithWorkspace)

	logger.Info("Workbrew API transport created",
		zap.String("base_url", baseURLWithWorkspace),
		zap.String("api_version", authConfig.APIVersion))

	return transport, nil
}

// GetHTTPClient returns the underlying resty HTTP client.
// Use this to access advanced resty features or customize the HTTP client directly.
//
// Returns:
//   - *resty.Client: The underlying resty client instance
func (t *Transport) GetHTTPClient() *resty.Client {
	return t.client
}

// GetLogger returns the configured zap logger instance.
// Use this to add custom logging within your application using the same logger.
//
// Returns:
//   - *zap.Logger: The configured logger instance
func (t *Transport) GetLogger() *zap.Logger {
	return t.logger
}

// SetWorkspace changes the active workspace for all subsequent API calls.
// This updates the base URL to target the specified workspace.
//
// Parameters:
//   - workspaceName: The name of the workspace to switch to
//
// Example:
//
//	transport.SetWorkspace("production-workspace")
func (t *Transport) SetWorkspace(workspaceName string) {
	baseURLWithWorkspace := fmt.Sprintf("%s/workspaces/%s", t.BaseURL, workspaceName)
	t.client.SetBaseURL(baseURLWithWorkspace)
	t.logger.Info("Workspace changed", zap.String("workspace", workspaceName))
}

// QueryBuilder creates a new query builder instance for constructing URL query parameters.
// The query builder provides a fluent interface for adding parameters with type safety.
//
// Returns:
//   - interfaces.ServiceQueryBuilder: A new query builder instance
//
// Example:
//
//	params := transport.QueryBuilder().
//	    AddString("name", "test").
//	    AddInt("limit", 100).
//	    AddBool("active", true).
//	    Build()
func (t *Transport) QueryBuilder() interfaces.ServiceQueryBuilder {
	return NewQueryBuilder()
}
