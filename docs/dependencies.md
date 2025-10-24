# Dependencies

`govman` is designed to be lightweight and has minimal external dependencies. This document outlines the key third-party libraries used in the project and their purpose.

## Core Dependencies

These libraries are essential to the core functionality of `govman`. They are listed in the `go.mod` file.

### `github.com/spf13/cobra`

-   **Purpose**: A powerful library for creating modern CLI applications.
-   **Role in `govman`**:
    -   Defines all user-facing commands (e.g., `install`, `use`, `list`).
    -   Handles command-line argument parsing, flag management, and sub-command routing.
    -   Generates help text and provides a structured way to build the entire CLI.
-   **Location**: Used extensively in the `internal/cli` package.

### `github.com/spf13/viper`

-   **Purpose**: A complete configuration solution for Go applications.
-   **Role in `govman`**:
    -   Manages the application's configuration, which is loaded from `~/.govman/config.yaml`.
    -   Handles default values, ensuring the application works out-of-the-box.
    -   Binds CLI flags to configuration settings (e.g., `--verbose` flag maps to `verbose: true`).
-   **Location**: Used primarily in the `internal/config` package.

## Transitive Dependencies

These are dependencies brought in by the core libraries. While `govman` does not interact with them directly, they are part of the compiled binary.

-   **`github.com/fsnotify/fsnotify`**: Used by Viper for watching for live changes to the configuration file.
-   **`github.com/inconshreveable/mousetrap`**: A helper for Cobra to provide better support on Windows.
-   **`github.com/spf13/pflag`**: A POSIX-compliant flag parser used by Cobra.
-   **`go.yaml.in/yaml.v3`**: Used by Viper for parsing the YAML configuration file.

## Why So Few Dependencies?

The decision to keep the dependency tree small is intentional and provides several benefits:

1.  **Security**: Fewer dependencies mean a smaller attack surface and less exposure to third-party vulnerabilities.
2.  **Stability**: The project is less likely to be affected by breaking changes in upstream libraries.
3.  **Performance**: A smaller binary size and faster compile times.
4.  **Maintainability**: It's easier to manage and audit a small number of well-vetted libraries.

All other functionalities, including the HTTP client, downloader, archive extractor, and symlink manager, are implemented using the Go standard library to ensure maximum reliability and control.