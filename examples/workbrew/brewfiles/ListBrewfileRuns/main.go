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
	brewfileLabel := os.Getenv("BREWFILE_LABEL") // e.g., "bundle-file"

	if apiKey == "" || workspace == "" || brewfileLabel == "" {
		log.Fatal("WORKBREW_API_KEY, WORKBREW_WORKSPACE, and BREWFILE_LABEL environment variables must be set")
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
	runs, err := brewfilesService.ListBrewfileRuns(ctx, brewfileLabel)
	if err != nil {
		log.Fatalf("Failed to list brewfile runs: %v", err)
	}

	fmt.Printf("Retrieved %d runs for brewfile '%s'\n", len(*runs), brewfileLabel)
	for i, run := range *runs {
		fmt.Printf("\nRun %d:\n", i+1)
		fmt.Printf("  Label: %s\n", run.Label)
		fmt.Printf("  Device: %s\n", run.Device)
		fmt.Printf("  Success: %t\n", run.Success)
		fmt.Printf("  Started At: %s\n", run.StartedAt)
		fmt.Printf("  Finished At: %s\n", run.FinishedAt)
		fmt.Printf("  Output: %s\n", run.Output)
	}
}
