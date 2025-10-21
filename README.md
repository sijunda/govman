# GOVMAN - Go Version Manager

<p align="center">
  <img src="https://img.shields.io/github/go-mod/go-version/sijunda/govman" alt="Go Version">
  <img src="https://img.shields.io/github/license/sijunda/govman" alt="License">
  <img src="https://img.shields.io/github/v/release/sijunda/govman" alt="Release">
  <img src="https://img.shields.io/github/downloads/sijunda/govman/total" alt="Downloads">
</p>

<img src="./govman.png" alt="Govman">

**The simplest, fastest, and most reliable Go version manager.**

GOVMAN is a lightweight, zero-dependency Go version manager that makes it effortless to install, manage, and switch between different Go versions. Built for developers who need reliability and performance without complexity.

## ✨ **Key Features**

⚡ **Lightning-fast installation and switching** between Go versions
🎯 **Zero configuration** - works out of the box, no setup required
📁 **Project-specific versions** with `.govman-version` file support
🚫 **No admin/sudo required** - fully userspace installation
💾 **Intelligent caching** with offline mode support
📦 **Parallel downloads** with automatic resume on failure
🌍 **Cross-platform support** (Windows, macOS, Linux, ARM)
🧹 **Built-in cleanup tools** to manage disk space efficiently

## 📦 **Installation**

### Quick Install

**macOS/Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.ps1 | iex
```

### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/sijunda/govman/releases)
2. Extract the binary to a directory in your PATH
3. Run `govman init` to set up shell integration

### Build from Source

```bash
git clone https://github.com/sijunda/govman.git
cd govman
go build -o govman ./cmd/govman
./govman init
```

## 🚀 **Quick Start**

```bash
# Install the latest Go version
govman install latest

# Install a specific version
govman install 1.21.5

# List all available versions
govman list --remote

# Switch to a specific version
govman use 1.21.5

# Set project-specific version
echo "1.21.5" > .govman-version
govman use  # Automatically uses project version

# Check current version
govman current

# Clean up cache and unused versions
govman clean
```

## 📚 **Commands**

### Installation & Management
```bash
govman install <version>         # Install a Go version
govman install latest            # Install latest stable version
govman uninstall <version>       # Remove an installed version
govman list                      # List installed versions
govman list --remote             # List all available versions
govman clean                     # Clean cache and temporary files
```

### Version Switching
```bash
govman use <version>             # Switch to a version (session-only)
govman use <version> --default   # Set as system default
govman use <version> --local     # Set for current project
govman current                   # Show active version and method
```

### Information & Utilities
```bash
govman info <version>            # Show version details and disk usage
govman refresh                   # Refresh version cache
govman selfupdate                # Update govman itself
govman init                      # Set up shell integration
```

## ⚙️ **Configuration**

GOVMAN uses a YAML configuration file located at `~/.govman/config.yaml`:

```yaml
# Installation directory for Go versions
install_dir: "~/.govman/versions"

# Cache directory for downloads
cache_dir: "~/.govman/cache"

# Default Go version (empty = none)
default_version: ""

# Download configuration
download:
  parallel: true          # Enable parallel downloads
  max_connections: 4      # Maximum concurrent connections
  timeout: 300s           # Download timeout
  retry_count: 3          # Number of retry attempts
  retry_delay: 5s         # Delay between retries

# Mirror configuration (for China users)
mirror:
  enabled: false
  url: "https://golang.google.cn/dl/"

# Auto-switch configuration
auto_switch:
  enabled: true
  project_file: ".govman-version"

# Shell integration
shell:
  auto_detect: true
  completion: true

# Go releases API
go_releases:
  api_url: "https://go.dev/dl/?mode=json&include=all"
  download_url: "https://go.dev/dl/%s"
  cache_expiry: 10m

# Self-update configuration
self_update:
  github_api_url: "https://api.github.com/repos/sijunda/govman/releases/latest"
  github_releases_url: "https://api.github.com/repos/sijunda/govman/releases?per_page=1"

# Logging
quiet: false    # Suppress normal output
verbose: false  # Enable verbose logging
```

## 🏗️ **Project Structure**

```
govman/
├── cmd/govman/            # Main application entry point
│   └── main.go
├── internal/              # Internal packages
│   ├── cli/               # Command-line interface
│   │   ├── cli.go         # Root command and banner
│   │   ├── command.go     # Command registration
│   │   ├── clean.go       # Cache cleanup
│   │   ├── current.go     # Current version info
│   │   ├── info.go        # Version information
│   │   ├── init.go        # Shell initialization
│   │   ├── install.go     # Install/uninstall commands
│   │   ├── list.go        # List versions
│   │   ├── refresh.go     # Refresh cache
│   │   ├── selfupdate.go  # Self-update functionality
│   │   └── use.go         # Version switching
│   ├── config/            # Configuration management
│   │   └── config.go      # Config loading and validation
│   ├── downloader/        # Download engine
│   │   └── downloader.go  # Parallel downloads with resume
│   ├── golang/            # Go releases API client
│   │   └── releases.go    # Version parsing and fetching
│   ├── logger/            # Structured logging
│   │   └── logger.go      # Multi-level logging system
│   ├── manager/           # Core version management
│   │   └── manager.go     # Install/uninstall/switch logic
│   ├── progress/          # Progress visualization
│   │   └── progress.go    # Download progress bars
│   ├── shell/             # Shell integration
│   │   └── shell.go       # Multi-shell support
│   ├── symlink/           # Symlink utilities
│   │   └── symlink.go     # Cross-platform symlinks
│   ├── util/              # Utility functions
│   │   └── format.go      # String formatting helpers
│   └── version/           # Version information
│       └── version.go     # Build version management
├── scripts/               # Installation scripts
│   ├── install.sh         # Unix installation
│   ├── install.ps1        # Windows PowerShell
│   ├── install.bat        # Windows batch
│   ├── uninstall.sh       # Unix uninstall
│   ├── uninstall.ps1      # Windows uninstall
│   └── uninstall.bat      # Windows batch uninstall
├── Dockerfile             # Container support
├── Makefile               # Build automation
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
└── config.yaml.example    # Example configuration
```

## 🔧 **Shell Integration**

GOVMAN supports automatic Go version switching when entering directories with a `.govman-version` file.

### Automatic Setup
```bash
govman init  # Automatically detects and configures your shell
```

### Supported Shells
- **Bash/Zsh**: Full auto-switching support
- **Fish**: Full auto-switching support
- **PowerShell**: Full auto-switching support
- **Command Prompt**: Limited auto-switching support

For manual setup instructions, see the [Shell Integration Guide](docs/shell-integration.md).

## 🧪 **Testing**

GOVMAN includes comprehensive tests for all components:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/logger
go test ./internal/downloader
go test ./internal/manager

# Run tests with verbose output
go test -v ./...

# Benchmark tests
go test -bench=. ./...
```

