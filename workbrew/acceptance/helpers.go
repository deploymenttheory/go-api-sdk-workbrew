package acceptance

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// SkipIfAPIKeyNotSet skips the test if the API key is not configured
func SkipIfAPIKeyNotSet(t *testing.T) {
	t.Helper()
	if !IsAPIKeySet() {
		t.Skip("WORKBREW_API_KEY not set, skipping acceptance test")
	}
}

// SkipIfWorkspaceNotSet skips the test if the workspace name is not configured
func SkipIfWorkspaceNotSet(t *testing.T) {
	t.Helper()
	if !IsWorkspaceSet() {
		t.Skip("WORKBREW_WORKSPACE_NAME not set, skipping acceptance test")
	}
}

// SkipIfNotConfigured skips the test if either API key or workspace is not configured
func SkipIfNotConfigured(t *testing.T) {
	t.Helper()
	if !IsConfigured() {
		t.Skip("WORKBREW_API_KEY or WORKBREW_WORKSPACE_NAME not set, skipping acceptance test")
	}
}

// RequireClient ensures the shared client is initialized
// Skips the test if the API key or workspace is not set or client initialization fails
func RequireClient(t *testing.T) {
	t.Helper()
	SkipIfNotConfigured(t)

	if Client == nil {
		err := InitClient()
		require.NoError(t, err, "Failed to initialize Workbrew client")
	}
}

// RateLimitedTest wraps a test function with rate limiting
// Automatically sleeps after test execution to respect API rate limits
func RateLimitedTest(t *testing.T, testFunc func(t *testing.T)) {
	t.Helper()
	defer func() {
		LogTestWarning(t, "Rate limiting: sleeping for %v", Config.RateLimitDelay)
		time.Sleep(Config.RateLimitDelay)
	}()
	testFunc(t)
}

// NewContext creates a context with timeout for acceptance tests
func NewContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), Config.RequestTimeout)
}

// LogResponse logs the response details if verbose mode is enabled
func LogResponse(t *testing.T, message string, details ...any) {
	t.Helper()
	if Config.Verbose {
		if len(details) > 0 {
			t.Logf(message, details...)
		} else {
			t.Log(message)
		}
	}
}

// AssertNoError is a helper that fails the test if an error occurs
// and logs additional context in verbose mode
func AssertNoError(t *testing.T, err error, msgAndArgs ...any) {
	t.Helper()
	if err != nil {
		if Config.Verbose {
			t.Logf("Error occurred: %v", err)
		}
	}
	require.NoError(t, err, msgAndArgs...)
}

// AssertNotNil is a helper that fails the test if the object is nil
func AssertNotNil(t *testing.T, object any, msgAndArgs ...any) {
	t.Helper()
	require.NotNil(t, object, msgAndArgs...)
}

// Cleanup registers a cleanup function that respects the SkipCleanup flag
func Cleanup(t *testing.T, cleanupFunc func()) {
	t.Helper()
	if !Config.SkipCleanup {
		t.Cleanup(cleanupFunc)
	} else if Config.Verbose {
		t.Log("Skipping cleanup due to WORKBREW_SKIP_CLEANUP=true")
	}
}

// isGitHubActions returns true if running in GitHub Actions
func isGitHubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

// LogTestStage logs a test stage with optional GitHub Actions notice annotation
func LogTestStage(t *testing.T, stage string, message string, details ...any) {
	t.Helper()

	formattedMsg := message
	if len(details) > 0 {
		formattedMsg = fmt.Sprintf(message, details...)
	}

	if isGitHubActions() {
		fmt.Printf("::notice title=%s::%s\n", stage, formattedMsg)
	}

	if Config.Verbose {
		t.Logf("🎯 [%s] %s", stage, formattedMsg)
	}
}

// LogTestSuccess logs a successful test operation
func LogTestSuccess(t *testing.T, message string, details ...any) {
	t.Helper()

	formattedMsg := message
	if len(details) > 0 {
		formattedMsg = fmt.Sprintf(message, details...)
	}

	if isGitHubActions() {
		fmt.Printf("::notice title=✅ Success::%s\n", formattedMsg)
	}

	if Config.Verbose {
		t.Logf("✅ %s", formattedMsg)
	}
}

// LogTestWarning logs a test warning
func LogTestWarning(t *testing.T, message string, details ...any) {
	t.Helper()

	formattedMsg := message
	if len(details) > 0 {
		formattedMsg = fmt.Sprintf(message, details...)
	}

	if isGitHubActions() {
		fmt.Printf("::warning title=⚠️  Warning::%s\n", formattedMsg)
	}

	if Config.Verbose {
		t.Logf("⚠️  %s", formattedMsg)
	}
}

// LogTestError logs a test error (does not fail test)
func LogTestError(t *testing.T, message string, details ...any) {
	t.Helper()

	formattedMsg := message
	if len(details) > 0 {
		formattedMsg = fmt.Sprintf(message, details...)
	}

	if isGitHubActions() {
		fmt.Printf("::error title=❌ Error::%s\n", formattedMsg)
	}

	if Config.Verbose {
		t.Logf("❌ %s", formattedMsg)
	}
}

// LogGroup starts a collapsible log group in GitHub Actions
func LogGroup(title string) {
	if isGitHubActions() {
		fmt.Printf("::group::%s\n", title)
	} else {
		fmt.Printf("\n=== %s ===\n", title)
	}
}

// LogGroupEnd ends a collapsible log group in GitHub Actions
func LogGroupEnd() {
	if isGitHubActions() {
		fmt.Println("::endgroup::")
	}
}
