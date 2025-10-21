package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"
)

var (
	currentGOOS        = runtime.GOOS
	execLookPath       = exec.LookPath
	userHomeDir        = os.UserHomeDir
	newlineRegex       = regexp.MustCompile(`\n{3,}`)
	configRemovalRegex = regexp.MustCompile(`(?ms)^[#\s]*(REM\s+)?GOVMAN - Go Version Manager.*?^[#\s]*(REM\s+)?END GOVMAN.*?$\n?`)
)

type Shell interface {
	Name() string
	DisplayName() string
	ConfigFile() string
	PathCommand(path string) string
	SetupCommands(binPath string) []string
	IsAvailable() bool
	ExecutePathCommand(path string) error
}

type BashShell struct{}
type ZshShell struct{}
type FishShell struct{}
type PowerShell struct{}
type CmdShell struct{}

// validateBinPath ensures the binary path is safe and exists
func validateBinPath(binPath string) error {
	if binPath == "" {
		return fmt.Errorf("binary path cannot be empty")
	}

	// Check for path traversal indicators
	if strings.Contains(binPath, "..") {
		return fmt.Errorf("invalid binary path (path traversal detected): %s", binPath)
	}

	// Clean the path to prevent directory traversal
	cleanPath := filepath.Clean(binPath)

	// Convert to absolute path for comparison
	absPath, err := filepath.Abs(binPath)
	if err != nil {
		return fmt.Errorf("unable to resolve absolute path: %w", err)
	}

	absCleanPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("unable to resolve clean absolute path: %w", err)
	}

	// Ensure no path traversal
	if absPath != absCleanPath {
		return fmt.Errorf("invalid binary path (path traversal detected): %s", binPath)
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("binary path does not exist: %s", absPath)
		}
		return fmt.Errorf("unable to access binary path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("binary path is not a directory: %s", absPath)
	}

	return nil
}

// escapeBashPath properly escapes a path for use in bash/zsh
func escapeBashPath(path string) string {
	// Escape special characters for bash/zsh
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		`$`, `\$`,
		"`", "\\`",
		`!`, `\!`,
	)
	return replacer.Replace(path)
}

// escapeFishPath properly escapes a path for use in fish
func escapeFishPath(path string) string {
	// Fish uses different escaping rules
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		`$`, `\$`,
	)
	return replacer.Replace(path)
}

// escapePowerShellPath properly escapes a path for use in PowerShell
func escapePowerShellPath(path string) string {
	// PowerShell escaping: backtick is the escape character
	// Order matters: escape backtick first
	replacer := strings.NewReplacer(
		"`", "``",
		`"`, "`\"",
		`$`, "`$",
	)
	return replacer.Replace(path)
}

// escapeCmdPath properly escapes a path for use in cmd
func escapeCmdPath(path string) string {
	// CMD uses % for variables
	return strings.ReplaceAll(path, "%", "%%")
}

// Detect determines the user's shell based on OS and environment variables,
// falling back to an available default when detection is inconclusive.
func Detect() Shell {
	if currentGOOS == "windows" {
		// Check for PowerShell Core first (preferred)
		if isCommandAvailable("pwsh") {
			return &PowerShell{}
		}
		if isCommandAvailable("powershell") {
			return &PowerShell{}
		}

		// Fallback to Command Prompt
		return &CmdShell{}
	}

	// For Unix-like systems, check SHELL environment variable
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		return detectAvailableShell()
	}

	shellName := filepath.Base(shellPath)
	switch shellName {
	case "zsh":
		if isCommandAvailable("zsh") {
			return &ZshShell{}
		}
	case "fish":
		if isCommandAvailable("fish") {
			return &FishShell{}
		}
	case "bash", "sh":
		if isCommandAvailable("bash") {
			return &BashShell{}
		}
	}

	// If the detected shell isn't available, find an alternative
	return detectAvailableShell()
}

// DetectAll returns a slice of supported shells that are available on the current system.
func DetectAll() []Shell {
	var shells []Shell

	if currentGOOS == "windows" {
		// Windows-specific shells
		shells = []Shell{
			&PowerShell{},
			&CmdShell{},
		}
	} else {
		// Unix-like shells
		shells = []Shell{
			&ZshShell{},
			&BashShell{},
			&FishShell{},
		}
	}

	var available []Shell
	for _, shell := range shells {
		if shell.IsAvailable() {
			available = append(available, shell)
		}
	}

	return available
}

