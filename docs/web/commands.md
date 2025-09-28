# 📖 Command Reference

Complete API documentation for all GOVMAN commands, options, and usage patterns.

## 📋 Table of Contents

- [Global Options](#-global-options)
- [Installation Commands](#-installation-commands)
- [Version Management](#-version-management)
- [Information Commands](#-information-commands)
- [Maintenance Commands](#-maintenance-commands)
- [Configuration Commands](#-configuration-commands)
- [Exit Codes](#-exit-codes)
- [Environment Variables](#-environment-variables)

## 🌐 Global Options

These options are available for all commands:

```bash
govman [global options] command [command options] [arguments...]
```

### **Global Flags**

| Flag | Description | Example |
|------|-------------|---------|
| `--config string` | Config file path | `govman --config /path/to/config.yaml list` |
| `--verbose` | Enable verbose output | `govman --verbose install 1.21.1` |
| `--quiet` | Suppress non-error output | `govman --quiet install 1.21.1` |
| `--help, -h` | Show help | `govman --help` |
| `--version, -v` | Show version | `govman --version` |

## 📦 Installation Commands

### **`govman install`**

Install one or more Go versions.

#### **Syntax**
```bash
govman install [version...] [flags]
```

#### **Examples**
```bash
# Install latest stable version
govman install latest

# Install specific version
govman install 1.21.1

# Install multiple versions
govman install 1.19.5 1.20.5 1.21.1

# Install latest patch version of 1.20
govman install 1.20
```

#### **Supported Version Formats**
- `latest` - Latest stable release
- `1.21.1` - Exact version
- `1.20` - Latest patch version (e.g., 1.20.5)
- Pre-release versions (e.g., `1.22rc1`, `1.22beta1`)

#### **Behavior**
- Downloads from official Go website
- Verifies checksums automatically
- Caches downloads for faster reinstallation
- Supports parallel downloads
- Resumes interrupted downloads

### **`govman uninstall`**

Remove installed Go versions.

#### **Syntax**
```bash
govman uninstall <version> [flags]
```

#### **Examples**
```bash
# Uninstall specific version
govman uninstall 1.20.5

# Uninstall multiple versions
govman uninstall 1.19.5 1.20.5
```

#### **Safety Features**
- Cannot uninstall currently active version
- Confirms before deletion
- Preserves user data and cache

## 🔄 Version Management

### **`govman use`**

Switch to a specific Go version with different activation modes.

#### **Syntax**
```bash
govman use <version> [flags]
```

#### **Activation Modes**

| Mode | Flag | Scope | Persistence | Example |
|------|------|-------|-------------|---------|
| **Session-only** | *(default)* | Current session | Temporary | `govman use 1.21.1` |
| **System default** | `--default, -d` | All new sessions | Permanent | `govman use 1.21.1 --default` |
| **Project local** | `--local, -l` | Current directory | Project-specific | `govman use 1.21.1 --local` |

#### **Examples**
```bash
# Session-only (temporary)
govman use 1.21.1

# Set as system default (permanent)
govman use 1.21.1 --default

# Set for current project (creates .govman-version)
govman use 1.21.1 --local

# Use default version (whatever is set as default)
govman use default
```

#### **Behavior**
- **Session-only**: Updates PATH for current terminal session
- **System default**: Creates symlink, affects all new terminals
- **Project local**: Creates `.govman-version` file, auto-switches when entering directory
- **Immediate effect**: Both `govman current` and `go version` show the same version instantly

### **`govman refresh`** ⭐ *NEW*

Manually refresh the current Go version based on directory context.

#### **Syntax**
```bash
govman refresh
```

#### **Examples**
```bash
# After removing .govman-version file
rm .govman-version
govman refresh  # Switches back to default version

# After manually editing .govman-version
echo "1.20.5" > .govman-version
govman refresh  # Switches to 1.20.5
```

#### **Use Cases**
- After manually editing `.govman-version` files
- When auto-switching doesn't trigger
- For manual control in edge cases
- In environments without shell integration

## 📊 Information Commands

### **`govman list`**

List Go versions (installed or available).

#### **Syntax**
```bash
govman list [flags]
```

#### **Flags**
| Flag | Description | Example |
|------|-------------|---------|
| `--remote` | Show available versions for download | `govman list --remote` |
| *(default)* | Show installed versions | `govman list` |

#### **Examples**
```bash
# List installed versions
govman list

# List available versions for download
govman list --remote

# List with additional details
govman --verbose list
```

#### **Output Format**
```bash
# Example output
ℹ️  📋 Installed Go Versions (3 total):
ℹ️  ────────────────────────────────────────────────────────────
ℹ️  → ✅ 1.21.1 [default]          193.35 MB   installed: 2024-09-22
ℹ️    💾 1.20.5                    203.29 MB   installed: 2024-09-22
ℹ️    💾 1.19.5                    321.92 MB   installed: 2024-09-22
ℹ️  ────────────────────────────────────────────────────────────
ℹ️  📊 Total disk usage: 718.56 MB across 3 versions
ℹ️  ✅ Currently active: Go 1.21.1
```

### **`govman current`**

Show currently active Go version and detailed information.

#### **Syntax**
```bash
govman current
```

#### **Examples**
```bash
# Show current version
govman current

# Detailed current version info
govman --verbose current
```

#### **Output Format**
```bash
# Example output
ℹ️  🔍 Current Go Environment:
ℹ️  ──────────────────────────────────────────────────
ℹ️  ✅ Version:        Go 1.21.1
ℹ️  📁 Install Path:    /Users/user/.govman/versions/go1.21.1
ℹ️  🖥️ Platform:        darwin/arm64
ℹ️  📅 Installed:       2024-09-22 20:16:59 WIB
ℹ️  💾 Disk Usage:      193.35 MB
ℹ️  🔄 Activation:      📱 Session-only (temporary)
ℹ️  ──────────────────────────────────────────────────
ℹ️  💡 Run 'go version' to verify your Go installation
```

#### **Activation Types**
- `📱 Session-only (temporary)` - Active for current session only
- `🏠 System default (persistent)` - Set as system-wide default
- `📁 Project local (.govman-version)` - Project-specific version

### **`govman info`**

Show detailed information about a specific Go version.

#### **Syntax**
```bash
govman info <version>
```

#### **Examples**
```bash
# Show info for specific version
govman info 1.21.1

# Show info for currently active version
govman info $(govman current --quiet)
```

#### **Output Includes**
- Installation date and size
- Platform and architecture
- Installation path
- Download URL and checksum
- Release notes link

## 🧹 Maintenance Commands

### **`govman clean`**

Clean cache and temporary files to reclaim disk space.

#### **Syntax**
```bash
govman clean [flags]
```

#### **Examples**
```bash
# Clean download cache
govman clean

# Clean with verbose output
govman --verbose clean
```

#### **What Gets Cleaned**
- Downloaded installation archives
- Temporary extraction files
- Cached version metadata
- Orphaned symlinks

#### **What's Preserved**
- Installed Go versions
- Configuration files
- Project `.govman-version` files

### **`govman selfupdate`**

Update GOVMAN to the latest version.

#### **Syntax**
```bash
govman selfupdate [flags]
```

#### **Examples**
```bash
# Update to latest version
govman selfupdate

# Check for updates without installing
govman selfupdate --check
```

#### **Behavior**
- Downloads latest release from GitHub
- Verifies download integrity
- Preserves existing configuration
- Updates shell integration if needed

## ⚙️ Configuration Commands

### **`govman init`**

Initialize shell integration for automatic version switching.

#### **Syntax**
```bash
govman init [flags]
```

#### **Flags**
| Flag | Description | Example |
|------|-------------|---------|
| `--force` | Overwrite existing configuration | `govman init --force` |
| `--shell string` | Specify shell type | `govman init --shell zsh` |

#### **Examples**
```bash
# Auto-detect shell and initialize
govman init

# Force reinitialize (overwrite existing)
govman init --force

# Initialize for specific shell
govman init --shell bash
```

#### **What It Does**
- Detects your shell automatically
- Adds shell integration to config files
- Enables automatic version switching
- Sets up PATH management
- Creates shell wrapper functions

#### **Supported Shells**
- **Bash** (`~/.bashrc`, `~/.bash_profile`)
- **Zsh** (`~/.zshrc`)
- **Fish** (`~/.config/fish/config.fish`)
- **PowerShell** (`$PROFILE`)
- **Command Prompt** (creates wrapper batch file)

## 🔧 Advanced Usage

### **Configuration File**

GOVMAN can be configured via `~/.govman/config.yaml`:

```yaml
# Example configuration
install_dir: ~/.govman/versions
cache_dir: ~/.govman/cache
default_version: "1.21.1"
verbose: false
quiet: false

download:
  parallel: true
  max_connections: 4
  timeout: 300s

auto_switch:
  enabled: true
  project_file: .govman-version

shell:
  auto_detect: true
  completion: true
```

### **Project Configuration**

Create `.govman-version` in your project root:

```bash
# Method 1: Using govman
govman use 1.21.1 --local

# Method 2: Manual creation
echo "1.21.1" > .govman-version

# Method 3: Team setup
git add .govman-version
git commit -m "Set Go version for project"
```

### **Environment Variables**

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `GOVMAN_ROOT` | Installation directory | `~/.govman` | `export GOVMAN_ROOT=/opt/govman` |
| `GOVMAN_CONFIG` | Config file path | `~/.govman/config.yaml` | `export GOVMAN_CONFIG=/etc/govman.yaml` |
| `GOVMAN_CACHE_DIR` | Cache directory | `~/.govman/cache` | `export GOVMAN_CACHE_DIR=/tmp/govman` |
| `GOVMAN_MIRROR_URL` | Download mirror | `https://go.dev/dl/` | `export GOVMAN_MIRROR_URL=https://golang.google.cn/dl/` |

## 📊 Exit Codes

GOVMAN uses standard exit codes:

| Code | Meaning | Example Scenario |
|------|---------|------------------|
| `0` | Success | Command completed successfully |
| `1` | General error | Invalid command or option |
| `2` | Misuse | Wrong number of arguments |
| `126` | Permission denied | Cannot write to installation directory |
| `127` | Command not found | Specified Go version not available |
| `130` | Interrupted | User pressed Ctrl+C |

### **Error Handling Examples**
```bash
# Check if command succeeded
if govman install 1.21.1; then
    echo "Installation successful"
else
    echo "Installation failed with exit code $?"
fi

# Use in scripts
govman use 1.21.1 --local || {
    echo "Failed to set local version"
    exit 1
}
```

## 🔄 Shell Integration Details

### **Auto-Switching Behavior**

When shell integration is enabled:

```bash
# Directory change triggers version check
cd my-project          # Checks for .govman-version
                       # Auto-switches if file exists

# Manual commands work immediately
govman use 1.20.5 --local   # Creates file + switches immediately
govman refresh               # Re-evaluates current directory
```

### **Shell-Specific Features**

#### **Bash/Zsh**
- `cd` command override for auto-switching
- Production-safe wrapper functions
- Error handling and recovery

#### **Fish**
- PWD change hooks for auto-switching
- Native Fish function syntax
- Fish-specific error handling

#### **PowerShell**
- `Set-Location` override for auto-switching
- Native PowerShell error handling
- Windows-specific PATH management

#### **Command Prompt**
- Batch file wrapper for basic functionality
- Manual refresh required
- Limited auto-switching

## 💡 Tips and Best Practices

### **Version Selection**
```bash
# Use specific versions for reproducibility
govman install 1.21.1    # Not just "latest"

# Pin project versions
echo "1.21.1" > .govman-version
git add .govman-version

# Test compatibility
govman use 1.20.5 --local  # Test with older version
```

### **Team Workflows**
```bash
# Team leader sets version
govman use 1.21.1 --local
git add .govman-version
git commit -m "Set team Go version"

# Team members get automatic version
git pull
cd project  # Automatically switches to 1.21.1
```

### **CI/CD Integration**
```bash
# Install specific version in CI
govman install 1.21.1
govman use 1.21.1

# Use project version
if [ -f .govman-version ]; then
    govman use $(cat .govman-version)
fi
```

---

**Need more help?** Check our [Examples](examples.md) or [FAQ](faq.md) for common usage patterns.