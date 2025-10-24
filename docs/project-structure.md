# Project Structure

This document outlines the directory structure of the `govman` project, explaining the purpose of each key directory and file.

```
govman/
├── cmd/govman/            # Main application entry point
│   └── main.go
├── internal/              # Internal packages (not for external use)
│   ├── cli/               # Command-line interface logic
│   ├── config/            # Configuration management
│   ├── downloader/        # Download engine
│   ├── golang/            # Go releases API client
│   ├── logger/            # Structured logging
│   ├── manager/           # Core version management logic
│   ├── progress/          # Download progress bars
│   ├── shell/             # Shell integration and auto-switching
│   ├── symlink/           # Cross-platform symlink utilities
│   ├── util/              # Shared utility functions
│   └── version/           # Build version information
├── scripts/               # Installation and uninstallation scripts
├── docs/                  # Project documentation
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── Makefile               # Build automation
└── README.md              # Project overview
```

## Top-Level Directories

### `cmd/govman/`

-   **Purpose**: The main entry point of the `govman` executable.
-   **`main.go`**: Contains the `main` function, which initializes and executes the root `cobra` command from the `internal/cli` package. It's responsible for handling top-level errors and setting the process exit code.

### `internal/`

-   **Purpose**: This directory contains all the core logic of the application. As per Go conventions, packages inside `internal/` are not importable by external applications.
-   See the detailed breakdown of `internal/` packages below.

### `scripts/`

-   **Purpose**: Contains shell scripts for easy installation and uninstallation on different platforms.
-   **`install.sh`**: Bash script for macOS and Linux.
-   **`install.ps1`**: PowerShell script for Windows.
-   **`install.bat`**: Batch script for older Windows systems.
-   **`uninstall.*`**: Corresponding scripts to remove `govman` safely.

### `docs/`

-   **Purpose**: Holds all technical and user-facing documentation for the project.

## `internal/` Packages

### `internal/cli`

-   **Responsibility**: Defines the entire command-line interface.
-   Each file (`install.go`, `use.go`, etc.) corresponds to a specific CLI command.
-   It uses the `cobra` library to structure commands, arguments, and flags.
-   This layer is responsible for parsing user input and calling the appropriate methods in the `manager` package.

### `internal/manager`

-   **Responsibility**: The brain of the application. It contains the core business logic for managing Go versions.
-   It acts as an orchestrator, coordinating actions between the `downloader`, `symlink`, `shell`, and `config` packages to fulfill user commands.

### `internal/config`

-   **Responsibility**: Manages application settings.
-   It uses `viper` to load, parse, and validate the `config.yaml` file. It also provides default values for all configuration options.

### `internal/downloader`

-   **Responsibility**: Handles all file download operations.
-   Implements features like parallel chunk downloading, resuming interrupted downloads, and verifying file integrity via checksums.

### `internal/golang`

-   **Responsibility**: A client for the official Go releases API.
-   Fetches the list of available Go versions and their metadata (download URLs, checksums, file sizes). It includes a caching mechanism to avoid excessive API calls.

### `internal/logger`

-   **Responsibility**: Provides a structured, multi-level logging system.
-   Supports different log levels (`Info`, `Verbose`, `Error`, `Debug`) and allows output to be controlled via CLI flags (`--quiet`, `--verbose`).

### `internal/progress`

-   **Responsibility**: Renders progress bars for downloads.
-   Provides a visual representation of download speed, ETA, and completion percentage.

### `internal/shell`

-   **Responsibility**: Manages shell-specific integration.
-   Detects the user's shell and provides the correct scripts and commands for `PATH` modification and automatic version switching.

### `internal/symlink`

-   **Responsibility**: A small, cross-platform utility for creating and managing symbolic links, which is the core mechanism for switching the active Go version.

### `internal/util`

-   **Responsibility**: Contains shared helper functions used across the application, such as `FormatBytes` for human-readable file sizes.

### `internal/version`

-   **Responsibility**: Manages the application's own version information.
-   Build-time variables (`Version`, `Commit`, `Date`) are injected into this package via linker flags in the `Makefile`.