// detectAvailableShell returns the first available shell from a prioritized list.
func detectAvailableShell() Shell {
	shells := []Shell{
		&BashShell{},
		&ZshShell{},
		&FishShell{},
	}

	for _, shell := range shells {
		if shell.IsAvailable() {
			return shell
		}
	}

	return &BashShell{}
}

// isCommandAvailable reports whether a command exists in the system PATH.
func isCommandAvailable(command string) bool {
	_, err := execLookPath(command)
	return err == nil
}

// fileExists checks if a file exists and is not a directory.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Name returns the identifier for Bash.
func (s *BashShell) Name() string {
	return "bash"
}

// DisplayName returns the human-friendly name for Bash.
func (s *BashShell) DisplayName() string {
	return "Bash"
}

// IsAvailable reports whether Bash is present in the system PATH.
func (s *BashShell) IsAvailable() bool {
	return isCommandAvailable("bash")
}

// ConfigFile returns the path to the Bash configuration file.
func (s *BashShell) ConfigFile() string {
	home, err := userHomeDir()
	if err != nil {
		return ".bashrc" // Fallback to relative path
	}

	candidates := []string{
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".bash_profile"),
		filepath.Join(home, ".profile"),
	}

	for _, candidate := range candidates {
		if fileExists(candidate) {
			return candidate
		}
	}

	// Default to .bashrc if none exist
	return filepath.Join(home, ".bashrc")
}

// PathCommand returns a Bash-compatible command to prepend binPath to PATH.
func (s *BashShell) PathCommand(path string) string {
	escapedPath := escapeBashPath(path)
	return fmt.Sprintf(`export PATH="%s:$PATH"`, escapedPath)
}

// SetupCommands returns the Bash shell configuration lines to integrate govman.
func (s *BashShell) SetupCommands(binPath string) []string {
	escapedPath := escapeBashPath(binPath)

	commands := []string{
		"# GOVMAN - Go Version Manager",
		fmt.Sprintf(`export PATH="%s:$PATH"`, escapedPath),
		"export GOTOOLCHAIN=local",
		"",
		"# Wrapper function for automatic PATH execution",
		"govman() {",
		fmt.Sprintf(`    local govman_bin="%s/govman"`, escapedPath),
		`    if [[ "$1" == "use" && "$#" -ge 2 && "$2" != "--help" && "$2" != "-h" ]]; then`,
		"        local output",
		`        output="$("$govman_bin" "$@" 2>&1)"`,
		"        local exit_code=$?",
		"        if [[ $exit_code -eq 0 ]]; then",
		`            local export_cmd=$(echo "$output" | grep -E '^export PATH=')`,
		`            if [[ -n "$export_cmd" ]]; then`,
		`                eval "$export_cmd"`,
		`                echo "✓ Go version switched successfully"`,
		"                return 0",
		"            fi",
		"        else",
		`            echo "$output" >&2`,
		"            return $exit_code",
		"        fi",
		"    fi",
		`    "$govman_bin" "$@"`,
		"}",
		"",
		"# Auto-switch Go versions based on .govman-version file",
		"govman_auto_switch() {",
		"    # Check if auto-switch is enabled in config",
		`    local config_file="$HOME/.govman/config.yaml"`,
		`    if [[ -f "$config_file" ]]; then`,
		`        local auto_switch_enabled=$(grep -E '^auto_switch:' -A 10 "$config_file" 2>/dev/null | grep -E '^[[:space:]]*enabled:' | head -1 | awk '{print $2}' | tr -d '[:space:]')`,
		`        if [[ "$auto_switch_enabled" != "true" ]]; then`,
		"            return 0",
		"        fi",
		"    fi",
		"",
		"    if [[ -f .govman-version ]]; then",
		`        local required_version=$(cat .govman-version 2>/dev/null | tr -d '\n\r' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')`,
		`        if [[ -n "$required_version" ]]; then`,
		"            if ! command -v go >/dev/null 2>&1; then",
		`                echo "Go not found. Switching to Go $required_version..."`,
		`                govman use "$required_version" >/dev/null 2>&1 || {`,
		`                    echo "Warning: Failed to switch to Go $required_version. Install it with 'govman install $required_version'" >&2`,
		"                }",
		"                return",
		"            fi",
		"",
		`            local current_version=$(go version 2>/dev/null | awk '{print $3}' | sed 's/go//')`,
		`            if [[ "$current_version" != "$required_version" ]]; then`,
		`                echo "Auto-switching to Go $required_version (required by .govman-version)"`,
		`                govman use "$required_version" >/dev/null 2>&1 || {`,
		`                    echo "Warning: Failed to switch to Go $required_version. Install it with 'govman install $required_version'" >&2`,
		"                }",
		"            fi",
		"        fi",
		"    fi",
		"}",
		"",
		"# Hook into cd command for auto-switching",
		"govman_cd() {",
		`    builtin cd "$@" && govman_auto_switch`,
		"}",
		"alias cd=govman_cd",
		"",
		"# Run auto-switch on shell startup",
		"govman_auto_switch",
		"# END GOVMAN",
	}

	return commands
}

