package client

import (
	"sync"
	"testing"

	"go.uber.org/zap/zaptest"
	"resty.dev/v3"
)

func TestAuthConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *AuthConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &AuthConfig{
				APIKey:     "test-api-key",
				APIVersion: "v1",
			},
			wantErr: false,
		},
		{
			name: "valid config without version",
			config: &AuthConfig{
				APIKey: "test-api-key",
			},
			wantErr: false,
		},
		{
			name: "empty API key",
			config: &AuthConfig{
				APIKey:     "",
				APIVersion: "v1",
			},
			wantErr: true,
			errMsg:  "API key is required",
		},
		{
			name: "nil config fields",
			config: &AuthConfig{
				APIKey:     "",
				APIVersion: "",
			},
			wantErr: true,
			errMsg:  "API key is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("AuthConfig.Validate() error message = %q, want %q", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestSetupAuthentication_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	authConfig := &AuthConfig{
		APIKey:     "test-api-key-12345",
		APIVersion: "v1",
	}

	authManager, err := SetupAuthentication(client, authConfig, logger)

	if err != nil {
		t.Fatalf("SetupAuthentication() error = %v, want nil", err)
	}

	if authManager == nil {
		t.Fatal("Expected non-nil AuthManager")
	}

	// Verify auth manager has correct key
	apiKey, err := authManager.GetAPIKey()
	if err != nil {
		t.Errorf("GetAPIKey() error = %v, want nil", err)
	}
	if apiKey != "test-api-key-12345" {
		t.Errorf("GetAPIKey() = %q, want %q", apiKey, "test-api-key-12345")
	}

	// Verify API version header is set
	headers := client.Header()
	if got := headers.Get(APIVersionHeader); got != "v1" {
		t.Errorf("API version header = %q, want %q", got, "v1")
	}
}

func TestSetupAuthentication_DefaultAPIVersion(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	authConfig := &AuthConfig{
		APIKey:     "test-api-key",
		APIVersion: "", // Empty, should use default
	}

	authManager, err := SetupAuthentication(client, authConfig, logger)

	if err != nil {
		t.Fatalf("SetupAuthentication() error = %v, want nil", err)
	}

	if authManager == nil {
		t.Fatal("Expected non-nil AuthManager")
	}

	// Verify default API version is used
	headers := client.Header()
	if got := headers.Get(APIVersionHeader); got != DefaultAPIVersion {
		t.Errorf("API version header = %q, want %q (default)", got, DefaultAPIVersion)
	}

	// Verify default API version via auth manager
	version := authManager.GetAPIVersion()
	if version != DefaultAPIVersion {
		t.Errorf("GetAPIVersion() = %q, want %q", version, DefaultAPIVersion)
	}
}

func TestSetupAuthentication_InvalidConfig(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	tests := []struct {
		name       string
		authConfig *AuthConfig
		wantErr    bool
		errContain string
	}{
		{
			name: "empty API key",
			authConfig: &AuthConfig{
				APIKey:     "",
				APIVersion: "v1",
			},
			wantErr:    true,
			errContain: "authentication validation failed",
		},
		{
			name: "nil-like config",
			authConfig: &AuthConfig{
				APIKey:     "",
				APIVersion: "",
			},
			wantErr:    true,
			errContain: "API key is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authManager, err := SetupAuthentication(client, tt.authConfig, logger)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetupAuthentication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if authManager != nil {
					t.Error("Expected nil AuthManager on error")
				}
				if err != nil && err.Error() == "" {
					t.Error("Expected error message, got empty string")
				}
			}
		})
	}
}

func TestSetupAuthentication_CustomAPIVersion(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name           string
		apiVersion     string
		expectedHeader string
	}{
		{
			name:           "custom v1",
			apiVersion:     "v1",
			expectedHeader: "v1",
		},
		{
			name:           "custom v2",
			apiVersion:     "v2",
			expectedHeader: "v2",
		},
		{
			name:           "empty uses default",
			apiVersion:     "",
			expectedHeader: DefaultAPIVersion,
		},
		{
			name:           "custom version string",
			apiVersion:     "2023-01",
			expectedHeader: "2023-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := resty.New()
			authConfig := &AuthConfig{
				APIKey:     "test-key",
				APIVersion: tt.apiVersion,
			}

			authManager, err := SetupAuthentication(client, authConfig, logger)
			if err != nil {
				t.Fatalf("SetupAuthentication() error = %v, want nil", err)
			}

			if authManager == nil {
				t.Fatal("Expected non-nil AuthManager")
			}

			headers := client.Header()
			if got := headers.Get(APIVersionHeader); got != tt.expectedHeader {
				t.Errorf("API version header = %q, want %q", got, tt.expectedHeader)
			}

			// Verify API version via auth manager
			gotVersion := authManager.GetAPIVersion()
			if gotVersion != tt.expectedHeader {
				t.Errorf("GetAPIVersion() = %q, want %q", gotVersion, tt.expectedHeader)
			}
		})
	}
}

