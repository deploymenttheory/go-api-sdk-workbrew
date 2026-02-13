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

	ctx := context.Background()
	brewfilesList, _, err := brewfilesService.ListBrewfiles(ctx)
	if err != nil {
		log.Fatalf("Failed to list brewfiles: %v", err)
	}

	fmt.Printf("Retrieved %d brewfiles\n", len(*brewfilesList))
	for i, bf := range *brewfilesList {
		fmt.Printf("\nBrewfile %d:\n", i+1)
		fmt.Printf("  Label: %s\n", bf.Label)
		fmt.Printf("  Last Updated By: %s\n", bf.LastUpdatedByUser)
		fmt.Printf("  Started At: %s\n", bf.StartedAt)
		fmt.Printf("  Finished At: %s\n", bf.FinishedAt)
		fmt.Printf("  Devices: %v\n", bf.Devices)
		fmt.Printf("  Run Count: %d\n", bf.RunCount)
	}
}