// ExecutePathCommand outputs the PATH command for automatic execution via eval.
func (s *BashShell) ExecutePathCommand(path string) error {
	if err := validateBinPath(path); err != nil {
		return err
	}

	pathCmd := s.PathCommand(path)

	// Output the command for eval
	fmt.Println(pathCmd)

	// Instructions to stderr so they don't interfere with eval
	fmt.Fprintf(os.Stderr, "# To apply to current session, run:\n")
	fmt.Fprintf(os.Stderr, "# eval \"$(govman use <version>)\"\n")

	return nil
}

// Name returns the identifier for Zsh.
func (s *ZshShell) Name() string {
	return "zsh"
}

// DisplayName returns the human-friendly name for Zsh.
func (s *ZshShell) DisplayName() string {
	return "Zsh"
}

// IsAvailable reports whether Zsh is present in the system PATH.
func (s *ZshShell) IsAvailable() bool {
	return isCommandAvailable("zsh")
}

// ConfigFile returns the path to the Zsh configuration file.
func (s *ZshShell) ConfigFile() string {
	home, err := userHomeDir()
	if err != nil {
		return ".zshrc"
	}
	return filepath.Join(home, ".zshrc")
}

// PathCommand returns a Zsh-compatible command to prepend binPath to PATH.
func (s *ZshShell) PathCommand(path string) string {
	escapedPath := escapeBashPath(path)
	return fmt.Sprintf(`export PATH="%s:$PATH"`, escapedPath)
}

// SetupCommands returns the Zsh configuration lines to integrate govman.
func (s *ZshShell) SetupCommands(binPath string) []string {
	escapedPath := escapeBashPath(binPath)

	commands := []string{
		"# GOVMAN - Go Version Manager",
		fmt.Sprintf(`export PATH="%s:$PATH"`, escapedPath),
		"export GOTOOLCHAIN=local",
		"",
		"# Wrapper function for automatic PATH execution",
		"govman() {",
		fmt.Sprintf(`    local govman_bin="%s/govman"`, escapedPath),
		`    if [[ "$1" == "use" && "$#" -ge 2 && "$2" != "--help" && "$2" != "-h" ]]; then`,
		"        local output",
		`        output="$("$govman_bin" "$@" 2>&1)"`,
		"        local exit_code=$?",
		"        if [[ $exit_code -eq 0 ]]; then",
		`            local export_cmd=$(echo "$output" | grep -E '^export PATH=')`,
		`            if [[ -n "$export_cmd" ]]; then`,
		`                eval "$export_cmd"`,
		`                echo "✓ Go version switched successfully"`,
		"                return 0",
		"            fi",
		"        else",
		`            echo "$output" >&2`,
		"            return $exit_code",
		"        fi",
		"    fi",
		`    "$govman_bin" "$@"`,
		"}",
		"",
		"# Auto-switch Go versions based on .govman-version file",
		"govman_auto_switch() {",
		"    # Check if auto-switch is enabled in config",
		`    local config_file="$HOME/.govman/config.yaml"`,
		`    if [[ -f "$config_file" ]]; then`,
		`        local auto_switch_enabled=$(grep -E '^auto_switch:' -A 10 "$config_file" 2>/dev/null | grep -E '^[[:space:]]*enabled:' | head -1 | awk '{print $2}' | tr -d '[:space:]')`,
		`        if [[ "$auto_switch_enabled" != "true" ]]; then`,
		"            return 0",
		"        fi",
		"    fi",
		"",
		"    if [[ -f .govman-version ]]; then",
		`        local required_version=$(cat .govman-version 2>/dev/null | tr -d '\n\r' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')`,
		`        if [[ -n "$required_version" ]]; then`,
		"            if ! command -v go >/dev/null 2>&1; then",
		`                echo "Go not found. Switching to Go $required_version..."`,
		`                govman use "$required_version" >/dev/null 2>&1 || {`,
		`                    echo "Warning: Failed to switch to Go $required_version. Install it with 'govman install $required_version'" >&2`,
		"                }",
		"                return",
		"            fi",
		"",
		`            local current_version=$(go version 2>/dev/null | awk '{print $3}' | sed 's/go//')`,
		`            if [[ "$current_version" != "$required_version" ]]; then`,
		`                echo "Auto-switching to Go $required_version (required by .govman-version)"`,
		`                govman use "$required_version" >/dev/null 2>&1 || {`,
		`                    echo "Warning: Failed to switch to Go $required_version. Install it with 'govman install $required_version'" >&2`,
		"                }",
		"            fi",
		"        fi",
		"    fi",
		"}",
		"",
		"# Zsh-specific: Hook into chpwd for directory changes",
		"autoload -U add-zsh-hook",
		"add-zsh-hook chpwd govman_auto_switch",
		"",
		"# Run auto-switch on shell startup",
		"govman_auto_switch",
		"# END GOVMAN",
	}

	return commands
}

