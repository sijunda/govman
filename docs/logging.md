# Logging System Documentation

## Overview

The govman logging system provides a centralized way to handle all user-facing messages with different verbosity levels. It supports three logging levels:


1. **Quiet Level** - Shows only errors
2. **Normal Level** - Shows essential information (default)
3. **Verbose Level** - Shows detailed information including debug messages, timing, internal logs, and progress bars
## Usage

### Command Line Flags

- `--verbose` or `-v`: Enable verbose output
- `--quiet` or `-q`: Enable quiet output (errors only)

### Log Levels

The logger provides several methods for different types of messages, clearly separated between user-facing messages and internal logs:

### User-Facing Messages

These messages are intended to be seen by end users and provide essential information about the operations being performed:

- Error Messages (`Error`, `ErrorWithHelp`)
- Informational Messages (`Info`)
- Success Messages (`Success`)
- Warning Messages (`Warning`)
- Progress Messages (`Progress`)
- Download Messages (`Download`)
- Extraction Messages (`Extract`)
- Verification Messages (`Verify`)

### Internal Logs

These messages are intended for developers and advanced users for debugging and detailed tracing. They are only shown in verbose mode:

- Verbose Messages (`Verbose`)
- Debug Messages (`Debug`)
- Step Messages (`Step`)
- Internal Progress Messages (`InternalProgress`)
- Timing Information (automatically shown with `StartTimer`/`StopTimer`)
- Progress Bars (automatically shown during long operations)

#### Error Messages
```go
_logger.Error("Error message")
_logger.ErrorWithHelp("Error message", "Help text")
```
Always shown regardless of log level.

#### Informational Messages
```go
_logger.Info("Informational message")
```
Shown in normal and verbose modes.

#### Success Messages
```go
_logger.Success("Success message")
```
Shown in normal and verbose modes.

#### Warning Messages
```go
_logger.Warning("Warning message")
```
Shown in normal and verbose modes.

#### Verbose Messages
```go
_logger.Verbose("Verbose message")
```
Only shown in verbose mode.

#### Debug Messages
```go
_logger.Debug("Debug message")
```
Only shown in verbose mode.

#### Progress Messages
```go
_logger.Progress("Progress message")
```
Shown in normal and verbose modes.

#### Download Messages
```go
_logger.Download("Download message")
```
Shown in normal and verbose modes.

#### Extraction Messages
```go
_logger.Extract("Extraction message")
```
Shown in normal and verbose modes.

#### Verification Messages
```go
_logger.Verify("Verification message")
```
Shown in normal and verbose modes.

#### Step Messages
```go
_logger.Step("Step message")
```
Only shown in verbose mode.

#### Internal Progress Messages
```go
_logger.InternalProgress("Internal progress message")
```
Only shown in verbose mode. Used for internal progress tracking that shouldn't be shown to users.

### Progress Bar
The progress bar is automatically displayed during downloads and other long-running operations,
but only when running in verbose mode. In normal mode, no progress bar is shown to keep the
output clean and focused on essential information.

### Timing Operations

The logger includes timing functionality for measuring operation duration:

```go
timer := _logger.StartTimer("operation name")
// ... perform operation ...
_logger.StopTimer(timer)
```

In verbose mode, this will automatically log the start and completion of the operation with timing information.

## Message Formatting

All messages follow a consistent format with appropriate emojis:

- ❌ Error: Error messages
- ℹ️  Info: Informational messages
- ✅ Success: Success messages
- ⚠️  Warning: Warning messages
- 🔍 [VERBOSE] Verbose: Verbose messages
- 🐛 [DEBUG] Debug: Debug messages
- 🔄 Progress: Progress messages
- 📦 Download: Download messages
- 📂 Extract: Extraction messages
- 🔍 Verify: Verification messages
- 📋 Step: Step messages
- 🧹 Clean: Clean messages
- 🔧 Init: Initialization messages
- 🐹 Use: Switch messages

## Best Practices

1. Use appropriate log levels for different types of messages
2. Provide helpful error messages with context
3. Use timing for long-running operations
4. Use consistent formatting with emojis
5. Provide help text with error messages when possible
6. Use step messages to indicate progress through multi-step operations
7. Clearly separate user-facing messages from internal logs
8. Use InternalProgress for internal progress tracking that shouldn't be shown to users
9. Only show progress bars in verbose mode to keep normal output clean

## Examples

```go
// Basic informational message
_logger.Info("Starting installation process")

// Error with help text
_logger.ErrorWithHelp("Failed to download file", "Check your internet connection and try again.")

// Timing a long operation
timer := _logger.StartTimer("download")
// ... download operation ...
_logger.StopTimer(timer)

// Step in a process
_logger.Step("Verifying checksum")