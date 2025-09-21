# PromptAI for Generating GOVMAN Landing Page

## Introduction
GOVMAN is a powerful, cross-platform Go version manager that simplifies the installation, management, and switching between multiple Go versions. Whether you're a developer working on multiple projects with different Go requirements or a system administrator managing Go environments, GOVMAN provides an intuitive and efficient solution.

Built with performance and ease-of-use in mind, GOVMAN offers lightning-fast installation, zero-configuration setup, and seamless shell integration. With support for project-specific version management through `.go-version` files, automatic version switching, and cross-platform compatibility, GOVMAN streamlines your Go development workflow.

## SEO Considerations
- **Meta Tags:** Include title, description, and keywords.
- **Alt Text:** Provide alt text for all images.
- **URL Structure:** Use a clean and descriptive URL structure.
- **Header Tags:** Use H1, H2, and H3 tags appropriately.
- **Content Quality:** Write high-quality, informative, and engaging content.
- **Internal Links:** Include internal links to other relevant pages within the site.
- **Mobile Responsiveness:** Ensure the page is mobile-friendly.
- **Loading Speed:** Optimize images and assets for faster loading times.
- **Schema Markup:** Add schema markup for better search engine understanding.

## Key Features
- **🚀 Lightning-fast Installation:** Install and switch between Go versions quickly with optimized download and extraction processes.
- **🎯 Zero Configuration:** Works out of the box with no setup required - just install and start using.
- **📁 Project-Specific Versions:** Support for defining Go versions on a per-project basis using a `.go-version` file for automatic version switching.
- **🚫 No Admin/Sudo Required:** Fully userspace installation for enhanced security and ease of use.
- **💾 Intelligent Caching:** Offline mode support with smart caching of downloaded versions.
- **📦 Parallel Downloads:** Fast parallel downloads with automatic resume on failure for reliable installations.
- **🌍 Cross-Platform Support:** Seamless operation across Windows, macOS, Linux, and ARM architectures.
- **🧹 Built-in Cleanup:** Efficient disk space management with built-in cache cleaning tools.
- **🐚 Automatic Shell Integration:** Seamless integration with bash, zsh, fish, and PowerShell environments.
- **🔄 Flexible Version Management:** Install, uninstall, list, and switch between multiple Go versions with ease.

## Quick Installation
The easiest way to install `govman` on Unix-like systems (Linux, macOS, FreeBSD) is by using the provided `install.sh` script:
```bash
curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```
This script will download the latest stable release of `govman` for your system, install it to `$HOME/.govman/bin`, and add it to your shell's `PATH`.