// ExecutePathCommand outputs the PATH command for automatic execution via eval.
func (s *ZshShell) ExecutePathCommand(path string) error {
	if err := validateBinPath(path); err != nil {
		return err
	}

	pathCmd := s.PathCommand(path)
	fmt.Println(pathCmd)

	fmt.Fprintf(os.Stderr, "# To apply to current session, run:\n")
	fmt.Fprintf(os.Stderr, "# eval \"$(govman use <version>)\"\n")

	return nil
}

// Name returns the identifier for Fish.
func (s *FishShell) Name() string {
	return "fish"
}

// DisplayName returns the human-friendly name for Fish.
func (s *FishShell) DisplayName() string {
	return "Fish"
}

// IsAvailable reports whether Fish is present in the system PATH.
func (s *FishShell) IsAvailable() bool {
	return isCommandAvailable("fish")
}

// ConfigFile returns the path to the Fish configuration file.
func (s *FishShell) ConfigFile() string {
	home, err := userHomeDir()
	if err != nil {
		return "config.fish"
	}
	return filepath.Join(home, ".config", "fish", "config.fish")
}

// PathCommand returns a Fish-compatible command to prepend binPath to PATH.
func (s *FishShell) PathCommand(path string) string {
	escapedPath := escapeFishPath(path)
	return fmt.Sprintf(`fish_add_path -p "%s"`, escapedPath)
}

// SetupCommands returns the Fish configuration lines to integrate govman.
func (s *FishShell) SetupCommands(binPath string) []string {
	escapedPath := escapeFishPath(binPath)

	commands := []string{
		"# GOVMAN - Go Version Manager",
		fmt.Sprintf(`fish_add_path -p "%s"`, escapedPath),
		"set -gx GOTOOLCHAIN local",
		"",
		"# Wrapper function for automatic PATH execution",
		"function govman",
		fmt.Sprintf(`    set govman_bin "%s/govman"`, escapedPath),
		`    if test "$argv[1]" = "use"; and test (count $argv) -ge 2; and test "$argv[2]" != "--help"; and test "$argv[2]" != "-h"`,
		"        set output ($govman_bin $argv 2>&1)",
		"        set exit_code $status",
		"        if test $exit_code -eq 0",
		"            for line in $output",
		"                if string match -qr '^fish_add_path' -- $line",
		"                    eval $line",
		`                    echo "✓ Go version switched successfully"`,
		"                    return 0",
		"                end",
		"            end",
		"        else",
		"            for line in $output",
		"                echo $line >&2",
		"            end",
		"            return $exit_code",
		"        end",
		"    end",
		"    $govman_bin $argv",
		"end",
		"",
		"# Auto-switch Go versions based on .govman-version file",
		"function govman_auto_switch",
		`    set config_file "$HOME/.govman/config.yaml"`,
		`    if test -f "$config_file"`,
		`        set auto_switch_enabled (grep -E '^auto_switch:' -A 10 "$config_file" 2>/dev/null | grep -E '^[[:space:]]*enabled:' | head -1 | awk '{print $2}' | tr -d '[:space:]')`,
		`        if test "$auto_switch_enabled" != "true"`,
		"            return 0",
		"        end",
		"    end",
		"",
		"    if test -f .govman-version",
		"        set required_version (string trim < .govman-version)",
		`        if test -n "$required_version"`,
		"            if not command -v go >/dev/null 2>&1",
		`                echo "Go not found. Switching to Go $required_version..."`,
		`                govman use "$required_version" >/dev/null 2>&1; or begin`,
		`                    echo "Warning: Failed to switch to Go $required_version. Install it with 'govman install $required_version'" >&2`,
		"                end",
		"                return",
		"            end",
		"",
		"            set current_version (go version 2>/dev/null | awk '{print $3}' | sed 's/go//')",
		`            if test "$current_version" != "$required_version"`,
		`                echo "Auto-switching to Go $required_version (required by .govman-version)"`,
		`                govman use "$required_version" >/dev/null 2>&1; or begin`,
		`                    echo "Warning: Failed to switch to Go $required_version. Install it with 'govman install $required_version'" >&2`,
		"                end",
		"            end",
		"        end",
		"    end",
		"end",
		"",
		"# Hook into cd command for auto-switching",
		"function cd",
		"    builtin cd $argv; and govman_auto_switch",
		"end",
		"",
		"# Fish-specific: Hook into directory changes",
		"function __govman_cd_hook --on-variable PWD",
		"    govman_auto_switch",
		"end",
		"",
		"# Run auto-switch on shell startup",
		"govman_auto_switch",
		"# END GOVMAN",
	}

	return commands
}

