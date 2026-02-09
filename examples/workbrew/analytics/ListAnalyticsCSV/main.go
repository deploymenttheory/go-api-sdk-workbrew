package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/analytics"
	"go.uber.org/zap"
)

func main() {
	// Retrieve API key and workspace from environment variables
	apiKey := os.Getenv("WORKBREW_API_KEY")
	workspace := os.Getenv("WORKBREW_WORKSPACE")

	if apiKey == "" || workspace == "" {
		log.Fatal("WORKBREW_API_KEY and WORKBREW_WORKSPACE environment variables must be set")
	}

	// Create logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create HTTP client
	httpClient, err := client.NewClient(apiKey, workspace,
		client.WithLogger(logger),
		client.WithBaseURL("https://console.workbrew.com"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create analytics service
	analyticsService := analytics.NewService(httpClient)

	// List analytics as CSV
	ctx := context.Background()
	csvData, err := analyticsService.ListAnalyticsCSV(ctx)
	if err != nil {
		log.Fatalf("Failed to list analytics CSV: %v", err)
	}

	// Print CSV data
	fmt.Printf("Analytics CSV (%d bytes):\n", len(csvData))
	fmt.Println(string(csvData))
}