### Test Coverage
- **Logger**: Complete test coverage for all log levels and concurrency
- **Downloader**: Tests for parallel downloads, resume, and error handling
- **Manager**: Version management, installation, and switching logic
- **Config**: Path expansion, validation, and default handling
- **Utils**: String formatting and utility functions
- **Version**: Version comparison and parsing algorithms

## 🔨 **Development**

### Prerequisites
- Go 1.25 or later
- Make (optional, for using Makefile)

### Building
```bash
# Build for current platform
go build -o govman ./cmd/govman

# Build for all platforms using Makefile
make build-all

# Run development version
go run ./cmd/govman --help
```

### Contributing

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/new-feature`
3. **Commit** your changes: `git commit -m 'Add new feature'`
4. **Push** to the branch: `git push origin feature/new-feature`
5. **Open** a Pull Request

### Code Standards
- Follow standard Go formatting (`go fmt`)
- Write comprehensive tests for new features
- Update documentation for user-facing changes
- Use conventional commit messages
- Ensure all tests pass before submitting

### Architecture
GOVMAN follows a clean architecture pattern:
- **CLI Layer**: Command parsing and user interaction
- **Manager Layer**: Core business logic for version management
- **Downloader Layer**: Handles downloads with resume capability
- **Config Layer**: Configuration management and validation
- **Shell Layer**: Multi-shell integration and auto-switching
- **Util Layer**: Shared utilities and helpers

## 🌍 **Supported Platforms**

| Platform | Architecture | Status |
|----------|--------------|--------|
| Linux    | amd64        | ✅ Fully Supported |
| Linux    | arm64        | ✅ Fully Supported |
| macOS    | amd64        | ✅ Fully Supported |
| macOS    | arm64 (M1/M2)| ✅ Fully Supported |
| Windows  | amd64        | ✅ Fully Supported |
| Windows  | arm64        | ✅ Fully Supported |
| FreeBSD  | amd64        | ✅ Fully Supported |

## 🚀 **Performance**

GOVMAN is designed for performance:
- **Parallel Downloads**: Up to 4 concurrent connections
- **Resume Support**: Interrupted downloads resume automatically
- **Smart Caching**: Avoids re-downloading existing files
- **Fast Switching**: Symlink-based version switching in milliseconds
- **Memory Efficient**: Minimal memory footprint with optimized data structures
- **Background Processing**: Non-blocking operations where possible

## 🛡️ **Security**

- **Path Traversal Protection**: Prevents malicious archive extraction
- **Checksum Verification**: SHA-256 validation for all downloads
- **Secure Downloads**: HTTPS-only with certificate validation
- **Sandboxed Extraction**: Safe archive handling with path validation
- **No Elevated Privileges**: Runs entirely in userspace

## 🔍 **Troubleshooting**

### Common Issues

**Permission Denied**
```bash
# Ensure govman binary is executable
chmod +x ~/.govman/bin/govman

# Check if ~/.govman/bin is in PATH
echo $PATH | grep -q "$HOME/.govman/bin" || echo "PATH not set correctly"
```

**Go Version Not Found**
```bash
# Refresh version cache
govman refresh

# Check if version exists
govman list --remote | grep <version>
```

**Shell Integration Not Working**
```bash
# Re-run initialization
govman init --force

# Manually source shell configuration
source ~/.bashrc  # or ~/.zshrc, ~/.config/fish/config.fish
```

### Debug Mode
```bash
# Enable verbose logging
govman --verbose <command>

# Check current configuration
govman current --verbose
```

## 📈 **Roadmap**

- [ ] GUI application for desktop users
- [ ] Plugin system for custom integrations
- [ ] Docker container with multiple Go versions
- [ ] Integration with popular IDEs
- [ ] Automatic security update notifications
- [ ] Custom mirror support for enterprise users
- [ ] Version constraint resolution (go.mod integration)


## 📄 **Changelog**

See [CHANGELOG.md](CHANGELOG.md) for a detailed history of changes.

## 🤝 **Contributing**

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📝 **License**

This project is licensed under the **MIT License**. See the [LICENSE.md](LICENSE.md) file for details.

---

<p align="center">
  <sub>Built with ❤️ by Muhammad Jundana Al Basyir</sub>
</p>