package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/events"
)

func main() {
	fmt.Println("=== Workbrew - Events Example ===")

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

	// Example 1: Get all events (no filter)
	fmt.Println("\n=== Example 1: Get All Events (No Filter) ===")

	allEvents, err := wbClient.Events.GetEvents(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to get events: %v", err)
	}

	fmt.Printf("Found %d events:\n\n", len(*allEvents))

	// Show first 5 events
	for i, event := range *allEvents {
		if i >= 5 {
			break
		}
		fmt.Printf("Event %d:\n", i+1)
		fmt.Printf("  ID: %s\n", event.ID)
		fmt.Printf("  Type: %s\n", event.EventType)
		fmt.Printf("  Occurred At: %s\n", event.OccurredAt.Format("2006-01-02 15:04:05"))

		if event.ActorType != nil {
			fmt.Printf("  Actor Type: %s\n", *event.ActorType)
		}
		if event.TargetType != nil {
			fmt.Printf("  Target Type: %s\n", *event.TargetType)
		}
		if event.TargetIdentifier != nil {
			fmt.Printf("  Target Identifier: %s\n", *event.TargetIdentifier)
		}

		fmt.Println()
	}

	// Example 2: Filter events by user actions
	fmt.Println("\n=== Example 2: Get User Events (Filtered) ===")

	userEvents, err := wbClient.Events.GetEvents(ctx, &events.RequestQueryOptions{
		Filter: "user",
	})
	if err != nil {
		log.Fatalf("Failed to get user events: %v", err)
	}

	fmt.Printf("Found %d user events\n", len(*userEvents))

	// Example 3: Filter events by system actions
	fmt.Println("\n=== Example 3: Get System Events (Filtered) ===")

	systemEvents, err := wbClient.Events.GetEvents(ctx, &events.RequestQueryOptions{
		Filter: "system",
	})
	if err != nil {
		log.Fatalf("Failed to get system events: %v", err)
	}

	fmt.Printf("Found %d system events\n", len(*systemEvents))

	// Example 4: Get events in CSV format
	fmt.Println("\n=== Example 4: Get Events (CSV) ===")

	csvData, err := wbClient.Events.GetEventsCSV(ctx, &events.RequestQueryOptions{
		Filter: "all",
	})
	if err != nil {
		log.Fatalf("Failed to get events CSV: %v", err)
	}

	fmt.Printf("Retrieved CSV data (%d bytes):\n", len(csvData))
	fmt.Println(string(csvData[:min(500, len(csvData))]) + "...")

	// Example 5: Get events CSV with download flag
	fmt.Println("\n=== Example 5: Get Events CSV (With Download Flag) ===")

	downloadCSV, err := wbClient.Events.GetEventsCSV(ctx, &events.RequestQueryOptions{
		Filter:   "user",
		Download: true, // Forces download=1 query parameter
	})
	if err != nil {
		log.Fatalf("Failed to get events CSV with download: %v", err)
	}

	fmt.Printf("Retrieved downloadable CSV data (%d bytes)\n", len(downloadCSV))

	// Example 6: Inspect event changes (for update events)
	fmt.Println("\n=== Example 6: Events with Changes ===")

	for i, event := range *allEvents {
		if i >= 10 { // Check first 10 events
			break
		}
		if event.Changes != nil && len(event.Changes) > 0 {
			fmt.Printf("\nEvent: %s\n", event.EventType)
			fmt.Printf("Changes:\n")
			changesJSON, _ := json.MarshalIndent(event.Changes, "  ", "  ")
			fmt.Printf("  %s\n", string(changesJSON))
		}
	}

	// Example 7: Pretty print JSON response
	fmt.Println("\n=== Example 7: Full JSON Response (First 3 Events) ===")

	if len(*allEvents) > 3 {
		sample := (*allEvents)[:3]
		jsonData, err := json.MarshalIndent(sample, "", "  ")
		if err != nil {
			log.Printf("Error marshaling response to JSON: %v", err)
		} else {
			fmt.Println(string(jsonData))
		}
	}

	fmt.Println("\n=== Example Complete ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