// ExecutePathCommand outputs the PATH command for automatic execution via eval.
func (s *FishShell) ExecutePathCommand(path string) error {
	if err := validateBinPath(path); err != nil {
		return err
	}

	pathCmd := s.PathCommand(path)
	fmt.Println(pathCmd)

	fmt.Fprintf(os.Stderr, "# To apply to current session, run:\n")
	fmt.Fprintf(os.Stderr, "# eval (govman use <version>)\n")

	return nil
}

// Name returns the identifier for PowerShell.
func (s *PowerShell) Name() string {
	return "powershell"
}

// DisplayName returns the human-friendly name for PowerShell.
func (s *PowerShell) DisplayName() string {
	return "PowerShell"
}

// IsAvailable reports whether PowerShell is available.
func (s *PowerShell) IsAvailable() bool {
	return isCommandAvailable("pwsh") || isCommandAvailable("powershell")
}

// ConfigFile returns the PowerShell profile path.
func (s *PowerShell) ConfigFile() string {
	home, err := userHomeDir()
	if err != nil {
		return "$PROFILE"
	}

	// Check for PowerShell Core first
	if isCommandAvailable("pwsh") {
		profilePath := filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
		return profilePath
	}

	// Fall back to Windows PowerShell
	profilePath := filepath.Join(home, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1")
	return profilePath
}

// PathCommand returns a PowerShell command to prepend binPath to PATH.
func (s *PowerShell) PathCommand(path string) string {
	escapedPath := escapePowerShellPath(path)
	return fmt.Sprintf(`$env:PATH = "%s;" + $env:PATH`, escapedPath)
}

// SetupCommands returns the PowerShell profile lines to integrate govman.
func (s *PowerShell) SetupCommands(binPath string) []string {
	escapedPath := escapePowerShellPath(binPath)

	commands := []string{
		"# GOVMAN - Go Version Manager",
		fmt.Sprintf(`$env:PATH = "%s;" + $env:PATH`, escapedPath),
		"$env:GOTOOLCHAIN = 'local'",
		"",
		"# Wrapper function for automatic PATH execution",
		"function govman {",
		fmt.Sprintf(`    $govman_bin = "%s\govman.exe"`, escapedPath),
		"    if ($args.Count -ge 2 -and $args[0] -eq 'use' -and $args[1] -ne '--help' -and $args[1] -ne '-h') {",
		"        try {",
		"            $output = & $govman_bin @args 2>&1",
		"            if ($LASTEXITCODE -eq 0) {",
		"                $pathCmd = $output | Where-Object { $_ -match '^\\$env:PATH = ' }",
		"                if ($pathCmd) {",
		"                    Invoke-Expression $pathCmd",
		"                    Write-Host '✓ Go version switched successfully' -ForegroundColor Green",
		"                    return",
		"                }",
		"            } else {",
		"                $output | ForEach-Object { Write-Error $_ }",
		"                exit $LASTEXITCODE",
		"            }",
		"        } catch {",
		"            Write-Error $_.Exception.Message",
		"            exit 1",
		"        }",
		"    }",
		"    & $govman_bin @args",
		"    exit $LASTEXITCODE",
		"}",
		"",
		"# Auto-switch Go versions based on .govman-version file",
		"function Invoke-GovmanAutoSwitch {",
		"    $configFile = \"$env:USERPROFILE\\.govman\\config.yaml\"",
		"    if (Test-Path $configFile) {",
		"        try {",
		"            $autoSwitchEnabled = $false",
		"            $content = Get-Content $configFile -Raw -ErrorAction Stop",
		"            if ($content -match '(?ms)auto_switch:.*?enabled:\\s*(true|false)') {",
		"                $autoSwitchEnabled = ($matches[1] -eq 'true')",
		"            }",
		"            if (-not $autoSwitchEnabled) {",
		"                return",
		"            }",
		"        } catch {",
		"            return",
		"        }",
		"    }",
		"",
		"    if (Test-Path .govman-version) {",
		"        try {",
		"            $requiredVersion = (Get-Content .govman-version -Raw -ErrorAction Stop).Trim()",
		"        } catch {",
		"            return",
		"        }",
		"",
		"        if ($requiredVersion) {",
		"            $currentVersion = $null",
		"            try {",
		"                $goVersionOutput = go version 2>$null",
		"                if ($LASTEXITCODE -eq 0 -and $goVersionOutput) {",
		"                    if ($goVersionOutput -match 'go version go([\\d\\.]+)') {",
		"                        $currentVersion = $matches[1]",
		"                    }",
		"                }",
		"            } catch {}",
		"",
		"            if (-not $currentVersion) {",
		"                Write-Host \"Go not found. Switching to Go $requiredVersion...\" -ForegroundColor Yellow",
		"                govman use $requiredVersion *>$null",
		"                if ($LASTEXITCODE -ne 0) {",
		"                    Write-Warning \"Failed to switch to Go $requiredVersion. Install it with 'govman install $requiredVersion'\"",
		"                }",
		"                return",
		"            }",
		"",
		"            if ($currentVersion -ne $requiredVersion) {",
		"                Write-Host \"Auto-switching to Go $requiredVersion (required by .govman-version)\" -ForegroundColor Yellow",
		"                govman use $requiredVersion *>$null",
		"                if ($LASTEXITCODE -ne 0) {",
		"                    Write-Warning \"Failed to switch to Go $requiredVersion. Install it with 'govman install $requiredVersion'\"",
		"                }",
		"            }",
		"        }",
		"    }",
		"}",
		"",
		"# Hook into prompt for auto-switching",
		"if (Get-Command prompt -ErrorAction SilentlyContinue) {",
		"    $Global:GovmanOriginalPrompt = $function:prompt",
		"    function global:prompt {",
		"        Invoke-GovmanAutoSwitch",
		"        if ($Global:GovmanOriginalPrompt) {",
		"            & $Global:GovmanOriginalPrompt",
		"        }",
		"    }",
		"}",
		"",
		"# Run auto-switch on shell startup",
		"Invoke-GovmanAutoSwitch",
		"# END GOVMAN",
	}

	return commands
}

// ExecutePathCommand outputs the PATH command for automatic execution.
func (s *PowerShell) ExecutePathCommand(path string) error {
	if err := validateBinPath(path); err != nil {
		return err
	}

	pathCmd := s.PathCommand(path)
	fmt.Println(pathCmd)

	fmt.Fprintf(os.Stderr, "# To apply to current session, run:\n")
	fmt.Fprintf(os.Stderr, "# govman use <version> | Invoke-Expression\n")

	return nil
}

// Name returns the identifier for Windows Command Prompt.
func (s *CmdShell) Name() string {
	return "cmd"
}

// DisplayName returns the human-friendly name for Command Prompt.
func (s *CmdShell) DisplayName() string {
	return "Command Prompt"
}

// IsAvailable reports whether cmd is available (Windows only).
func (s *CmdShell) IsAvailable() bool {
	return currentGOOS == "windows"
}

// ConfigFile returns a description of where cmd configuration is managed.
func (s *CmdShell) ConfigFile() string {
	return "Environment Variables (System Properties)"
}

// PathCommand returns a cmd.exe command to prepend binPath to PATH.
func (s *CmdShell) PathCommand(path string) string {
	escapedPath := escapeCmdPath(path)
	return fmt.Sprintf(`set PATH=%s;%%PATH%%`, escapedPath)
}

// SetupCommands returns guidance for integrating govman with Command Prompt.
func (s *CmdShell) SetupCommands(binPath string) []string {
	escapedPath := escapeCmdPath(binPath)

	return []string{
		"REM GOVMAN - Go Version Manager",
		fmt.Sprintf(`set PATH=%s;%%PATH%%`, escapedPath),
		"set GOTOOLCHAIN=local",
		"",
		"REM Note: CMD has limited scripting capabilities.",
		"REM For best experience, use PowerShell instead.",
		"REM END GOVMAN",
	}
}

// ExecutePathCommand outputs the PATH command for Command Prompt.
func (s *CmdShell) ExecutePathCommand(path string) error {
	if err := validateBinPath(path); err != nil {
		return err
	}

	pathCmd := s.PathCommand(path)
	fmt.Println(pathCmd)

	fmt.Fprintln(os.Stderr, "REM To apply to current session, copy and run:")
	fmt.Fprintf(os.Stderr, "REM %s\n", pathCmd)

	return nil
}

// InitializeShell sets up shell integration for govman.
func InitializeShell(shell Shell, binPath string, force bool) error {
	// Validate the binary path first
	if err := validateBinPath(binPath); err != nil {
		return fmt.Errorf("invalid binary path: %w", err)
	}

	switch shell.Name() {
	case "powershell":
		return initializePowerShell(shell, binPath, force)
	case "cmd":
		return initializeCmdShell(shell, binPath, force)
	default:
		return initializeUnixShell(shell, binPath, force)
	}
}

// initializeUnixShell writes govman integration to the shell config file.
func initializeUnixShell(shell Shell, binPath string, force bool) error {
	configFile := shell.ConfigFile()

	// Create config directory if needed
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	// Verify we can write to the directory
	testFile := filepath.Join(configDir, ".govman_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("insufficient permissions to write to %s: %w", configDir, err)
	}
	os.Remove(testFile)

	// Read existing content
	var existingContent string
	if content, err := os.ReadFile(configFile); err == nil {
		existingContent = string(content)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Check if govman is already configured
	if containsGovmanConfig(existingContent) {
		if !force {
			return fmt.Errorf("govman is already configured in %s (use --force to override)", configFile)
		}
		existingContent = removeExistingConfig(existingContent)
	}

	// Prepare new configuration
	setupCommands := shell.SetupCommands(binPath)
	newConfig := "\n" + strings.Join(setupCommands, "\n") + "\n"

	// Combine content
	finalContent := strings.TrimSpace(existingContent) + newConfig

	// Write to file with proper permissions
	if err := os.WriteFile(configFile, []byte(finalContent), 0644); err != nil {
		return fmt.Errorf("failed to write config to %s: %w", configFile, err)
	}

	fmt.Printf("✅ Successfully configured %s\n", shell.DisplayName())
	fmt.Printf("📝 Configuration added to: %s\n", configFile)
	fmt.Printf("🔄 Reload your shell or run: source %s\n", configFile)

	return nil
}

// initializePowerShell writes configuration to PowerShell profile.
func initializePowerShell(shell Shell, binPath string, force bool) error {
	profilePath := shell.ConfigFile()

	// Create profile directory if needed
	profileDir := filepath.Dir(profilePath)
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	// Verify write permissions
	testFile := filepath.Join(profileDir, ".govman_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("insufficient permissions to write to %s: %w", profileDir, err)
	}
	os.Remove(testFile)

	// Read existing content
	var existingContent string
	if content, err := os.ReadFile(profilePath); err == nil {
		existingContent = string(content)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to read profile: %w", err)
	}

	// Check if govman is already configured
	if containsGovmanConfig(existingContent) {
		if !force {
			return fmt.Errorf("govman is already configured in PowerShell profile (use --force to override)")
		}
		existingContent = removeExistingConfig(existingContent)
	}

	// Prepare new configuration
	setupCommands := shell.SetupCommands(binPath)
	newConfig := "\r\n" + strings.Join(setupCommands, "\r\n") + "\r\n"

	// Combine content
	finalContent := strings.TrimSpace(existingContent) + newConfig

	// Write to file
	if err := os.WriteFile(profilePath, []byte(finalContent), 0644); err != nil {
		return fmt.Errorf("failed to write PowerShell profile: %w", err)
	}

	fmt.Printf("✅ Successfully configured PowerShell\n")
	fmt.Printf("📝 Configuration added to: %s\n", profilePath)
	fmt.Printf("🔄 Reload PowerShell or run: . $PROFILE\n")

	return nil
}

// initializeCmdShell creates a batch wrapper for Command Prompt.
func initializeCmdShell(shell Shell, binPath string, force bool) error {
	wrapperPath := filepath.Join(binPath, "govman_wrapper.bat")

	// Check if wrapper exists
	if !force && fileExists(wrapperPath) {
		return fmt.Errorf("wrapper already exists at %s (use --force to override)", wrapperPath)
	}

	// Verify write permissions
	testFile := filepath.Join(binPath, ".govman_test")
	if err := os.WriteFile(testFile, []byte("test"), 0755); err != nil {
		return fmt.Errorf("insufficient permissions to write to %s: %w", binPath, err)
	}
	os.Remove(testFile)

	// Create wrapper content using template for better maintainability
	tmpl := `@echo off
setlocal enabledelayedexpansion

REM GOVMAN Wrapper for Command Prompt
set "GOVMAN_BIN={{.BinPath}}\govman.exe"

if "%~1"=="use" (
    if not "%~2"=="" (
        if not "%~2"=="--help" (
            if not "%~2"=="-h" (
                "%GOVMAN_BIN%" use %2 > "%TEMP%\govman_path.tmp" 2>nul
                if !errorlevel! equ 0 (
                    for /f "usebackq delims=" %%i in ("%TEMP%\govman_path.tmp") do (
                        echo %%i | findstr /b "set PATH=" >nul && %%i
                    )
                    del "%TEMP%\govman_path.tmp" 2>nul
                    echo Go version switched successfully
                    exit /b 0
                ) else (
                    if exist "%TEMP%\govman_path.tmp" del "%TEMP%\govman_path.tmp" 2>nul
                )
            )
        )
    )
)

"%GOVMAN_BIN%" %*
exit /b %errorlevel%
`

	// Parse and execute template
	t, err := template.New("wrapper").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse wrapper template: %w", err)
	}

	var buf strings.Builder
	data := struct {
		BinPath string
	}{
		BinPath: binPath,
	}

	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to generate wrapper: %w", err)
	}

	// Write wrapper file
	if err := os.WriteFile(wrapperPath, []byte(buf.String()), 0755); err != nil {
		return fmt.Errorf("failed to create wrapper: %w", err)
	}

	fmt.Printf("✅ Created govman wrapper: %s\n\n", wrapperPath)
	fmt.Println("📝 To complete setup, add govman to your PATH:")
	fmt.Printf("   setx PATH \"%%PATH%%;%s\"\n\n", binPath)
	fmt.Println("💡 For better experience, consider using PowerShell instead:")
	fmt.Println("   powershell -Command \"govman init\"")

	return nil
}

