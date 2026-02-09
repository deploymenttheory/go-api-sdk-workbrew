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

	brewCommandsService := brewcommands.NewService(httpClient)

	ctx := context.Background()
	commands, err := brewCommandsService.ListBrewCommands(ctx)
	if err != nil {
		log.Fatalf("Failed to list brew commands: %v", err)
	}

	fmt.Printf("Retrieved %d brew commands\n", len(*commands))
	for i, cmd := range *commands {
		fmt.Printf("\nCommand %d:\n", i+1)
		fmt.Printf("  Command: %s\n", cmd.Command)
		fmt.Printf("  Label: %s\n", cmd.Label)
		fmt.Printf("  Last Updated By: %s\n", cmd.LastUpdatedByUser)
		fmt.Printf("  Started At: %s\n", cmd.StartedAt.String())
		fmt.Printf("  Finished At: %s\n", cmd.FinishedAt.String())
		fmt.Printf("  Devices: %v\n", cmd.Devices)
		fmt.Printf("  Run Count: %d\n", cmd.RunCount)
	}
}
