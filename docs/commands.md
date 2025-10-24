# Govman Commands

This document provides a detailed reference for all `govman` CLI commands.

## Core Concepts

-   **Version String**: A Go version, such as `1.25.1`, `1.22`, or `latest`.
-   **Activation**: The process of making a specific Go version available in the shell's `PATH`.
-   **Scope**: Activation can be `session-only` (temporary), `system-default` (persistent), or `project-local` (tied to a directory).

---

## `govman install`

Downloads and installs one or more Go versions.

### Usage

```bash
govman install [version...]
```

### Arguments

-   `version...`: One or more version strings to install. `latest` is a special keyword for the most recent stable version.

### Features

-   **Parallel Downloads**: Downloads multiple versions concurrently.
-   **Resumable**: Automatically resumes interrupted downloads.
-   **Checksum Validation**: Ensures the integrity of downloaded files.
-   **Caching**: Avoids re-downloading archives that are already present in the cache.

### Examples

```bash
# Install the latest stable Go version
govman install latest

# Install a specific version
govman install 1.25.1

# Install multiple versions at once
govman install 1.25.1 1.20.12

# Install a pre-release version
govman install 1.22rc1
```

---

## `govman uninstall`

Removes an installed Go version.

### Usage

```bash
govman uninstall <version>
```

### Aliases

-   `remove`
-   `rm`

### Safety

-   Prevents removal of the currently active Go version.
-   Performs a complete cleanup of the version's installation directory.

### Example

```bash
govman uninstall 1.20.12
```

---

## `govman use`

Switches the active Go version.

### Usage

```bash
govman use <version> [flags]
```

### Activation Modes (Flags)

-   `--default` or `-d`: Sets the version as the system-wide default for all new shell sessions.
-   `--local` or `-l`: Sets the version for the current project by creating a `.govman-version` file.
-   (no flag): Activates the version for the current shell session only.

### Examples

```bash
# Activate for the current session
govman use 1.25.1

# Set as the system default
govman use 1.25.1 --default

# Set for the current project
govman use 1.22.4 --local

# Switch to the default version
govman use default
```

---

## `govman list`

Lists installed or remote Go versions.

### Usage

```bash
govman list [flags]
```

### Aliases

-   `ls`

### Flags

-   `--remote` or `-r`: Lists all available versions from Go's official release source.
-   `--stable-only`: (Remote only) Shows only stable, production-ready versions.
-   `--beta`: (Remote only) Includes beta/rc versions.
-   `--pattern <glob>`: (Remote only) Filters remote versions using a glob pattern (e.g., `1.25.*`).

### Examples

```bash
# List installed versions
govman list

# List available remote versions
govman list --remote

# Find all 1.25 patch releases
govman list --remote --pattern "1.25.*"

# List only stable versions
govman list --remote --stable-only

# Include beta/rc versions
govman list --remote --beta
```

---

## `govman current`

Displays detailed information about the currently active Go version.

### Usage

```bash
govman current
```

### Information Displayed

-   Version number and release status
-   Installation path and size
-   Platform (OS/Arch)
-   Installation date and source
-   Activation method (system, project, or session)

---

## `govman info`

Displays detailed information about a specific *installed* Go version.

### Usage

```bash
govman info <version>
```

### Information Displayed

-   Version number and release details
-   Complete installation path and directory structure
-   Platform architecture and OS compatibility
-   Installation date, size, and disk usage
-   Binary locations and environment details
-   Active status (whether currently in use)
-   Age warnings for versions older than 6 months

---

## `govman clean`

Removes cached downloads and temporary data to free up disk space.

### What Gets Cleaned

-   Downloaded Go archive files (.tar.gz, .zip)
-   Temporary extraction directories
-   Incomplete or corrupted downloads
-   Obsolete cache metadata and checksums

### Safety

-   This command is safe to run at any time.
-   It **does not** remove any installed Go versions.
-   Your project files and configurations are preserved.
-   Only temporary cache files are removed.

---

## `govman init`

Sets up shell integration for automatic version switching.

### Usage

```bash
govman init [flags]
```

### Flags

-   `--force` or `-f`: Overwrites any existing `govman` configuration in your shell profile.
-   `--shell <name>`: Manually specifies the shell (e.g., `bash`, `zsh`, `fish`, `powershell`).

### Supported Shells

-   **Bash** (.bashrc, .bash_profile)
-   **Zsh** (.zshrc)
-   **Fish** (config.fish)
-   **PowerShell** (profile)

### Functionality

-   Detects your shell and modifies the appropriate configuration file.
-   Adds `~/.govman/bin` to your PATH.
-   Sets up automatic hooks for version switching based on `.govman-version` files.
-   Creates wrapper functions for seamless integration.

---

## `govman selfupdate`

Updates `govman` to the latest version.

### Usage

```bash
govman selfupdate [flags]
```

### Flags

-   `--check`: Checks for a new version without installing it.
-   `--force`: Reinstalls the latest version even if you are already up-to-date.
-   `--prerelease`: Includes pre-release versions in the update check.

### Features

-   Automatic platform detection and binary selection
-   Safe backup and rollback on failure
-   Integrity verification and secure downloads
-   Detailed release notes and changelog display

---

## `govman refresh`

Manually triggers the auto-switching mechanism in the current directory.

### Usage

```bash
govman refresh
```

### Behavior

-   Re-evaluates the current directory for `.govman-version` files
-   Switch to the appropriate version (local or default)
-   Useful after adding/removing `.govman-version` files
-   Equivalent to the auto-switch that happens on `cd`

### Use Case

Useful if you manually create or modify a `.govman-version` file and want to immediately switch to the specified version without changing directories.