package client

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// Client represents the HTTP client for Workbrew API
type Client struct {
	client        *resty.Client
	logger        *zap.Logger
	authConfig    *AuthConfig
	BaseURL       string
	globalHeaders map[string]string
	userAgent     string
}

// NewClient creates a new Workbrew API client
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

	// Format: "go-api-sdk-workbrew/1.0.0"
	userAgent := fmt.Sprintf("%s/%s", UserAgentBase, Version)

	// Create resty client
	restyClient := resty.New()
	restyClient.SetTimeout(DefaultTimeout * time.Second)
	restyClient.SetRetryCount(MaxRetries)
	restyClient.SetRetryWaitTime(RetryWaitTime * time.Second)
	restyClient.SetRetryMaxWaitTime(RetryMaxWaitTime * time.Second)
	restyClient.SetHeader("User-Agent", userAgent)

	// Create client instance
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

	// Setup authentication
	if err := SetupAuthentication(restyClient, authConfig, logger); err != nil {
		return nil, fmt.Errorf("failed to setup authentication: %w", err)
	}

	// Set base URL with workspace
	baseURLWithWorkspace := fmt.Sprintf("%s/workspaces/%s", client.BaseURL, workspaceName)
	restyClient.SetBaseURL(baseURLWithWorkspace)

	logger.Info("Workbrew API client created",
		zap.String("base_url", baseURLWithWorkspace),
		zap.String("api_version", authConfig.APIVersion))

	return client, nil
}

// GetHTTPClient returns the underlying resty client
func (c *Client) GetHTTPClient() *resty.Client {
	return c.client
}

// GetLogger returns the logger instance
func (c *Client) GetLogger() *zap.Logger {
	return c.logger
}

// SetWorkspace changes the workspace for subsequent API calls
func (c *Client) SetWorkspace(workspaceName string) {
	baseURLWithWorkspace := fmt.Sprintf("%s/workspaces/%s", c.BaseURL, workspaceName)
	c.client.SetBaseURL(baseURLWithWorkspace)
	c.logger.Info("Workspace changed", zap.String("workspace", workspaceName))
}

// QueryBuilder creates a new query builder instance
func (c *Client) QueryBuilder() interfaces.ServiceQueryBuilder {
	return NewQueryBuilder()
}
