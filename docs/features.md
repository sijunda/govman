# GOVMAN Features

GOVMAN is a powerful and cross-platform Go version manager that simplifies the installation, management, and switching between multiple Go versions. Here's a comprehensive list of its features:

## Core Features

### 1. Version Management
- **Install Go versions**: Easily install any Go version including stable releases, beta, and release candidates
- **Uninstall Go versions**: Remove installed Go versions to free up disk space
- **Switch between versions**: Quickly switch between different Go versions for different projects
- **List versions**: View all installed Go versions or browse available versions for download

### 2. Flexible Version Selection
- **Project-specific versions**: Set Go versions on a per-project basis using `.go-version` files
- **System-wide defaults**: Set a default Go version for your entire system
- **Session-only switching**: Temporarily switch to a Go version for the current session

### 3. Cross-Platform Support
- **Multi-OS compatibility**: Works seamlessly on Windows, macOS, and Linux
- **Architecture support**: Supports multiple architectures including ARM
- **Shell integration**: Automatic integration with bash, zsh, fish, and PowerShell

### 4. Smart Installation System
- **Parallel downloads**: Utilize fast parallel downloads for quicker installations
- **Resume capability**: Automatically resume failed downloads
- **Offline mode**: Intelligent caching allows working offline with previously downloaded versions
- **Zero configuration**: Works out of the box with no setup required

### 5. Performance & Efficiency
- **Lightning-fast operations**: Optimized for speed in all operations
- **Intelligent caching**: Cache downloaded versions to avoid re-downloading
- **Disk space management**: Built-in cleanup tools to manage disk space efficiently
- **No admin/sudo required**: Fully userspace installation for enhanced security

## Command Line Interface

### Installation Commands
- `govman install [version...]` - Install one or more Go versions
- `govman uninstall <version>` - Uninstall a Go version
- `govman clean` - Clean download cache to free up disk space

### Version Management Commands
- `govman use <version>` - Switch to a specific Go version
- `govman current` - Show the currently active Go version
- `govman list` - List installed Go versions
- `govman info <version>` - Show detailed information about an installed Go version

### Version Discovery
- `govman list --remote` - List available Go versions for download
- `govman list --remote --pattern <pattern>` - Filter remote versions by pattern
- `govman list --remote --beta` - Include beta/rc versions in remote listing

### System Management
- `govman init` - Initialize shell integration
- `govman selfupdate` - Update govman to the latest version
- `govman --help` - Display help information

## Advanced Features

### 1. Shell Integration
- **Automatic PATH management**: Seamlessly manages your PATH environment variable
- **Auto-switching**: Automatically switches Go versions based on project configuration
- **Multi-shell support**: Works with bash, zsh, fish, and PowerShell

### 2. Configuration Options
- **Flexible configuration**: Uses viper for configuration management
- **Custom config file**: Support for custom configuration files with `--config` flag
- **Verbose/quiet modes**: Control output verbosity with `--verbose` and `--quiet` flags

### 3. Project-Level Configuration
- **Local version files**: Use `.go-version` files to specify project-specific Go versions
- **Directory-based switching**: Automatically use the correct Go version when entering project directories

### 4. Self-Management
- **Self-updating**: Built-in self-update mechanism to keep govman up to date
- **Prerelease support**: Option to update to prerelease versions
- **Force updates**: Force reinstallation even if already at the latest version

## User Experience

### 1. Informative Output
- **Colored output**: Color-coded messages for better readability
- **Emoji support**: Emojis for quick visual recognition of message types
- **Progress indication**: Clear progress indicators during long operations
- **Helpful error messages**: Contextual error messages with suggestions

### 2. Multiple Verbosity Levels
- **Quiet mode**: Show only errors (`--quiet`)
- **Normal mode**: Show essential information (default)
- **Verbose mode**: Show detailed information including debug messages (`--verbose`)

### 3. Safety Features
- **Backup protection**: Backup current binary during self-updates
- **Permission checking**: Clear error messages for permission-related issues
- **Validation**: Checksum verification for downloaded binaries

## Technical Features

### 1. Robust Architecture
- **Modular design**: Well-organized codebase with clear separation of concerns
- **Cobra CLI framework**: Built on the robust Cobra command-line library
- **Viper configuration**: Flexible configuration management

### 2. Development Tools
- **Makefile automation**: Comprehensive Makefile for building, testing, and development
- **Testing support**: Built-in support for unit, integration, and benchmark tests
- **Code quality tools**: Integrated linting and formatting tools

### 3. Build System
- **Cross-compilation**: Support for building on multiple platforms
- **Version injection**: Build-time version information injection
- **Release management**: Automated release processes

## Installation Flexibility

### 1. Multiple Installation Methods
- **Script-based installation**: One-line installation script for Unix-like systems
- **Manual installation**: Direct binary installation
- **Source compilation**: Build from source code

### 2. Flexible Installation Paths
- **User directory installation**: Install to user directory without admin privileges
- **System-wide installation**: Install to system directories with sudo
- **Custom paths**: Support for custom installation paths

This comprehensive feature set makes GOVMAN one of the most powerful and user-friendly Go version managers available, suitable for both individual developers and enterprise environments.