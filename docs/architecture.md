# Govman Architecture

This document provides a high-level overview of the `govman` architecture, its core components, and how they interact.

## Guiding Principles

The architecture is designed to be:

- **Modular**: Each component has a distinct responsibility, making the codebase easy to understand, maintain, and test.
- **Extensible**: New commands, shells, or features can be added with minimal changes to existing code.
- **Cross-Platform**: The codebase is designed to work seamlessly across Windows, macOS, and Linux.
- **User-Centric**: The tool prioritizes a simple and intuitive user experience with zero-config defaults.

## Core Components

The project is divided into several key layers, each with a specific purpose:

1.  **CLI Layer (`internal/cli`)**:
    -   **Purpose**: Handles all user-facing interactions, command parsing, and flag management.
    -   **Technology**: Built on top of the `cobra` library.
    -   **Responsibilities**:
        -   Defines all commands (`install`, `use`, `list`, etc.).
        -   Parses command-line arguments and flags.
        -   Orchestrates calls to the `manager` layer to execute business logic.
        -   Uses the `logger` to provide user feedback.

2.  **Manager Layer (`internal/manager`)**:
    -   **Purpose**: This is the core business logic layer of the application. It orchestrates all version management operations.
    -   **Responsibilities**:
        -   Implements the logic for installing, uninstalling, and switching Go versions.
        -   Coordinates between the `downloader`, `symlink`, `shell`, and `golang` packages.
        -   Manages the state of installed versions and the active version.

3.  **Downloader Layer (`internal/downloader`)**:
    -   **Purpose**: Manages the downloading and extraction of Go archives.
    -   **Features**:
        -   Parallel downloads for faster performance.
        -   Automatic resume of interrupted downloads.
        -   SHA-256 checksum verification to ensure integrity.
        -   Intelligent caching to avoid re-downloading.

4.  **Configuration Layer (`internal/config`)**:
    -   **Purpose**: Manages application settings.
    -   **Technology**: Uses the `viper` library.
    -   **Responsibilities**:
        -   Loads configuration from `~/.govman/config.yaml`.
        -   Provides sensible defaults if no config file is present.
        -   Manages settings like installation directories, download behavior, and mirrors.

5.  **Shell Integration Layer (`internal/shell`)**:
    -   **Purpose**: Handles automatic version switching and PATH management.
    -   **Responsibilities**:
        -   Detects the user's shell (Bash, Zsh, Fish, PowerShell).
        -   Provides the necessary scripts to hook into the shell's environment.
        -   Enables automatic switching based on `.govman-version` files.

6.  **Go Releases Layer (`internal/golang`)**:
    -   **Purpose**: Interacts with the official Go releases API.
    -   **Responsibilities**:
        -   Fetches the list of available Go versions.
        -   Provides metadata for specific versions (e.g., download URL, checksum).
        -   Implements caching to reduce API calls.

## Data Flow

### `govman install latest`

1.  **CLI**: The `install` command is executed with the argument `latest`.
2.  **Manager**: The `Install` method is called with `latest`.
3.  **Manager -> GoReleases**: The manager resolves `latest` to a specific version (e.g., `1.25.1`) by fetching remote versions.
4.  **Manager -> Downloader**: The manager gets the download URL and file info from the `golang` package and passes it to the `Download` method.
5.  **Downloader**:
    -   Downloads the Go archive, showing progress via the `progress` package.
    -   Verifies the SHA-256 checksum.
    -   Extracts the archive to the appropriate version directory (e.g., `~/.govman/versions/go1.25.1`).
6.  **CLI -> Logger**: The CLI reports success to the user.

### `govman use 1.25.1`

1.  **CLI**: The `use` command is executed.
2.  **Manager**: The `Use` method is called.
3.  **Manager**: It verifies that version `1.25.1` is installed.
4.  **Manager -> Symlink**: It creates or updates the `~/.govman/bin/go` symlink to point to the `go` binary inside `~/.govman/versions/go1.25.1/bin/go`.
5.  **Manager -> Shell**: It generates the appropriate shell command to update the `PATH` for the current session.
6.  **CLI -> Logger**: The CLI reports success to the user.