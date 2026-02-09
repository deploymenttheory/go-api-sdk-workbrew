package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewcommands"
	"go.uber.org/zap"
)

func main() {
	apiKey := os.Getenv("WORKBREW_API_KEY")
	workspace := os.Getenv("WORKBREW_WORKSPACE")
	brewCommandLabel := os.Getenv("BREW_COMMAND_LABEL") // e.g., "outdated"

	if apiKey == "" || workspace == "" || brewCommandLabel == "" {
		log.Fatal("WORKBREW_API_KEY, WORKBREW_WORKSPACE, and BREW_COMMAND_LABEL environment variables must be set")
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

	brewCommandsService := brewcommands.NewService(httpClient)

	ctx := context.Background()
	csvData, err := brewCommandsService.ListBrewCommandRunsCSV(ctx, brewCommandLabel)
	if err != nil {
		log.Fatalf("Failed to list brew command runs CSV: %v", err)
	}

	fmt.Printf("Brew Command Runs CSV for '%s' (%d bytes):\n", brewCommandLabel, len(csvData))
	fmt.Println(string(csvData))
}
