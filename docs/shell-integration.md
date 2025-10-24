# Shell Integration

`govman` features powerful shell integration that enables automatic Go version switching as you navigate your filesystem.

## How It Works

When you run `govman init`, the tool automatically detects your shell (e.g., Bash, Zsh, PowerShell) and adds a small script to your shell's configuration file (e.g., `.zshrc`, `.bash_profile`, `profile.ps1`).

This script does two things:
1.  Adds the `~/.govman/bin` directory, your `GOBIN` (if set), your `GOPATH/bin` (if Go is available), and the default `$HOME/go/bin` to your `PATH`.
2.  Hooks into your shell's `cd` (change directory) command or prompt.

When you `cd` into a directory that contains a `.govman-version` file, the hook is triggered, and it automatically runs `govman use` to activate the version specified in that file. The activation is for the current session only, so it doesn't change your system-wide default.

## Supported Shells

`govman` provides first-class support for the most popular shells:

| Shell        | Automatic Setup | Auto-Switching | Notes                                  |
|--------------|-----------------|----------------|----------------------------------------|
| **Zsh**      | ✅ Yes          | ✅ Yes         | Recommended for macOS and Linux users. |
| **Bash**     | ✅ Yes          | ✅ Yes         |                                        |
| **Fish**     | ✅ Yes          | ✅ Yes         |                                        |
| **PowerShell** | ✅ Yes          | ✅ Yes         | Recommended for Windows users.         |
| **Cmd.exe**  | ⚠️ Limited     | ❌ No          | Not recommended. Lacks hooking support.|

## Setup

### Automatic Setup (Recommended)

The easiest way to set up shell integration is to run:

```bash
govman init
```

This command will guide you through the process. You will need to restart your shell for the changes to take effect.

### Manual Setup

If the automatic setup fails, or if you prefer to manage your shell configuration manually, you can add the required scripts yourself.

Run `govman init --shell <your-shell-name>` to see the exact lines you need to add to your configuration file.

For example, for Zsh:
```bash
govman init --shell zsh
```
This will output the script block that you can copy and paste into your `~/.zshrc` file.

### Bash

Add the following to your `~/.bashrc`, `~/.bash_profile`, or `~/.profile`:

```bash
# GOVMAN - Go Version Manager
export PATH="$HOME/.govman/bin:$PATH"
export GOTOOLCHAIN=local

# Go tool binaries on PATH
if [[ -n "$GOBIN" ]]; then
    export PATH="$GOBIN:$PATH"
fi
if command -v go >/dev/null 2>&1; then
    export PATH="$(go env GOPATH)/bin:$PATH"
fi
# Default GOPATH location
export PATH="$HOME/go/bin:$PATH"

# Wrapper function for automatic PATH execution
govman() {
    local govman_bin="$HOME/.govman/bin/govman"
    if [[ "$1" == "use" && "$#" -ge 2 && "$2" != "--help" && "$2" != "-h" ]]; then
        local output
        output="$("$govman_bin" "$@" 2>&1)"
        local exit_code=$?
        if [[ $exit_code -eq 0 ]]; then
            local export_cmd=$(echo "$output" | grep -E '^export PATH=')
            if [[ -n "$export_cmd" ]]; then
                eval "$export_cmd"
                echo "✓ Go version switched successfully"
                return 0
            fi
        else
            echo "$output" >&2
            return $exit_code
        fi
    fi
    "$govman_bin" "$@"
}

# Auto-switch Go versions based on .govman-version file
govman_auto_switch() {
    # Check if auto-switch is enabled in config
    local config_file="$HOME/.govman/config.yaml"
    if [[ -f "$config_file" ]]; then
        local auto_switch_enabled=$(grep -E '^auto_switch:' -A 10 "$config_file" 2>/dev/null | grep -E '^[[:space:]]*enabled:' | head -1 | awk '{print $2}' | tr -d '[:space:]')
        if [[ "$auto_switch_enabled" != "true" ]]; then
            return 0
        fi
    fi

    if [[ -f .govman-version ]]; then
        local required_version=$(cat .govman-version 2>/dev/null | tr -d '\n\r' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        if [[ -n "$required_version" ]]; then
            if ! command -v go >/dev/null 2>&1; then
                echo "Go not found. Switching to Go $required_version..."
                govman use "$required_version" >/dev/null 2>&1 || {
                    echo "Warning: Failed to switch to Go $required_version. Install it with 'govman install $required_version'" >&2
                }
                return
            fi

            local current_version=$(go version 2>/dev/null | awk '{print $3}' | sed 's/go//')
            if [[ "$current_version" != "$required_version" ]]; then
                echo "Auto-switching to Go $required_version (required by .govman-version)"
                govman use "$required_version" >/dev/null 2>&1 || {
                    echo "Warning: Failed to switch to Go $required_version. Install it with 'govman install $required_version'" >&2
                }
            fi
        fi
    fi
}

# Hook into cd command for auto-switching
govman_cd() {
    builtin cd "$@" && govman_auto_switch
}
alias cd=govman_cd

# Run auto-switch on shell startup
govman_auto_switch
# END GOVMAN
```

