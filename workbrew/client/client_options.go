package client

import (
	"time"

	"go.uber.org/zap"
)

// ClientOption is a function type for configuring the Client
type ClientOption func(*Client) error

// WithBaseURL sets a custom base URL for the API client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		c.BaseURL = baseURL
		c.logger.Info("Base URL configured", zap.String("base_url", baseURL))
		return nil
	}
}

// WithAPIVersion sets a custom API version
func WithAPIVersion(version string) ClientOption {
	return func(c *Client) error {
		c.authConfig.APIVersion = version
		c.logger.Info("API version configured", zap.String("api_version", version))
		return nil
	}
}

// WithTimeout sets a custom timeout for HTTP requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.client.SetTimeout(timeout)
		c.logger.Info("HTTP timeout configured", zap.Duration("timeout", timeout))
		return nil
	}
}

// WithRetryCount sets the number of retries for failed requests
func WithRetryCount(count int) ClientOption {
	return func(c *Client) error {
		c.client.SetRetryCount(count)
		c.logger.Info("Retry count configured", zap.Int("retry_count", count))
		return nil
	}
}

// WithLogger sets a custom logger for the client
func WithLogger(logger *zap.Logger) ClientOption {
	return func(c *Client) error {
		c.logger = logger
		c.logger.Info("Custom logger configured")
		return nil
	}
}

// WithDebug enables debug mode which logs request and response details
func WithDebug() ClientOption {
	return func(c *Client) error {
		c.client.SetDebug(true)
		c.logger.Info("Debug mode enabled")
		return nil
	}
}

// WithUserAgent sets a custom user agent string
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.client.SetHeader("User-Agent", userAgent)
		c.logger.Info("User agent configured", zap.String("user_agent", userAgent))
		return nil
	}
}
