package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewfiles"
	"go.uber.org/zap"
)

func main() {
	apiKey := os.Getenv("WORKBREW_API_KEY")
	workspace := os.Getenv("WORKBREW_WORKSPACE")

	if apiKey == "" || workspace == "" {
		log.Fatal("WORKBREW_API_KEY and WORKBREW_WORKSPACE environment variables must be set")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	httpClient, err := client.NewClient(apiKey, workspace,
		client.WithLogger(logger),
		client.WithBaseURL("https://console.workbrew.com"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	brewfilesService := brewfiles.NewService(httpClient)

	// Create brewfile request
	// Option 1: Assign to specific devices (comma-separated serial numbers)
	deviceSerials := "TC6R2DHVHG"
	request := &brewfiles.CreateBrewfileRequest{
		Label:               "my-brewfile",
		Content:             "brew \"wget\"\nbrew \"htop\"",
		DeviceSerialNumbers: &deviceSerials,
		// Option 2: Assign to device group (uncomment and set DeviceSerialNumbers to nil)
		// DeviceGroupID: stringPtr("ddba0af6-bd3c-5abf-8311-e62dc6bd9fbc"),
	}

	ctx := context.Background()
	response, _, err := brewfilesService.CreateBrewfile(ctx, request)
	if err != nil {
		log.Fatalf("Failed to create brewfile: %v", err)
	}

	fmt.Printf("Success: %s\n", response.Message)
}
