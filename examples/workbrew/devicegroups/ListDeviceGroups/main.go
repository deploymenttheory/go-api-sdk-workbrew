package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devicegroups"
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

	deviceGroupsService := devicegroups.NewService(httpClient)

	ctx := context.Background()
	groups, err := deviceGroupsService.ListDeviceGroups(ctx)
	if err != nil {
		log.Fatalf("Failed to list device groups: %v", err)
	}

	fmt.Printf("Retrieved %d device groups\n", len(*groups))
	for i, group := range *groups {
		fmt.Printf("\nGroup %d:\n", i+1)
		fmt.Printf("  ID: %s\n", group.ID)
		fmt.Printf("  Name: %s\n", group.Name)
		fmt.Printf("  Devices: %v\n", group.Devices)
	}
}
