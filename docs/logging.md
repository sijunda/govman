# Logging System

The govman logging system provides a clear separation between user-facing messages and internal technical logs. This allows for clean, user-friendly output while still providing detailed information for debugging and development.

## Log Levels

The logger supports three log levels:

1. **QuietLevel** - Shows only errors
2. **NormalLevel** - Shows essential information (default)
3. **VerboseLevel** - Shows detailed information

## Output Separation

The logger distinguishes between two types of output:

### Normal Output (User-facing)
- Clear, friendly, and non-technical messages
- Suitable for CLI, UI, or command-line tools for users
- Functions: `Info`, `Success`, `Warning`, `Error`, `Progress`, `Download`, `Extract`, `Verify`

### Verbose Output (Technical)
- Detailed technical logs with file names, internal errors, code paths, etc.
- Suitable for development, debugging, CI logs, or advanced diagnostics
- Functions: `Verbose`, `Debug`, `Step`, `InternalProgress`

## Usage

### Basic Usage

```go
import "github.com/sijunda/govman/internal/logger"

// Get the global logger instance
log := logger.Get()

// Set log level
log.SetLevel(logger.VerboseLevel)

// Log user-facing messages
log.Info("Application started")
log.Success("Operation completed")
log.Warning("This is a warning")
log.Error("An error occurred")

// Log technical messages
log.Verbose("Starting process")
log.Debug("Debug information")
log.Step("Processing step 1")
log.InternalProgress("50% complete")
```

### Using Separate Writers

```go
// Create buffers for separate outputs
normalOutput := &bytes.Buffer{}
verboseOutput := &bytes.Buffer{}

// Set separate writers
log.SetNormalWriter(normalOutput)
log.SetVerboseWriter(verboseOutput)

// Now user-facing messages go to normalOutput
// and technical messages go to verboseOutput
```

## Functions

### User-facing Functions
- `Info(format string, args ...interface{})` - General information
- `Success(format string, args ...interface{})` - Success messages
- `Warning(format string, args ...interface{})` - Warning messages
- `Error(format string, args ...interface{})` - Error messages
- `ErrorWithHelp(errorMsg, helpMsg string, args ...interface{})` - Error messages with help text
- `Progress(format string, args ...interface{})` - Progress updates
- `Download(format string, args ...interface{})` - Download status
- `Extract(format string, args ...interface{})` - Extraction status
- `Verify(format string, args ...interface{})` - Verification status

### Technical Functions
- `Verbose(format string, args ...interface{})` - Verbose technical information
- `Debug(format string, args ...interface{})` - Debug information
- `Step(format string, args ...interface{})` - Step in a process
- `InternalProgress(format string, args ...interface{})` - Internal progress updates
- `StartTimer(name string) *Timer` - Start a timer
- `StopTimer(t *Timer)` - Stop a timer and log duration

## Example Output

### Normal Output (User-facing)
```
ℹ️  Application started successfully
✅ Download completed
⚠️  This is a warning message
❌ Error: An error occurred
```

### Verbose Output (Technical)
```
🔍 [VERBOSE] Starting download process
🐛 [DEBUG] Debug information: file size is 10MB
📋 Step 1: Initialize connection
🔄 [INTERNAL] Processed 50% of data