package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
)

func main() {
	fmt.Println("=== Workbrew - Get Devices Example ===")

	// API credentials
	apiKey := "your-api-key-here"
	workspaceName := "your-workspace"

	// Create client
	wbClient, err := workbrew.NewClient(
		apiKey,
		workspaceName,
		client.WithDebug(), // Enable debug mode to see HTTP requests
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Get all devices in JSON format
	fmt.Println("\n=== Example 1: Get All Devices (JSON) ===")

	devices, err := wbClient.Devices.GetDevices(ctx)
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}

	fmt.Printf("Found %d devices:\n\n", len(*devices))

	for i, device := range *devices {
		fmt.Printf("Device %d:\n", i+1)
		fmt.Printf("  Serial Number: %s\n", device.SerialNumber)
		fmt.Printf("  Device Type: %s\n", device.DeviceType)
		fmt.Printf("  OS Version: %s\n", device.OSVersion)
		fmt.Printf("  Homebrew Version: %s\n", device.HomebrewVersion)
		fmt.Printf("  Workbrew Version: %s\n", device.WorkbrewVersion)
		fmt.Printf("  Formulae Count: %d\n", device.FormulaeCount)
		fmt.Printf("  Casks Count: %d\n", device.CasksCount)
		fmt.Printf("  Last Seen: %s\n", device.LastSeenAt.String())
		fmt.Printf("  Command Last Run: %s\n", device.CommandLastRunAt.String())

		if device.MDMUserOrDeviceName != nil {
			fmt.Printf("  MDM Name: %s\n", *device.MDMUserOrDeviceName)
		}

		if len(device.Groups) > 0 {
			fmt.Printf("  Groups: %v\n", device.Groups)
		}

		fmt.Println()
	}

	// Example 2: Get devices in CSV format
	fmt.Println("\n=== Example 2: Get All Devices (CSV) ===")

	csvData, err := wbClient.Devices.GetDevicesCSV(ctx)
	if err != nil {
		log.Fatalf("Failed to get devices CSV: %v", err)
	}

	fmt.Printf("Retrieved CSV data (%d bytes):\n", len(csvData))
	fmt.Println(string(csvData))

	// Example 3: Pretty print JSON response
	fmt.Println("\n=== Example 3: Full JSON Response ===")
	jsonData, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
	} else {
		fmt.Println(string(jsonData))
	}

	fmt.Println("\n=== Example Complete ===")
}