### Zsh

Add the following to your `~/.zshrc`:

```zsh
# GOVMAN - Go Version Manager
export PATH="$HOME/.govman/bin:$PATH"
export GOTOOLCHAIN=local

# Go tool binaries on PATH
if [[ -n "$GOBIN" ]]; then
    export PATH="$GOBIN:$PATH"
fi
if command -v go >/dev/null 2>&1; then
    export PATH="$(go env GOPATH)/bin:$PATH"
fi
# Default GOPATH location
export PATH="$HOME/go/bin:$PATH"

# Wrapper function for automatic PATH execution
govman() {
    local govman_bin="$HOME/.govman/bin/govman"
    if [[ "$1" == "use" && "$#" -ge 2 && "$2" != "--help" && "$2" != "-h" ]]; then
        local output
        output="$("$govman_bin" "$@" 2>&1)"
        local exit_code=$?
        if [[ $exit_code -eq 0 ]]; then
            local export_cmd=$(echo "$output" | grep -E '^export PATH=')
            if [[ -n "$export_cmd" ]]; then
                eval "$export_cmd"
                echo "✓ Go version switched successfully"
                return 0
            fi
        else
            echo "$output" >&2
            return $exit_code
        fi
    fi
    "$govman_bin" "$@"
}

# Auto-switch Go versions based on .govman-version file
govman_auto_switch() {
    # Check if auto-switch is enabled in config
    local config_file="$HOME/.govman/config.yaml"
    if [[ -f "$config_file" ]]; then
        local auto_switch_enabled=$(grep -E '^auto_switch:' -A 10 "$config_file" 2>/dev/null | grep -E '^[[:space:]]*enabled:' | head -1 | awk '{print $2}' | tr -d '[:space:]')
        if [[ "$auto_switch_enabled" != "true" ]]; then
            return 0
        fi
    fi

    if [[ -f .govman-version ]]; then
        local required_version=$(cat .govman-version 2>/dev/null | tr -d '\n\r' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        if [[ -n "$required_version" ]]; then
            if ! command -v go >/dev/null 2>&1; then
                echo "Go not found. Switching to Go $required_version..."
                govman use "$required_version" >/dev/null 2>&1 || {
                    echo "Warning: Failed to switch to Go $required_version. Install it with 'govman install $required_version'" >&2
                }
                return
            fi

            local current_version=$(go version 2>/dev/null | awk '{print $3}' | sed 's/go//')
            if [[ "$current_version" != "$required_version" ]]; then
                echo "Auto-switching to Go $required_version (required by .govman-version)"
                govman use "$required_version" >/dev/null 2>&1 || {
                    echo "Warning: Failed to switch to Go $required_version. Install it with 'govman install $required_version'" >&2
                }
            fi
        fi
    fi
}

# Zsh-specific: Hook into chpwd for directory changes
if [[ -n "$ZSH_VERSION" ]]; then
  autoload -U add-zsh-hook
  add-zsh-hook chpwd govman_auto_switch
fi

# Run auto-switch on shell startup
govman_auto_switch
# END GOVMAN
```