// containsGovmanConfig checks if content contains govman configuration.
func containsGovmanConfig(content string) bool {
	markers := []string{
		"GOVMAN - Go Version Manager",
		"govman_auto_switch",
		"Invoke-GovmanAutoSwitch",
		"__govman_cd_hook",
	}

	for _, marker := range markers {
		if strings.Contains(content, marker) {
			return true
		}
	}

	return false
}

// removeExistingConfig removes existing govman configuration from content.
func removeExistingConfig(content string) string {
	// Use the pre-compiled regex for better performance
	cleanedContent := configRemovalRegex.ReplaceAllString(content, "")

	// Clean up excessive newlines
	cleanedContent = newlineRegex.ReplaceAllString(cleanedContent, "\n\n")

	return strings.TrimSpace(cleanedContent)
}

// GetShellInstructions returns manual setup instructions for a shell.
func GetShellInstructions(shell Shell, binPath string) string {
	var instructions strings.Builder

	instructions.WriteString(fmt.Sprintf("Manual setup for %s:\n\n", shell.DisplayName()))
	instructions.WriteString(fmt.Sprintf("1. Edit: %s\n\n", shell.ConfigFile()))
	instructions.WriteString("2. Add these lines:\n\n")

	commands := shell.SetupCommands(binPath)
	for _, cmd := range commands {
		instructions.WriteString(fmt.Sprintf("   %s\n", cmd))
	}

	instructions.WriteString("\n3. Reload your shell:\n")

	switch shell.Name() {
	case "fish":
		instructions.WriteString("   source ~/.config/fish/config.fish\n")
	case "powershell":
		instructions.WriteString("   . $PROFILE\n")
	case "cmd":
		instructions.WriteString("   (Restart Command Prompt)\n")
	default:
		instructions.WriteString(fmt.Sprintf("   source %s\n", shell.ConfigFile()))
	}

	return instructions.String()
}
