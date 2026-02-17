package client

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// Client represents the HTTP client for Workbrew API.
// It provides methods for making HTTP requests to the Workbrew API with built-in
// authentication, retry logic, and request/response logging.
type Client struct {
	client        *resty.Client
	logger        *zap.Logger
	authConfig    *AuthConfig
	BaseURL       string
	globalHeaders map[string]string
	userAgent     string
}

// NewClient creates a new Workbrew API client with the provided API key and workspace.
//
// Parameters:
//   - apiKey: Your Workbrew API key (required)
//   - workspaceName: The name of the workspace to use (required)
//   - options: Optional client configuration options
//
// Returns:
//   - *Client: Configured API client instance
//   - error: Any error encountered during client creation
//
// The client is configured with:
//   - Default timeout of 120 seconds
//   - Automatic retry on transient failures (up to 3 retries)
//   - Gzip compression support
//   - Bearer token authentication
//   - Production-ready logger (use WithLogger to customize)
//
// Example:
//
//	client, err := client.NewClient("your-api-key", "your-workspace")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Example with options:
//
//	client, err := client.NewClient(
//	    "your-api-key",
//	    "your-workspace",
//	    client.WithTimeout(60 * time.Second),
//	    client.WithRetryCount(5),
//	    client.WithDebug(),
//	)
func NewClient(apiKey string, workspaceName string, options ...ClientOption) (*Client, error) {
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

	client := &Client{
		client:        restyClient,
		logger:        logger,
		authConfig:    authConfig,
		BaseURL:       DefaultBaseURL,
		globalHeaders: make(map[string]string),
		userAgent:     userAgent,
	}

	// Apply any additional options
	for _, option := range options {
		if err := option(client); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	if err := SetupAuthentication(restyClient, authConfig, logger); err != nil {
		return nil, fmt.Errorf("failed to setup authentication: %w", err)
	}

	baseURLWithWorkspace := fmt.Sprintf("%s/workspaces/%s", client.BaseURL, workspaceName)
	restyClient.SetBaseURL(baseURLWithWorkspace)

	logger.Info("Workbrew API client created",
		zap.String("base_url", baseURLWithWorkspace),
		zap.String("api_version", authConfig.APIVersion))

	return client, nil
}

// GetHTTPClient returns the underlying resty HTTP client.
// Use this to access advanced resty features or customize the HTTP client directly.
//
// Returns:
//   - *resty.Client: The underlying resty client instance
func (c *Client) GetHTTPClient() *resty.Client {
	return c.client
}

// GetLogger returns the configured zap logger instance.
// Use this to add custom logging within your application using the same logger.
//
// Returns:
//   - *zap.Logger: The configured logger instance
func (c *Client) GetLogger() *zap.Logger {
	return c.logger
}

// SetWorkspace changes the active workspace for all subsequent API calls.
// This updates the base URL to target the specified workspace.
//
// Parameters:
//   - workspaceName: The name of the workspace to switch to
//
// Example:
//
//	client.SetWorkspace("production-workspace")
func (c *Client) SetWorkspace(workspaceName string) {
	baseURLWithWorkspace := fmt.Sprintf("%s/workspaces/%s", c.BaseURL, workspaceName)
	c.client.SetBaseURL(baseURLWithWorkspace)
	c.logger.Info("Workspace changed", zap.String("workspace", workspaceName))
}

// QueryBuilder creates a new query builder instance for constructing URL query parameters.
// The query builder provides a fluent interface for adding parameters with type safety.
//
// Returns:
//   - interfaces.ServiceQueryBuilder: A new query builder instance
//
// Example:
//
//	params := client.QueryBuilder().
//	    AddString("name", "test").
//	    AddInt("limit", 100).
//	    AddBool("active", true).
//	    Build()
func (c *Client) QueryBuilder() interfaces.ServiceQueryBuilder {
	return NewQueryBuilder()
}