### Fish

Add the following to your `~/.config/fish/config.fish`:

```fish
# GOVMAN - Go Version Manager
fish_add_path -p "$HOME/.govman/bin"
set -gx GOTOOLCHAIN local

# Go tool binaries on PATH
if set -q GOBIN
    fish_add_path -p "$GOBIN"
end
if type -q go
    fish_add_path -p (go env GOPATH)"/bin"
end
# Default GOPATH location
fish_add_path -p "$HOME/go/bin"

# Wrapper function for automatic PATH execution
function govman
    set govman_bin "$HOME/.govman/bin/govman"
    if test "$argv[1]" = "use"; and test (count $argv) -ge 2; and test "$argv[2]" != "--help"; and test "$argv[2]" != "-h"
        set output ($govman_bin $argv 2>&1)
        set exit_code $status
        if test $exit_code -eq 0
            for line in $output
                if string match -qr '^fish_add_path' -- $line
                    eval $line
                    echo "✓ Go version switched successfully"
                    return 0
                end
            end
        else
            for line in $output
                echo $line >&2
            end
            return $exit_code
        end
    end
    $govman_bin $argv
end

# Auto-switch Go versions based on .govman-version file
function govman_auto_switch
    set config_file "$HOME/.govman/config.yaml"
    if test -f "$config_file"
        set auto_switch_enabled (grep -E '^auto_switch:' -A 10 "$config_file" 2>/dev/null | grep -E '^[[:space:]]*enabled:' | head -1 | awk '{print $2}' | tr -d '[:space:]')
        if test "$auto_switch_enabled" != "true"
            return 0
        end
    end

    if test -f .govman-version
        set required_version (string trim < .govman-version)
        if test -n "$required_version"
            if not command -v go >/dev/null 2>&1
                echo "Go not found. Switching to Go $required_version..."
                govman use "$required_version" >/dev/null 2>&1; or begin
                    echo "Warning: Failed to switch to Go $required_version. Install it with 'govman install $required_version'" >&2
                end
                return
            end

            set current_version (go version 2>/dev/null | awk '{print $3}' | sed 's/go//')
            if test "$current_version" != "$required_version"
                echo "Auto-switching to Go $required_version (required by .govman-version)"
                govman use "$required_version" >/dev/null 2>&1; or begin
                    echo "Warning: Failed to switch to Go $required_version. Install it with 'govman install $required_version'" >&2
                end
            end
        end
    end
end

# Hook into cd command for auto-switching
function cd
    builtin cd $argv; and govman_auto_switch
end

# Fish-specific: Hook into directory changes via PWD variable
function __govman_cd_hook --on-variable PWD
    govman_auto_switch
end

# Run auto-switch on shell startup
govman_auto_switch
# END GOVMAN
```

### PowerShell

Add the following to your PowerShell profile (find it with `echo $PROFILE`):