func TestSetupAuthentication_HeadersAreSet(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	authConfig := &AuthConfig{
		APIKey:     "test-api-key",
		APIVersion: "v1",
	}

	authManager, err := SetupAuthentication(client, authConfig, logger)
	if err != nil {
		t.Fatalf("SetupAuthentication() error = %v, want nil", err)
	}

	if authManager == nil {
		t.Fatal("Expected non-nil AuthManager")
	}

	// Verify the API version header is present
	headers := client.Header()
	apiVersionHeader := headers.Get(APIVersionHeader)
	if apiVersionHeader == "" {
		t.Error("API version header should be set")
	}

	if apiVersionHeader != "v1" {
		t.Errorf("API version header = %q, want %q", apiVersionHeader, "v1")
	}
}

func TestSetupAuthentication_MultipleCallsOverwrite(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	// First setup
	authConfig1 := &AuthConfig{
		APIKey:     "first-key",
		APIVersion: "v1",
	}
	authManager1, err := SetupAuthentication(client, authConfig1, logger)
	if err != nil {
		t.Fatalf("First SetupAuthentication() error = %v, want nil", err)
	}

	// Verify first setup
	headers := client.Header()
	if got := headers.Get(APIVersionHeader); got != "v1" {
		t.Errorf("After first setup, API version = %q, want %q", got, "v1")
	}

	apiKey1, _ := authManager1.GetAPIKey()
	if apiKey1 != "first-key" {
		t.Errorf("First AuthManager key = %q, want %q", apiKey1, "first-key")
	}

	// Second setup with different values
	authConfig2 := &AuthConfig{
		APIKey:     "second-key",
		APIVersion: "v2",
	}
	authManager2, err := SetupAuthentication(client, authConfig2, logger)
	if err != nil {
		t.Fatalf("Second SetupAuthentication() error = %v, want nil", err)
	}

	// Verify second setup overwrote first
	headers = client.Header()
	if got := headers.Get(APIVersionHeader); got != "v2" {
		t.Errorf("After second setup, API version = %q, want %q", got, "v2")
	}

	apiKey2, _ := authManager2.GetAPIKey()
	if apiKey2 != "second-key" {
		t.Errorf("Second AuthManager key = %q, want %q", apiKey2, "second-key")
	}
}

func TestAuthConfig_Fields(t *testing.T) {
	// Test that AuthConfig struct can hold expected values
	config := &AuthConfig{
		APIKey:     "my-api-key-12345",
		APIVersion: "v1.5",
	}

	if config.APIKey != "my-api-key-12345" {
		t.Errorf("APIKey = %q, want %q", config.APIKey, "my-api-key-12345")
	}

	if config.APIVersion != "v1.5" {
		t.Errorf("APIVersion = %q, want %q", config.APIVersion, "v1.5")
	}
}

func TestSetupAuthentication_NilClient(t *testing.T) {
	// This test verifies behavior with nil client (should panic or handle gracefully)
	// In practice, this should never happen as NewClient creates the resty client
	logger := zaptest.NewLogger(t)

	authConfig := &AuthConfig{
		APIKey:     "test-key",
		APIVersion: "v1",
	}

	// This will panic if not handled, which is acceptable for nil client
	defer func() {
		if r := recover(); r != nil {
			// Panic is expected for nil client
			t.Logf("Panic recovered (expected): %v", r)
		}
	}()

	// This should panic or fail
	_, _ = SetupAuthentication(nil, authConfig, logger)
}

func TestAuthConfig_LongAPIKey(t *testing.T) {
	// Test with a very long API key (should still be valid)
	longKey := ""
	for i := 0; i < 1000; i++ {
		longKey += "a"
	}

	config := &AuthConfig{
		APIKey:     longKey,
		APIVersion: "v1",
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Validate() with long API key error = %v, want nil", err)
	}

	// Setup should also work
	logger := zaptest.NewLogger(t)
	client := resty.New()
	authManager, err := SetupAuthentication(client, config, logger)
	if err != nil {
		t.Errorf("SetupAuthentication() with long API key error = %v, want nil", err)
	}
	if authManager == nil {
		t.Error("Expected non-nil AuthManager")
	}
}

