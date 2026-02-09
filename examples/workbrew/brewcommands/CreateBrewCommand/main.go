package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewcommands"
)

func main() {
	fmt.Println("=== Workbrew - Brew Commands Example ===")

	// API credentials
	apiKey := "your-api-key-here"
	workspaceName := "your-workspace"

	// Create client
	wbClient, err := workbrew.NewClient(
		apiKey,
		workspaceName,
		client.WithDebug(),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Get all brew commands
	fmt.Println("\n=== Example 1: Get All Brew Commands ===")

	commands, err := wbClient.BrewCommands.GetBrewCommands(ctx)
	if err != nil {
		log.Fatalf("Failed to get brew commands: %v", err)
	}

	fmt.Printf("Found %d brew commands:\n\n", len(*commands))

	for i, cmd := range *commands {
		fmt.Printf("Command %d:\n", i+1)
		fmt.Printf("  Command: %s\n", cmd.Command)
		fmt.Printf("  Label: %s\n", cmd.Label)
		fmt.Printf("  Last Updated By: %s\n", cmd.LastUpdatedByUser)
		fmt.Printf("  Run Count: %d\n", cmd.RunCount)
		fmt.Printf("  Started At: %s\n", cmd.StartedAt.String())
		fmt.Printf("  Finished At: %s\n", cmd.FinishedAt.String())
		if len(cmd.Devices) > 0 {
			fmt.Printf("  Devices: %v\n", cmd.Devices)
		}
		fmt.Println()
	}

	// Example 2: Create a new brew command
	fmt.Println("\n=== Example 2: Create New Brew Command ===")

	recurrence := brewcommands.RecurrenceOnce
	createRequest := &brewcommands.CreateBrewCommandRequest{
		Arguments:  "install wget",
		Recurrence: &recurrence,
	}

	createResp, err := wbClient.BrewCommands.CreateBrewCommand(ctx, createRequest)
	if err != nil {
		// Check for specific error types
		if client.IsForbidden(err) {
			fmt.Printf("Error: This feature requires a paid plan\n")
		} else if client.IsValidationError(err) {
			fmt.Printf("Validation error: %v\n", err)
		} else {
			log.Fatalf("Failed to create brew command: %v", err)
		}
	} else {
		fmt.Printf("✓ %s\n", createResp.Message)
	}

	// Example 3: Get brew command runs
	fmt.Println("\n=== Example 3: Get Brew Command Runs ===")

	if len(*commands) > 0 {
		firstCommand := (*commands)[0]
		fmt.Printf("Getting runs for command: %s\n", firstCommand.Label)

		runs, err := wbClient.BrewCommands.GetBrewCommandRuns(ctx, firstCommand.Label)
		if err != nil {
			log.Printf("Failed to get brew command runs: %v", err)
		} else {
			fmt.Printf("Found %d runs:\n\n", len(*runs))

			for i, run := range *runs {
				fmt.Printf("Run %d:\n", i+1)
				fmt.Printf("  Device: %s\n", run.Device)
				fmt.Printf("  Success: %v\n", run.Success)
				fmt.Printf("  Started: %s\n", run.StartedAt.String())
				fmt.Printf("  Finished: %s\n", run.FinishedAt.String())
				if run.Output != "" {
					fmt.Printf("  Output: %s\n", run.Output)
				}
				fmt.Println()
			}
		}
	}

	// Example 4: Get brew commands CSV
	fmt.Println("\n=== Example 4: Get Brew Commands (CSV) ===")

	csvData, err := wbClient.BrewCommands.GetBrewCommandsCSV(ctx)
	if err != nil {
		log.Printf("Failed to get brew commands CSV: %v", err)
	} else {
		fmt.Printf("Retrieved CSV data (%d bytes):\n", len(csvData))
		fmt.Println(string(csvData))
	}

	// Example 5: Pretty print JSON
	fmt.Println("\n=== Example 5: Full JSON Response ===")
	jsonData, err := json.MarshalIndent(commands, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
	} else {
		fmt.Println(string(jsonData))
	}

	fmt.Println("\n=== Example Complete ===")
}
