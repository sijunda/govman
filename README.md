# GOVMAN - Go Version Manager

## Project Overview

`govman` is a cross-platform command-line interface (CLI) tool designed to simplify the installation, management, and switching of multiple Go programming language versions. It provides features such as:

*   **Version Management:** Install, uninstall, and switch between different Go versions.
*   **Project-Specific Versions:** Support for defining Go versions on a per-project basis using a `.govman-version` file.
*   **Cross-Platform Compatibility:** Works across Windows, macOS, and Linux operating systems.
*   **Automatic Shell Integration:** Seamless integration with various shell environments.
*   **Efficient Downloads:** Utilizes fast parallel downloads with resume capabilities.

The project is written in **Go** and leverages popular Go libraries such as **Cobra** for building the CLI and **Viper** for configuration management. It interacts with the official Go download API to fetch available releases.

## Quick Installation

The easiest way to install `govman` on Unix-like systems (Linux, macOS, FreeBSD) is by using the provided `install.sh` script:

```bash
curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

This script will download the latest stable release of `govman` for your system, install it to `$HOME/.govman/bin`, and add it to your shell's `PATH`.

## Uninstallation

To uninstall `govman`, you can use the provided `uninstall.sh` script:

```bash
curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/uninstall.sh | bash
```

This script will remove the `govman` binary and remove it from your shell's `PATH`. It will also ask if you want to remove the entire `$HOME/.govman` directory which contains all installed Go versions and cached files.

## Building and Running

The project uses a comprehensive `Makefile` to manage various development and build tasks.

### Prerequisites

*   Go (version 1.25.1 or higher, as specified in `go.mod`)
*   `git` (for version information in builds)
*   `curl` (for the `install.sh` script)

### Development Setup

To set up the development environment and install necessary tools:

```bash
make dev-setup
```

This command will download Go module dependencies and install development tools like `golangci-lint`, `goreleaser`, `goimports`, etc.

### Dependency Management

To download and verify Go module dependencies:

```bash
make deps
```

### Building

*   **Build for the current platform:**
    ```bash
    make build
    ```
    The executable will be placed in the `build/` directory.

*   **Build for all supported platforms:**
    ```bash
    make build-all
    ```
    Cross-compiled binaries will be placed in the `dist/` directory.

### Installation (from source)

*   **Install to your GOPATH/bin:**
    ```bash
    make install
    ```

*   **Install to `/usr/local/bin` (requires `sudo`):**
    ```bash
    make install-local
    ```

### Running

After installation, you can run `govman` commands from your terminal:

```bash
govman --help
```

## Usage

Here are some common `govman` commands:

*   **Install a specific Go version:**
    ```bash
    govman install 1.20.1
    govman install latest
    ```

*   **Switch to a specific Go version:**
    ```bash
    govman use 1.20.1
    ```

*   **Set a Go version as default:**
    ```bash
    govman use 1.20.1 --default
    ```

*   **Set a project-specific Go version:**
    ```bash
    govman use 1.20.1 --local
    ```

*   **List installed Go versions:**
    ```bash
    govman list
    ```

*   **List available Go versions for download:**
    ```bash
    govman list --remote
    ```

*   **Uninstall a Go version:**
    ```bash
    govman uninstall 1.20.1
    ```

*   **Show current active Go version:**
    ```bash
    govman current
    ```

*   **Clean cached Go archives:**
    ```bash
    govman clean
    ```

## Configuration

`govman` uses `viper` for configuration. The default configuration file is located at `$HOME/.govman/config.yaml`. You can also specify a custom configuration file using the `--config` flag.

## Testing

The project includes various test targets:

*   **Run unit tests:**
    ```bash
    make test
    ```

*   **Run tests with coverage analysis:**
    ```bash
    make test-coverage
    ```

*   **Run integration tests:**
    ```bash
    make test-integration
    ```

*   **Run all tests (unit, integration, benchmark):**
    ```bash
    make test-all
    ```

## Code Quality and Linting

*   **Format code using `goimports`:**
    ```bash
    make fmt
    ```

*   **Run `go vet` for static analysis:**
    ```bash
    make vet
    ```

*   **Run comprehensive linting using `golangci-lint`:**
    ```bash
    make lint
    ```

*   **Run all validation checks (fmt, vet, lint):**
    ```bash
    make validate
    ```

## Development Conventions

*   **Build System:** Uses a `Makefile` for consistent task automation.
*   **Code Formatting:** Enforces code formatting using `goimports`.
*   **Linting:** Utilizes `golangci-lint` for comprehensive code quality checks.
*   **Static Analysis:** Employs `go vet` for identifying suspicious constructs.
*   **Testing:** Follows standard Go testing practices with `_test.go` files. Integration tests are specifically tagged.
*   **CLI Framework:** Built with the Cobra library for structured command-line interfaces.
*   **Configuration:** Uses Viper for flexible application configuration.
*   **Version Information:** Build-time version, commit, and branch information are injected into the binary using `ldflags`.
*   **Logging System:** Features a dual-output logging system with separate user-facing and technical logs. See [Logging Documentation](docs/logging.md) for details.

## Contributing

Contributions are welcome! Please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project. (Note: `CONTRIBUTING.md` is a placeholder and needs to be created if not present).

## License

This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.