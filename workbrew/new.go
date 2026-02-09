package workbrew

import (
	"fmt"
	"os"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/analytics"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewcommands"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewconfigurations"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewfiles"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewtaps"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/casks"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devicegroups"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/events"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/formulae"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/licenses"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/vulnerabilities"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/vulnerabilitychanges"
)

// Client is the main entry point for the Workbrew API SDK
// It aggregates all service clients and provides a unified interface
type Client struct {
	*client.Client

	// Services
	Analytics            *analytics.Service
	BrewCommands         *brewcommands.Service
	BrewConfigurations   *brewconfigurations.Service
	Brewfiles            *brewfiles.Service
	BrewTaps             *brewtaps.Service
	Casks                *casks.Service
	DeviceGroups         *devicegroups.Service
	Devices              *devices.Service
	Events               *events.Service
	Formulae             *formulae.Service
	Licenses             *licenses.Service
	Vulnerabilities      *vulnerabilities.Service
	VulnerabilityChanges *vulnerabilitychanges.Service
}

// NewClient creates a new Workbrew API client
//
// Parameters:
//   - apiKey: The bearer token for authentication
//   - workspaceName: The workspace slug to operate on
//   - options: Optional client configuration options
//
// Example:
//
//	client, err := workbrew.NewClient(
//	    "your-api-key",
//	    "your-workspace",
//	    workbrew.WithDebug(),
//	)
func NewClient(apiKey string, workspaceName string, options ...client.ClientOption) (*Client, error) {
	// Create base HTTP client
	httpClient, err := client.NewClient(apiKey, workspaceName, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	// Initialize service clients
	c := &Client{
		Client:               httpClient,
		Analytics:            analytics.NewService(httpClient),
		BrewCommands:         brewcommands.NewService(httpClient),
		BrewConfigurations:   brewconfigurations.NewService(httpClient),
		Brewfiles:            brewfiles.NewService(httpClient),
		BrewTaps:             brewtaps.NewService(httpClient),
		Casks:                casks.NewService(httpClient),
		DeviceGroups:         devicegroups.NewService(httpClient),
		Devices:              devices.NewService(httpClient),
		Events:               events.NewService(httpClient),
		Formulae:             formulae.NewService(httpClient),
		Licenses:             licenses.NewService(httpClient),
		Vulnerabilities:      vulnerabilities.NewService(httpClient),
		VulnerabilityChanges: vulnerabilitychanges.NewService(httpClient),
	}

	return c, nil
}

// NewClientFromEnv creates a new client using environment variables
//
// Required environment variables:
//   - WORKBREW_API_KEY: The bearer token for authentication
//   - WORKBREW_WORKSPACE: The workspace slug
//
// Optional environment variables:
//   - WORKBREW_BASE_URL: Custom base URL (defaults to https://console.workbrew.com)
//   - WORKBREW_API_VERSION: API version (defaults to v0)
//
// Example:
//
//	client, err := workbrew.NewClientFromEnv()
func NewClientFromEnv(options ...client.ClientOption) (*Client, error) {
	apiKey := os.Getenv("WORKBREW_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("WORKBREW_API_KEY environment variable is required")
	}

	workspaceName := os.Getenv("WORKBREW_WORKSPACE")
	if workspaceName == "" {
		return nil, fmt.Errorf("WORKBREW_WORKSPACE environment variable is required")
	}

	// Check for optional environment variables and append to options
	if baseURL := os.Getenv("WORKBREW_BASE_URL"); baseURL != "" {
		options = append(options, client.WithBaseURL(baseURL))
	}

	if apiVersion := os.Getenv("WORKBREW_API_VERSION"); apiVersion != "" {
		options = append(options, client.WithAPIVersion(apiVersion))
	}

	return NewClient(apiKey, workspaceName, options...)
}
