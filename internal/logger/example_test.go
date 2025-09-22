package logger

import (
	"bytes"
	"fmt"
)

// Example demonstrates how to use the logger with separate normal and verbose outputs
func Example() {
	// Create a new logger
	l := New()

	// Create buffers to capture normal and verbose output separately
	normalOutput := &bytes.Buffer{}
	verboseOutput := &bytes.Buffer{}

	// Set separate writers for normal and verbose output
	l.SetNormalWriter(normalOutput)
	l.SetVerboseWriter(verboseOutput)

	// Set to verbose level to see all messages
	l.SetLevel(VerboseLevel)

	// Log user-facing messages (these will go to normal output)
	l.Info("Application started successfully")
	l.Success("Download completed")
	l.Warning("This is a warning message")
	l.Error("An error occurred")

	// Log technical messages (these will go to verbose output)
	l.Verbose("Starting download process")
	l.Debug("Debug information: file size is 10MB")
	l.Step("Step 1: Initialize connection")
	l.InternalProgress("Processed 50%% of data")

	// Show what would be displayed to a regular user
	fmt.Println("=== Normal Output (User-facing) ===")
	fmt.Print(normalOutput.String())

	// Show what would be displayed to a developer
	fmt.Println("\n=== Verbose Output (Technical) ===")
	fmt.Print(verboseOutput.String())

	// Output:
	// === Normal Output (User-facing) ===
	// ℹ️  Application started successfully
	// ✅ Download completed
	// ⚠️  This is a warning message
	// ❌ Error: An error occurred
	//
	// === Verbose Output (Technical) ===
	// 🔍 [VERBOSE] Starting download process
	// 🐛 [DEBUG] Debug information: file size is 10MB
	// 📋 Step 1: Initialize connection
	// 🔄 [INTERNAL] Processed 50% of data
}