For Windows users, you can download the binary directly from the [GitHub releases page](https://github.com/sijunda/govman/releases).

## Getting Started
Getting started with GOVMAN is straightforward. After installation, you can immediately begin managing Go versions.

### Initial Setup
After installation, initialize shell integration for automatic version switching:
```bash
govman init
```
This command adds GOVMAN to your shell configuration and enables auto-switching based on `.go-version` files.

### Basic Usage Examples
1. **Install a Go version:**
   ```bash
   govman install latest
   ```

2. **Switch to a specific Go version:**
   ```bash
   govman use 1.25.1
   ```

3. **List installed versions:**
   ```bash
   govman list
   ```

4. **Set project-specific version:**
   ```bash
   govman use 1.25.1 --local
   ```
   This creates a `.go-version` file in your project directory.

## Documentation
GOVMAN comes with comprehensive documentation to help you get the most out of the tool.

### Installation
GOVMAN can be installed in multiple ways:
1. **Using the installation script (recommended for Unix-like systems):**
   ```bash
   curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/install.sh | bash
   ```
2. **Manual installation:** Download binaries directly from the GitHub releases page.
3. **From source:** Clone the repository and build from source using `make build`.

### Configuration
GOVMAN uses Viper for configuration management. The default configuration file is located at `$HOME/.govman/config.yaml`. You can specify a custom configuration file using the `--config` flag.

Key configuration options include:
- `install_dir` - Directory where Go versions are installed
- `cache_dir` - Directory for caching downloaded archives
- `default_version` - Default Go version to use

### Usage
GOVMAN provides intuitive commands for all version management tasks:
- Install versions with `govman install`
- Switch versions with `govman use`
- List versions with `govman list`
- Manage project-specific versions with `.go-version` files

### Troubleshooting
Common issues and solutions:
- **Permission errors:** Ensure you have write permissions to the installation directory
- **Shell integration not working:** Run `govman init` and restart your shell
- **Network issues:** Check your internet connection and proxy settings if behind a corporate firewall

## API Reference
GOVMAN provides a comprehensive command-line interface for managing Go versions.

### Core Commands
- `govman install [version...]` - Install one or more Go versions
- `govman uninstall <version>` - Uninstall a Go version
- `govman use <version>` - Switch to a specific Go version
- `govman list` - List installed Go versions
- `govman current` - Show current active Go version
- `govman info <version>` - Show detailed information about an installed Go version
- `govman clean` - Clean download cache
- `govman init` - Initialize shell integration
- `govman selfupdate` - Update govman to the latest version

### Command Flags
- `--verbose` or `-v` - Enable verbose output for detailed logging
- `--quiet` or `-q` - Enable quiet output (errors only)
- `--config` - Specify a custom configuration file
- `--default` or `-d` - Set as system default version (with `use` command)
- `--local` or `-l` - Set as project-local version (with `use` command)
- `--remote` or `-r` - List available versions for download (with `list` command)
- `--force` or `-f` - Force re-initialization (with `init` command)
- `--shell` - Target shell for initialization (with `init` command)
- `--check` - Check for updates without installing (with `selfupdate` command)

### Environment Variables
- `GOVMAN_CONFIG` - Path to configuration file
- `GOVMAN_HOME` - GOVMAN installation directory
- `GOVMAN_NO_COLOR` - Disable colored output

## FAQ
### What is GOVMAN?
GOVMAN is a cross-platform command-line tool for managing multiple Go versions. It allows you to easily install, switch between, and manage different versions of Go on your system.

### How do I install GOVMAN?
The easiest way to install GOVMAN on Unix-like systems is by using the provided installation script:
```bash
curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/install.sh | bash
```
For Windows users, you can download the binary directly from the releases page on GitHub.

### Can I use GOVMAN on Windows?
Yes, GOVMAN fully supports Windows along with macOS and Linux. It works seamlessly across all major operating systems and architectures including ARM.

### How do I report issues or contribute to the project?
You can report issues or contribute to the project by visiting the GitHub repository at https://github.com/sijunda/govman. Feel free to open issues for bug reports or feature requests, and submit pull requests for contributions.

### How does GOVMAN handle shell integration?
GOVMAN automatically integrates with popular shells including bash, zsh, fish, and PowerShell. The `govman init` command sets up the necessary configurations for automatic version switching based on project-specific `.go-version` files.

### Does GOVMAN require administrator privileges?
No, GOVMAN operates entirely in userspace and does not require administrator or sudo privileges for normal operation. This enhances security and makes it easier to install and use.

## Contact
For support, feature requests, or contributions, please reach out through the following channels:
- **GitHub Issues:** [https://github.com/sijunda/govman/issues](https://github.com/sijunda/govman/issues)
- **Email:** [sijun.danang@gmail.com](mailto:sijun.danang@gmail.com)

## Meta Tags
- **Title:** GOVMAN - Go Version Manager | Simplify Go Version Management
- **Description:** GOVMAN is a cross-platform CLI tool for managing multiple Go versions. Install, switch, and manage Go versions easily.
- **Keywords:** GOVMAN, Go Version Manager, Go, CLI, version management, cross-platform, shell integration, efficient downloads

## Alt Text
Ensure all images have descriptive alt text. For example:
- ![GOVMAN Logo](logo.png "GOVMAN Logo")

## URL Structure
Use a clean and descriptive URL structure. For example:
- `https://govman.io/`
- `https://govman.io/features`
- `https://govman.io/installation`

## Header Tags
Use H1, H2, and H3 tags appropriately. For example:
- `# GOVMAN - Go Version Manager`
- `## Quick Installation`
- `### Prerequisites`

## Content Quality
Write high-quality, informative, and engaging content. Ensure the page is well-written and easy to understand.

## Internal Links
Include internal links to other relevant pages within the site. For example:
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)

## Mobile Responsiveness
Ensure the page is mobile-friendly. Use responsive design techniques to ensure the page looks good on all devices.

## Loading Speed
Optimize images and assets for faster loading times. Use techniques such as image compression and lazy loading.

## Schema Markup
Add schema markup for better search engine understanding. For example:
- Use JSON-LD to mark up the page structure and content.
- Include schema for software application, organization, and product.
```bash
curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/install.sh | bash
```
This script will download the latest stable release of `govman` for your system, install it to `$HOME/.govman/bin`, and add it to your shell's `PATH`.

## Building and Running
### Prerequisites
- Go (version 1.25.1 or higher, as specified in `go.mod`)
- `git` (for version information in builds)
- `curl` (for the `install.sh` script)

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
- **Build for the current platform:**
  ```bash
  make build
  ```
  The executable will be placed in the `build/` directory.
- **Build for all supported platforms:**
  ```bash
  make build-all
  ```
  Cross-compiled binaries will be placed in the `dist/` directory.

### Installation (from source)
- **Install to your GOPATH/bin:**
  ```bash
  make install
  ```
- **Install to `/usr/local/bin` (requires `sudo`):**
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
- **Install a specific Go version:**
  ```bash
  govman install 1.20.1
  govman install latest
  ```
- **Switch to a specific Go version:**
  ```bash
  govman use 1.20.1
  ```
- **Set a Go version as default:**
  ```bash
  govman use 1.20.1 --default
  ```
- **Set a project-specific Go version:**
  ```bash
  govman use 1.20.1 --local
  ```
- **List installed Go versions:**
  ```bash
  govman list
  ```
- **List available Go versions for download:**
  ```bash
  govman list --remote
  ```
- **Uninstall a Go version:**
  ```bash
  govman uninstall 1.20.1
  ```
- **Show current active Go version:**
  ```bash
  govman current
  ```
- **Clean cached Go archives:**
  ```bash
  govman clean
  ```

## Configuration
`govman` uses `viper` for configuration. The default configuration file is located at `$HOME/.govman/config.yaml`. You can also specify a custom configuration file using the `--config` flag.

## Testing
The project includes various test targets:
- **Run unit tests:**
  ```bash
  make test
  ```
- **Run tests with coverage analysis:**
  ```bash
  make test-coverage
  ```
- **Run integration tests:**
  ```bash
  make test-integration
  ```
- **Run all tests (unit, integration, benchmark):**
  ```bash
  make test-all
  ```

## Code Quality and Linting
- **Format code using `goimports`:**
  ```bash
  make fmt
  ```
- **Run `go vet` for static analysis:**
  ```bash
  make vet
  ```
- **Run comprehensive linting using `golangci-lint`:**
  ```bash
  make lint
  ```
- **Run all validation checks (fmt, vet, lint):**
  ```bash
  make validate
  ```

## Development Conventions
- **Build System:** Uses a `Makefile` for consistent task automation.
- **Code Formatting:** Enforces code formatting using `goimports`.
- **Linting:** Utilizes `golangci-lint` for comprehensive code quality checks.
- **Static Analysis:** Employs `go vet` for identifying suspicious constructs.
- **Testing:** Follows standard Go testing practices with `_test.go` files. Integration tests are specifically tagged.
- **CLI Framework:** Built with the Cobra library for structured command-line interfaces.
- **Configuration:** Uses Viper for flexible application configuration.
- **Version Information:** Build-time version, commit, and branch information are injected into the binary using `ldflags`.

## Contributing
Contributions are welcome! Please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project. (Note: `CONTRIBUTING.md` is a placeholder and needs to be created if not present).

## License
This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.