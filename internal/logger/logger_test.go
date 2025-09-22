package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoggerOutputs(t *testing.T) {
	// Create a new logger
	l := New()

	// Create buffers to capture output
	normalBuf := &bytes.Buffer{}
	verboseBuf := &bytes.Buffer{}

	// Set the writers
	l.SetNormalWriter(normalBuf)
	l.SetVerboseWriter(verboseBuf)

	// Set to verbose level to see all messages
	l.SetLevel(VerboseLevel)

	// Test user-facing messages (should go to normal output)
	l.Info("This is an info message")
	l.Success("This is a success message")
	l.Warning("This is a warning message")
	l.Error("This is an error message")

	// Test technical messages (should go to verbose output)
	l.Verbose("This is a verbose message")
	l.Debug("This is a debug message")
	l.Step("This is a step message")
	l.InternalProgress("This is an internal progress message")

	// Check normal output
	normalOutput := normalBuf.String()
	if !strings.Contains(normalOutput, "This is an info message") {
		t.Errorf("Info message not found in normal output: %s", normalOutput)
	}
	if !strings.Contains(normalOutput, "This is a success message") {
		t.Errorf("Success message not found in normal output: %s", normalOutput)
	}
	if !strings.Contains(normalOutput, "This is a warning message") {
		t.Errorf("Warning message not found in normal output: %s", normalOutput)
	}
	if !strings.Contains(normalOutput, "This is an error message") {
		t.Errorf("Error message not found in normal output: %s", normalOutput)
	}

	// Check verbose output
	verboseOutput := verboseBuf.String()
	if !strings.Contains(verboseOutput, "This is a verbose message") {
		t.Errorf("Verbose message not found in verbose output: %s", verboseOutput)
	}
	if !strings.Contains(verboseOutput, "This is a debug message") {
		t.Errorf("Debug message not found in verbose output: %s", verboseOutput)
	}
	if !strings.Contains(verboseOutput, "This is a step message") {
		t.Errorf("Step message not found in verbose output: %s", verboseOutput)
	}
	if !strings.Contains(verboseOutput, "This is an internal progress message") {
		t.Errorf("Internal progress message not found in verbose output: %s", verboseOutput)
	}

	// Verify that user-facing messages are not in verbose output
	if strings.Contains(verboseOutput, "This is an info message") {
		t.Errorf("Info message found in verbose output when it should be in normal output: %s", verboseOutput)
	}

	// Verify that technical messages are not in normal output
	if strings.Contains(normalOutput, "This is a verbose message") {
		t.Errorf("Verbose message found in normal output when it should be in verbose output: %s", normalOutput)
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	// Create a new logger
	l := New()

	// Create buffers to capture output
	normalBuf := &bytes.Buffer{}
	verboseBuf := &bytes.Buffer{}

	// Set the writers
	l.SetNormalWriter(normalBuf)
	l.SetVerboseWriter(verboseBuf)

	// Test with quiet level - only errors should show
	l.SetLevel(QuietLevel)
	l.Info("This info should not appear")
	l.Error("This error should appear")

	normalOutput := normalBuf.String()
	if strings.Contains(normalOutput, "This info should not appear") {
		t.Errorf("Info message appeared in quiet mode: %s", normalOutput)
	}
	if !strings.Contains(normalOutput, "This error should appear") {
		t.Errorf("Error message did not appear in quiet mode: %s", normalOutput)
	}

	// Clear buffers
	normalBuf.Reset()
	verboseBuf.Reset()

	// Test with normal level - info, success, warning, error should show
	l.SetLevel(NormalLevel)
	l.Info("This info should appear")
	l.Verbose("This verbose should not appear")

	normalOutput = normalBuf.String()
	verboseOutput := verboseBuf.String()

	if !strings.Contains(normalOutput, "This info should appear") {
		t.Errorf("Info message did not appear in normal mode: %s", normalOutput)
	}
	if strings.Contains(verboseOutput, "This verbose should not appear") {
		t.Errorf("Verbose message appeared in normal mode: %s", verboseOutput)
	}
}
