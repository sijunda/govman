# Quick Start

This guide will get you up and running with `govman` in minutes.

## 1. Install `govman`

You can install `govman` with a single command.

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.ps1 | iex
```

After installation, restart your terminal or run `source ~/.bashrc` (or the equivalent for your shell) to update your `PATH`.

## 2. Install the Latest Go Version

Use the `install` command to download and set up the latest stable version of Go.

```bash
govman install latest
```

## 3. Activate the New Version

Make the newly installed version available in your current shell session.

```bash
govman use latest
```
To make it the default for all future sessions, use the `--default` flag:
```bash
govman use latest --default
```

## 4. Verify Your Installation

Check that the correct Go version is now active.

```bash
go version
```

You should see the version you just installed.

## 5. Project-Specific Versions

`govman` can automatically switch Go versions based on your project's needs.

1.  Navigate to your project directory:
    ```bash
    cd /path/to/your/project
    ```

2.  Create a `.govman-version` file:
    ```bash
    echo "1.22.4" > .govman-version
    ```
    (First, make sure version `1.22.4` is installed with `govman install 1.22.4`)

3.  `govman` will now automatically use Go `1.22.4` whenever you are in this directory or its subdirectories.

That's it! You're now ready to manage your Go versions like a pro.