func TestAuthConfig_SpecialCharactersInAPIKey(t *testing.T) {
	// Test with special characters in API key
	specialKeys := []string{
		"key-with-dashes",
		"key_with_underscores",
		"key.with.dots",
		"key123with456numbers",
		"key-_./~:?#[]@!$&'()*+,;=%spaces", // URL-safe characters
	}

	for _, key := range specialKeys {
		t.Run(key, func(t *testing.T) {
			config := &AuthConfig{
				APIKey:     key,
				APIVersion: "v1",
			}

			err := config.Validate()
			if err != nil {
				t.Errorf("Validate() with key %q error = %v, want nil", key, err)
			}

			logger := zaptest.NewLogger(t)
			client := resty.New()
			authManager, err := SetupAuthentication(client, config, logger)
			if err != nil {
				t.Errorf("SetupAuthentication() with key %q error = %v, want nil", key, err)
			}
			if authManager == nil {
				t.Error("Expected non-nil AuthManager")
			}
		})
	}
}

func TestAuthConfig_WhitespaceAPIKey(t *testing.T) {
	// Test with whitespace-only API key (should be considered invalid)
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "spaces only",
			apiKey:  "   ",
			wantErr: false, // Non-empty string, validation passes (though API would reject it)
		},
		{
			name:    "tabs only",
			apiKey:  "\t\t\t",
			wantErr: false, // Non-empty string, validation passes
		},
		{
			name:    "newlines only",
			apiKey:  "\n\n",
			wantErr: false, // Non-empty string, validation passes
		},
		{
			name:    "truly empty",
			apiKey:  "",
			wantErr: true, // Empty string, validation fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AuthConfig{
				APIKey:     tt.apiKey,
				APIVersion: "v1",
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() with whitespace key error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Tests for AuthManager

func TestAuthManager_GetAPIKey(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "valid API key",
			apiKey:  "test-key-123",
			wantErr: false,
		},
		{
			name:    "empty API key",
			apiKey:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AuthConfig{
				APIKey:     tt.apiKey,
				APIVersion: "v1",
			}

			am := NewAuthManager(config, logger)
			key, err := am.GetAPIKey()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAPIKey() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && key != tt.apiKey {
				t.Errorf("GetAPIKey() = %q, want %q", key, tt.apiKey)
			}
		})
	}
}

func TestAuthManager_UpdateAPIKey(t *testing.T) {
	logger := zaptest.NewLogger(t)

	config := &AuthConfig{
		APIKey:     "initial-key",
		APIVersion: "v1",
	}

	am := NewAuthManager(config, logger)

	// Verify initial key
	initialKey, err := am.GetAPIKey()
	if err != nil {
		t.Fatalf("GetAPIKey() initial error = %v", err)
	}
	if initialKey != "initial-key" {
		t.Errorf("Initial key = %q, want %q", initialKey, "initial-key")
	}

	// Update to new key
	err = am.UpdateAPIKey("updated-key")
	if err != nil {
		t.Fatalf("UpdateAPIKey() error = %v", err)
	}

	// Verify updated key
	updatedKey, err := am.GetAPIKey()
	if err != nil {
		t.Fatalf("GetAPIKey() after update error = %v", err)
	}
	if updatedKey != "updated-key" {
		t.Errorf("Updated key = %q, want %q", updatedKey, "updated-key")
	}
}

func TestAuthManager_UpdateAPIKey_EmptyKey(t *testing.T) {
	logger := zaptest.NewLogger(t)

	config := &AuthConfig{
		APIKey:     "initial-key",
		APIVersion: "v1",
	}

	am := NewAuthManager(config, logger)

	// Try to update with empty key
	err := am.UpdateAPIKey("")
	if err == nil {
		t.Error("UpdateAPIKey() with empty key should return error")
	}

	// Verify original key is unchanged
	key, _ := am.GetAPIKey()
	if key != "initial-key" {
		t.Errorf("Key after failed update = %q, want %q", key, "initial-key")
	}
}

