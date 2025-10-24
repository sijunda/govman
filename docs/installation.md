# Installation

`govman` can be installed on Windows, macOS, and Linux.

## Quick Install (Recommended)

The quickest way to install `govman` is by using the official installation scripts.

### macOS / Linux

Run the following command in your terminal:

```bash
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

This script will:
1.  Detect your OS and architecture.
2.  Download the latest `govman` binary to `~/.govman/bin`.
3.  Add `~/.govman/bin` to your shell's `PATH` by modifying your profile file (`.bashrc`, `.zshrc`, etc.).
4.  Run `govman init` to set up shell integration.

### Windows (PowerShell)

Run the following command in PowerShell:

```powershell
irm https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.ps1 | iex
```

This script will:
1.  Download the latest `govman` binary to `%USERPROFILE%\\.govman\\bin`.
2.  Add the directory to your user `PATH` environment variable.
3.  Run `govman init` to configure your PowerShell profile for auto-switching.

## Manual Installation

If you prefer to install `govman` manually:

1.  **Download the Binary**:
    Go to the [GitHub Releases page](https://github.com/sijunda/govman/releases) and download the appropriate binary for your operating system and architecture.

2.  **Place it in your PATH**:
    Move the downloaded binary to a directory that is included in your system's `PATH`. A common choice is `/usr/local/bin` on Linux/macOS or a custom scripts folder on Windows.

    **macOS / Linux:**
    ```bash
    # Example:
    mv ./govman-darwin-arm64 /usr/local/bin/govman
    chmod +x /usr/local/bin/govman
    ```

3.  **Initialize the Shell**:
    Run `govman init` to configure your shell for automatic version switching. This step is crucial for the best experience.
    ```bash
    govman init
    ```
    Follow the on-screen instructions and restart your terminal.
    
    > **Note**: For detailed information about shell integration, including manual setup instructions for specific shells (Bash, Zsh, Fish, PowerShell), see the [Shell Integration](shell-integration.md) documentation.

## Build from Source

If you have Go installed, you can build `govman` from source.

```bash
# 1. Clone the repository
git clone https://github.com/sijunda/govman.git
cd govman

# 2. Build the binary
go build -o govman ./cmd/govman

# 3. Move the binary to your PATH
mv ./govman /usr/local/bin/

# 4. Initialize your shell
govman init
```

## Uninstallation

To uninstall `govman`, use the corresponding uninstallation script.

### macOS / Linux
```bash
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/uninstall.sh | bash
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/sijunda/govman/main/scripts/uninstall.ps1 | iex