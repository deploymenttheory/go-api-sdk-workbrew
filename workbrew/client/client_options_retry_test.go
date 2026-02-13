package client

import (
	"testing"
	"time"

	"go.uber.org/zap/zaptest"
)

func TestWithRetryWaitTime(t *testing.T) {
	tests := []struct {
		name     string
		waitTime time.Duration
	}{
		{
			name:     "1 second wait",
			waitTime: 1 * time.Second,
		},
		{
			name:     "5 second wait",
			waitTime: 5 * time.Second,
		},
		{
			name:     "10 second wait",
			waitTime: 10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t)

			client, err := NewClient(
				"test-api-key",
				"test-workspace",
				WithLogger(logger),
				WithRetryWaitTime(tt.waitTime),
			)

			if err != nil {
				t.Fatalf("NewClient() error = %v, want nil", err)
			}

			if client == nil {
				t.Fatal("NewClient() returned nil client")
			}

			// Verify the client was created successfully with the option applied
			if client.client == nil {
				t.Error("Client's internal HTTP client is nil")
			}
		})
	}
}

func TestWithRetryMaxWaitTime(t *testing.T) {
	tests := []struct {
		name        string
		maxWaitTime time.Duration
	}{
		{
			name:        "5 second max wait",
			maxWaitTime: 5 * time.Second,
		},
		{
			name:        "30 second max wait",
			maxWaitTime: 30 * time.Second,
		},
		{
			name:        "1 minute max wait",
			maxWaitTime: 60 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t)

			client, err := NewClient(
				"test-api-key",
				"test-workspace",
				WithLogger(logger),
				WithRetryMaxWaitTime(tt.maxWaitTime),
			)

			if err != nil {
				t.Fatalf("NewClient() error = %v, want nil", err)
			}

			if client == nil {
				t.Fatal("NewClient() returned nil client")
			}

			// Verify the client was created successfully with the option applied
			if client.client == nil {
				t.Error("Client's internal HTTP client is nil")
			}
		})
	}
}

func TestWithRetryConfiguration_Combined(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Test combining all retry options
	client, err := NewClient(
		"test-api-key",
		"test-workspace",
		WithLogger(logger),
		WithRetryCount(5),
		WithRetryWaitTime(3*time.Second),
		WithRetryMaxWaitTime(30*time.Second),
	)

	if err != nil {
		t.Fatalf("NewClient() with combined retry options error = %v, want nil", err)
	}

	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}

	// Verify the client was created successfully
	if client.client == nil {
		t.Error("Client's internal HTTP client is nil")
	}
}

func TestWithRetryWaitTime_ZeroDuration(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Test with zero duration (edge case)
	client, err := NewClient(
		"test-api-key",
		"test-workspace",
		WithLogger(logger),
		WithRetryWaitTime(0),
	)

	if err != nil {
		t.Fatalf("NewClient() with zero wait time error = %v, want nil", err)
	}

	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
}

func TestWithRetryMaxWaitTime_ZeroDuration(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Test with zero duration (edge case)
	client, err := NewClient(
		"test-api-key",
		"test-workspace",
		WithLogger(logger),
		WithRetryMaxWaitTime(0),
	)

	if err != nil {
		t.Fatalf("NewClient() with zero max wait time error = %v, want nil", err)
	}

	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
}
