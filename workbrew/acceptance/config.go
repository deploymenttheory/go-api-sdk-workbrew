package acceptance

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
)

// TestConfig holds configuration for acceptance tests
type TestConfig struct {
	APIKey                string
	WorkspaceName         string
	BaseURL               string
	RateLimitDelay        time.Duration
	RequestTimeout        time.Duration
	SkipCleanup           bool
	Verbose               bool
	KnownDeviceSerial     string
	KnownGroupName        string
	KnownBrewfileName     string
	KnownCaskName         string
	KnownFormulaName      string
	KnownEventID          string
	KnownVulnerabilityCVE string
}

var (
	// Config is the global test configuration
	Config *TestConfig
	// Client is the shared Workbrew client for acceptance tests
	Client *client.Transport
)

// init initializes the test configuration from environment variables
func init() {
	Config = &TestConfig{
		APIKey:                getEnv("WORKBREW_API_KEY", ""),
		WorkspaceName:         getEnv("WORKBREW_WORKSPACE_NAME", ""),
		BaseURL:               getEnv("WORKBREW_BASE_URL", "https://console.workbrew.com"),
		RateLimitDelay:        getDurationEnv("WORKBREW_RATE_LIMIT_DELAY", 2*time.Second), // Conservative default
		RequestTimeout:        getDurationEnv("WORKBREW_REQUEST_TIMEOUT", 30*time.Second),
		SkipCleanup:           getBoolEnv("WORKBREW_SKIP_CLEANUP", false),
		Verbose:               getBoolEnv("WORKBREW_VERBOSE", false),
		KnownDeviceSerial:     getEnv("WORKBREW_TEST_DEVICE_SERIAL", ""),
		KnownGroupName:        getEnv("WORKBREW_TEST_GROUP_NAME", ""),
		KnownBrewfileName:     getEnv("WORKBREW_TEST_BREWFILE_NAME", ""),
		KnownCaskName:         getEnv("WORKBREW_TEST_CASK_NAME", "docker"), // Common cask
		KnownFormulaName:      getEnv("WORKBREW_TEST_FORMULA_NAME", "git"), // Common formula
		KnownEventID:          getEnv("WORKBREW_TEST_EVENT_ID", ""),
		KnownVulnerabilityCVE: getEnv("WORKBREW_TEST_CVE", ""),
	}
}

// InitClient initializes the shared Workbrew client
// Returns an error if the API key or workspace name is not set or client creation fails
func InitClient() error {
	if Config.APIKey == "" {
		return fmt.Errorf("WORKBREW_API_KEY environment variable is not set")
	}
	if Config.WorkspaceName == "" {
		return fmt.Errorf("WORKBREW_WORKSPACE_NAME environment variable is not set")
	}

	var err error
	Client, err = client.NewTransport(
		Config.APIKey,
		Config.WorkspaceName,
		client.WithBaseURL(Config.BaseURL),
		client.WithTimeout(Config.RequestTimeout),
	)
	if err != nil {
		return fmt.Errorf("failed to create Workbrew client: %w", err)
	}

	if Config.Verbose {
		log.Printf("Acceptance test client initialized with base URL: %s, workspace: %s", Config.BaseURL, Config.WorkspaceName)
	}

	return nil
}

// IsAPIKeySet returns true if the API key is configured
func IsAPIKeySet() bool {
	return Config.APIKey != ""
}

// IsWorkspaceSet returns true if the workspace name is configured
func IsWorkspaceSet() bool {
	return Config.WorkspaceName != ""
}

// IsConfigured returns true if both API key and workspace are configured
func IsConfigured() bool {
	return IsAPIKeySet() && IsWorkspaceSet()
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getBoolEnv retrieves a boolean environment variable or returns a default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("Warning: invalid boolean value for %s: %s, using default: %v", key, value, defaultValue)
			return defaultValue
		}
		return parsed
	}
	return defaultValue
}

// getDurationEnv retrieves a duration environment variable or returns a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		parsed, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("Warning: invalid duration value for %s: %s, using default: %v", key, value, defaultValue)
			return defaultValue
		}
		return parsed
	}
	return defaultValue
}