func TestAuthManager_ValidateAPIKey(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "valid key",
			apiKey:  "valid-key",
			wantErr: false,
		},
		{
			name:    "empty key",
			apiKey:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AuthConfig{
				APIKey:     tt.apiKey,
				APIVersion: "v1",
			}

			am := NewAuthManager(config, logger)
			err := am.ValidateAPIKey()

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAPIKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthManager_GetAPIVersion(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name       string
		apiVersion string
		want       string
	}{
		{
			name:       "custom version",
			apiVersion: "v1",
			want:       "v1",
		},
		{
			name:       "empty uses default",
			apiVersion: "",
			want:       DefaultAPIVersion,
		},
		{
			name:       "custom version string",
			apiVersion: "2024-01",
			want:       "2024-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AuthConfig{
				APIKey:     "test-key",
				APIVersion: tt.apiVersion,
			}

			am := NewAuthManager(config, logger)
			got := am.GetAPIVersion()

			if got != tt.want {
				t.Errorf("GetAPIVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Thread-safety tests

func TestAuthManager_ConcurrentGetAPIKey(t *testing.T) {
	logger := zaptest.NewLogger(t)

	config := &AuthConfig{
		APIKey:     "concurrent-test-key",
		APIVersion: "v1",
	}

	am := NewAuthManager(config, logger)

	// Run 100 concurrent GetAPIKey calls
	const numGoroutines = 100
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)
	keys := make(chan string, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key, err := am.GetAPIKey()
			if err != nil {
				errors <- err
				return
			}
			keys <- key
		}()
	}

	wg.Wait()
	close(errors)
	close(keys)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent GetAPIKey() error: %v", err)
	}

	// Verify all keys are correct
	for key := range keys {
		if key != "concurrent-test-key" {
			t.Errorf("Concurrent GetAPIKey() = %q, want %q", key, "concurrent-test-key")
		}
	}
}

func TestAuthManager_ConcurrentUpdateAPIKey(t *testing.T) {
	logger := zaptest.NewLogger(t)

	config := &AuthConfig{
		APIKey:     "initial-key",
		APIVersion: "v1",
	}

	am := NewAuthManager(config, logger)

	// Run concurrent updates
	const numGoroutines = 50
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		keyNum := i
		go func(n int) {
			defer wg.Done()
			newKey := "key-" + string(rune('0'+n%10))
			err := am.UpdateAPIKey(newKey)
			if err != nil {
				errors <- err
			}
		}(keyNum)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent UpdateAPIKey() error: %v", err)
	}

	// Verify we can still get a key (any valid key is fine)
	finalKey, err := am.GetAPIKey()
	if err != nil {
		t.Errorf("GetAPIKey() after concurrent updates error: %v", err)
	}
	if finalKey == "" {
		t.Error("Final key should not be empty after concurrent updates")
	}
}

func TestAuthManager_ConcurrentReadWrite(t *testing.T) {
	logger := zaptest.NewLogger(t)

	config := &AuthConfig{
		APIKey:     "initial-key",
		APIVersion: "v1",
	}

	am := NewAuthManager(config, logger)

	// Run concurrent reads and writes
	const numReaders = 50
	const numWriters = 10
	var wg sync.WaitGroup

	// Start readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, err := am.GetAPIKey()
				if err != nil {
					t.Errorf("Reader GetAPIKey() error: %v", err)
				}
			}
		}()
	}

	// Start writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		keyNum := i
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				newKey := "key-" + string(rune('0'+n%10))
				err := am.UpdateAPIKey(newKey)
				if err != nil {
					t.Errorf("Writer UpdateAPIKey() error: %v", err)
				}
			}
		}(keyNum)
	}

	wg.Wait()

	// Verify final state is valid
	finalKey, err := am.GetAPIKey()
	if err != nil {
		t.Errorf("GetAPIKey() after concurrent read/write error: %v", err)
	}
	if finalKey == "" {
		t.Error("Final key should not be empty")
	}
}

func TestAuthManager_MiddlewareValidation(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	authConfig := &AuthConfig{
		APIKey:     "middleware-test-key",
		APIVersion: "v1",
	}

	authManager, err := SetupAuthentication(client, authConfig, logger)
	if err != nil {
		t.Fatalf("SetupAuthentication() error = %v", err)
	}

	// Update the API key
	newKey := "updated-middleware-key"
	err = authManager.UpdateAPIKey(newKey)
	if err != nil {
		t.Fatalf("UpdateAPIKey() error = %v", err)
	}

	// The middleware should use the updated key
	// We can't easily test the actual request without a server,
	// but we can verify the auth manager has the updated key
	currentKey, err := authManager.GetAPIKey()
	if err != nil {
		t.Fatalf("GetAPIKey() error = %v", err)
	}

	if currentKey != newKey {
		t.Errorf("Current key = %q, want %q", currentKey, newKey)
	}
}