```powershell
# GOVMAN - Go Version Manager
$env:PATH = "$env:USERPROFILE\.govman\bin;" + $env:PATH
$env:GOTOOLCHAIN = 'local'

# Go tool binaries on PATH
if ($env:GOBIN) { $env:PATH = "$env:GOBIN;" + $env:PATH }
if (Get-Command go -ErrorAction SilentlyContinue) { $env:PATH = "$(go env GOPATH)\bin;" + $env:PATH }
$env:PATH = "$env:USERPROFILE\go\bin;" + $env:PATH

# Wrapper function for automatic PATH execution
function govman {
    $govman_bin = "$env:USERPROFILE\.govman\bin\govman.exe"
    if ($args.Count -ge 2 -and $args[0] -eq 'use' -and $args[1] -ne '--help' -and $args[1] -ne '-h') {
        try {
            $output = & $govman_bin @args 2>&1
            if ($LASTEXITCODE -eq 0) {
                $pathCmd = $output | Where-Object { $_ -match '^\$env:PATH = ' }
                if ($pathCmd) {
                    Invoke-Expression $pathCmd
                    Write-Host '✓ Go version switched successfully' -ForegroundColor Green
                    return
                }
            } else {
                $output | ForEach-Object { Write-Error $_ }
                exit $LASTEXITCODE
            }
        } catch {
            Write-Error $_.Exception.Message
            exit 1
        }
    }
    & $govman_bin @args
    exit $LASTEXITCODE
}

# Auto-switch Go versions based on .govman-version file
function Invoke-GovmanAutoSwitch {
    $configFile = "$env:USERPROFILE\.govman\config.yaml"
    if (Test-Path $configFile) {
        try {
            $autoSwitchEnabled = $false
            $content = Get-Content $configFile -Raw -ErrorAction Stop
            if ($content -match '(?ms)auto_switch:.*?enabled:\s*(true|false)') {
                $autoSwitchEnabled = ($matches[1] -eq 'true')
            }
            if (-not $autoSwitchEnabled) {
                return
            }
        } catch {
            return
        }
    }

    if (Test-Path .govman-version) {
        try {
            $requiredVersion = (Get-Content .govman-version -Raw -ErrorAction Stop).Trim()
        } catch {
            return
        }

        if ($requiredVersion) {
            $currentVersion = $null
            try {
                $goVersionOutput = go version 2>$null
                if ($LASTEXITCODE -eq 0 -and $goVersionOutput) {
                    if ($goVersionOutput -match 'go version go([\d\.]+)') {
                        $currentVersion = $matches[1]
                    }
                }
            } catch {}

            if (-not $currentVersion) {
                Write-Host "Go not found. Switching to Go $requiredVersion..." -ForegroundColor Yellow
                govman use $requiredVersion *>$null
                if ($LASTEXITCODE -ne 0) {
                    Write-Warning "Failed to switch to Go $requiredVersion. Install it with 'govman install $requiredVersion'"
                }
                return
            }

            if ($currentVersion -ne $requiredVersion) {
                Write-Host "Auto-switching to Go $requiredVersion (required by .govman-version)" -ForegroundColor Yellow
                govman use $requiredVersion *>$null
                if ($LASTEXITCODE -ne 0) {
                    Write-Warning "Failed to switch to Go $requiredVersion. Install it with 'govman install $requiredVersion'"
                }
            }
        }
    }
}

# Hook into prompt for auto-switching
if (Get-Command prompt -ErrorAction SilentlyContinue) {
    $Global:GovmanOriginalPrompt = $function:prompt
    function global:prompt {
        Invoke-GovmanAutoSwitch
        if ($Global:GovmanOriginalPrompt) {
            & $Global:GovmanOriginalPrompt
        }
    }
}

# Run auto-switch on shell startup
Invoke-GovmanAutoSwitch
# END GOVMAN
```

## Troubleshooting

If auto-switching isn't working, try these steps:

1.  **Force Re-initialization**:
    This will remove any old `govman` configuration and write a fresh one.
    ```bash
    govman init --force
    ```

2.  **Verify Your `PATH`**:
    Ensure that `~/.govman/bin`, your `GOBIN` (if set), your `GOPATH/bin` (often `~/go/bin`), or the default `$HOME/go/bin` appear in your `PATH`.
    ```bash
    echo $PATH
    ```
    If any of those are missing, your shell configuration file may not be sourced correctly.

3.  **Check File Permissions**:
    Ensure your shell configuration file (e.g., `~/.zshrc`) is readable and that `govman` has permission to write to it.

4.  **Restart Your Terminal**:
    Shell configuration is only loaded when the shell starts. Make sure you've opened a new terminal window after running `govman init`.