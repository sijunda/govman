package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sijunda/govman/internal/logger"
)

func main() {
	// Create a new logger
	l := logger.New()

	// Create buffers to capture normal and verbose output separately
	normalOutput := &bytes.Buffer{}
	verboseOutput := &bytes.Buffer{}

	// Set separate writers for normal and verbose output
	l.SetNormalWriter(normalOutput)
	l.SetVerboseWriter(verboseOutput)

	// Set to verbose level to see all messages
	l.SetLevel(logger.VerboseLevel)

	// Simulate some application activity
	fmt.Println("=== Simulating Application Activity ===")

	// Log user-facing messages (these will go to normal output)
	l.Info("Application started successfully")
	l.Download("Downloading Go 1.21.0")
	l.Extract("Extracting archive...")
	l.Success("Installation completed successfully")
	l.Warning("This is a warning message")
	l.Error("An error occurred")

	// Log technical messages (these will go to verbose output)
	l.Verbose("Starting download process")
	l.Debug("Debug information: file size is 10MB")
	l.Step("Step 1: Initialize connection")
	l.Step("Step 2: Download file")
	l.Step("Step 3: Verify checksum")
	l.InternalProgress("Processed 50%% of data")

	fmt.Println("\n=== Normal Output (User-facing) ===")
	// Print normal output line by line to preserve formatting
	normalLines := strings.Split(strings.TrimSpace(normalOutput.String()), "\n")
	for _, line := range normalLines {
		fmt.Println(line)
	}

	fmt.Println("\n=== Verbose Output (Technical) ===")
	// Print verbose output line by line to preserve formatting
	verboseLines := strings.Split(strings.TrimSpace(verboseOutput.String()), "\n")
	for _, line := range verboseLines {
		fmt.Println(line)
	}
}
