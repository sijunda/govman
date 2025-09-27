package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// Shell represents a shell environment with its specific configurations
type Shell interface {
	Name() string
	DisplayName() string
	ConfigFile() string
	PathCommand(path string) string
	SetupCommands(binPath string) []string
	IsAvailable() bool
}

// Shell implementations
type BashShell struct{}
type ZshShell struct{}
type FishShell struct{}
type PowerShell struct{}
type CmdShell struct{} // Windows Command Prompt

// Detect automatically detects the current shell with fallback logic
func Detect() Shell {
	if runtime.GOOS == "windows" {
		// Try PowerShell first, then Command Prompt
		if isCommandAvailable("pwsh") || isCommandAvailable("powershell") {
			return &PowerShell{}
		}
		return &CmdShell{}
	}

	// Check SHELL environment variable
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		// Fallback: try to detect available shells
		return detectAvailableShell()
	}

	shellName := filepath.Base(shellPath)
	switch shellName {
	case "zsh":
		return &ZshShell{}
	case "fish":
		return &FishShell{}
	case "bash", "sh":
		return &BashShell{}
	default:
		// Try to detect what's actually available
		return detectAvailableShell()
	}
}

// DetectAll returns all available shells on the system
func DetectAll() []Shell {
	shells := []Shell{
		&BashShell{},
		&ZshShell{},
		&FishShell{},
	}

	if runtime.GOOS == "windows" {
		shells = append(shells, &PowerShell{}, &CmdShell{})
	}

	var available []Shell
	for _, shell := range shells {
		if shell.IsAvailable() {
			available = append(available, shell)
		}
	}

	return available
}

// Helper function to detect available shell as fallback
func detectAvailableShell() Shell {
	shells := []Shell{
		&ZshShell{},
		&BashShell{},
		&FishShell{},
	}

	for _, shell := range shells {
		if shell.IsAvailable() {
			return shell
		}
	}

	// Ultimate fallback
	return &BashShell{}
}

// Helper function to check if a command is available
func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// BashShell implementation
func (s *BashShell) Name() string {
	return "bash"
}

func (s *BashShell) DisplayName() string {
	return "Bash"
}

func (s *BashShell) IsAvailable() bool {
	return isCommandAvailable("bash")
}

func (s *BashShell) ConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Use current directory as fallback instead of /tmp
		home = "."
	}
	// Priority order: .bashrc > .bash_profile > .profile
	candidates := []string{
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".bash_profile"),
		filepath.Join(home, ".profile"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// Default to .bashrc if none exist
	return candidates[0]
}

func (s *BashShell) PathCommand(path string) string {
	return fmt.Sprintf(`export PATH="%s:$PATH"`, path)
}

func (s *BashShell) SetupCommands(binPath string) []string {
	return []string{
		"# GOVMAN - Go Version Manager",
		"export PATH=\"$HOME/.govman/bin:$PATH\"",
		"",
		"govman_auto_switch() {",
		"  # Check if a .govman-version file exists in the current directory",
		"  if [ -f \".govman-version\" ]; then",
		"    local version",
		"    # Read the version from the .govman-version file and remove any carriage return characters",
		"    version=$(<.govman-version)",
		"    version=\"${version//$'\\r'/}\"",
		"",
		"    # If a version is specified, try to switch to it",
		"    if [ -n \"$version\" ]; then",
		"      local output",
		"      # Run 'govman use <version>' and capture the output",
		"      output=$(command govman use \"$version\" 2>/dev/null)",
		"      # Extract only the first line (the export command) if it starts with 'export PATH='",
		"      if [[ $? -eq 0 && \"$output\" == export\\ PATH=* ]]; then",
		"        local export_cmd=$(echo \"$output\" | head -n 1)",
		"        eval \"$export_cmd\"",
		"      fi",
		"    fi",
		"  else",
		"    # If no .govman-version file found, switch to the default version",
		"    local output",
		"    output=$(command govman use default 2>/dev/null)",
		"    if [[ $? -eq 0 && \"$output\" == export\\ PATH=* ]]; then",
		"      local export_cmd=$(echo \"$output\" | head -n 1)",
		"      eval \"$export_cmd\"",
		"    fi",
		"  fi",
		"}",
		"",
		"# Production-safe govman wrapper - handles use and refresh commands",
		"govman() {",
		"  local cmd=\"$1\"",
		"  ",
		"  # Handle commands that output shell commands",
		"  if [[ \"$cmd\" == \"use\" || \"$cmd\" == \"refresh\" ]]; then",
		"    local output ret",
		"    output=$(command govman \"$@\")",
		"    ret=$?",
		"    ",
		"    if [[ $ret -eq 0 ]]; then",
		"      # Check if output contains an export command",
		"      if [[ \"$output\" == export\\ PATH=* ]]; then",
		"        # Extract and execute the export command safely",
		"        local export_cmd=$(echo \"$output\" | head -n 1)",
		"        eval \"$export_cmd\" 2>/dev/null || {",
		"          echo \"Warning: Failed to update PATH\" >&2",
		"        }",
		"        # Show the rest of the output (status messages)",
		"        echo \"$output\" | tail -n +2",
		"      else",
		"        # No export command, just show the output",
		"        echo \"$output\"",
		"      fi",
		"    else",
		"      # Command failed, show error output",
		"      echo \"$output\" >&2",
		"      return $ret",
		"    fi",
		"  else",
		"    # For other govman commands, just run them normally",
		"    command govman \"$@\"",
		"  fi",
		"}",
		"",
		"# Override the 'cd' command to automatically switch Go version on directory change",
		"cd() {",
		"  if builtin cd \"$@\"; then",
		"    govman_auto_switch",
		"    return 0",
		"  else",
		"    return 1",
		"  fi",
		"}",
		"",
		"# Run once at shell startup to set the PATH to the correct Go version",
		"govman_auto_switch",
		"# END GOVMAN",
	}
}

// ZshShell implementation
func (s *ZshShell) Name() string {
	return "zsh"
}

func (s *ZshShell) DisplayName() string {
	return "Zsh"
}

func (s *ZshShell) IsAvailable() bool {
	return isCommandAvailable("zsh")
}

func (s *ZshShell) ConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Use current directory as fallback instead of /tmp
		home = "."
	}
	return filepath.Join(home, ".zshrc")
}

func (s *ZshShell) PathCommand(path string) string {
	return fmt.Sprintf(`export PATH="%s:$PATH"`, path)
}

func (s *ZshShell) SetupCommands(binPath string) []string {
	return []string{
		"# GOVMAN - Go Version Manager",
		"export PATH=\"$HOME/.govman/bin:$PATH\"",
		"",
		"govman_auto_switch() {",
		"  # Check if a .govman-version file exists in the current directory",
		"  if [ -f \".govman-version\" ]; then",
		"    local version",
		"    # Read the version from the .govman-version file and remove any carriage return characters",
		"    version=$(<.govman-version)",
		"    version=\"${version//$'\\r'/}\"",
		"",
		"    # If a version is specified, try to switch to it",
		"    if [ -n \"$version\" ]; then",
		"      local output",
		"      # Run 'govman use <version>' and capture the output",
		"      output=$(command govman use \"$version\" 2>/dev/null)",
		"      # Extract only the first line (the export command) if it starts with 'export PATH='",
		"      if [[ $? -eq 0 && \"$output\" == export\\ PATH=* ]]; then",
		"        local export_cmd=$(echo \"$output\" | head -n 1)",
		"        eval \"$export_cmd\"",
		"      fi",
		"    fi",
		"  else",
		"    # If no .govman-version file found, switch to the default version",
		"    local output",
		"    output=$(command govman use default 2>/dev/null)",
		"    if [[ $? -eq 0 && \"$output\" == export\\ PATH=* ]]; then",
		"      local export_cmd=$(echo \"$output\" | head -n 1)",
		"      eval \"$export_cmd\"",
		"    fi",
		"  fi",
		"}",
		"",
		"# Production-safe govman wrapper - handles use and refresh commands",
		"govman() {",
		"  local cmd=\"$1\"",
		"  ",
		"  # Handle commands that output shell commands",
		"  if [[ \"$cmd\" == \"use\" || \"$cmd\" == \"refresh\" ]]; then",
		"    local output ret",
		"    output=$(command govman \"$@\")",
		"    ret=$?",
		"    ",
		"    if [[ $ret -eq 0 ]]; then",
		"      # Check if output contains an export command",
		"      if [[ \"$output\" == export\\ PATH=* ]]; then",
		"        # Extract and execute the export command safely",
		"        local export_cmd=$(echo \"$output\" | head -n 1)",
		"        eval \"$export_cmd\" 2>/dev/null || {",
		"          echo \"Warning: Failed to update PATH\" >&2",
		"        }",
		"        # Show the rest of the output (status messages)",
		"        echo \"$output\" | tail -n +2",
		"      else",
		"        # No export command, just show the output",
		"        echo \"$output\"",
		"      fi",
		"    else",
		"      # Command failed, show error output",
		"      echo \"$output\" >&2",
		"      return $ret",
		"    fi",
		"  else",
		"    # For other govman commands, just run them normally",
		"    command govman \"$@\"",
		"  fi",
		"}",
		"",
		"# Override the 'cd' command to automatically switch Go version on directory change",
		"cd() {",
		"  if builtin cd \"$@\"; then",
		"    govman_auto_switch",
		"    return 0",
		"  else",
		"    return 1",
		"  fi",
		"}",
		"",
		"# For zsh users: hook into directory change events so pushd/popd also trigger version switching",
		"if [ -n \"$ZSH_VERSION\" ]; then",
		"  autoload -U add-zsh-hook 2>/dev/null",
		"  if type add-zsh-hook >/dev/null 2>&1; then",
		"    add-zsh-hook chpwd govman_auto_switch",
		"  fi",
		"fi",
		"",
		"# Run once at shell startup to set the PATH to the correct Go version",
		"govman_auto_switch",
		"# END GOVMAN",
	}
}

// FishShell implementation
func (s *FishShell) Name() string {
	return "fish"
}

func (s *FishShell) DisplayName() string {
	return "Fish"
}

func (s *FishShell) IsAvailable() bool {
	return isCommandAvailable("fish")
}

func (s *FishShell) ConfigFile() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "fish", "config.fish")
}

func (s *FishShell) PathCommand(path string) string {
	return fmt.Sprintf(`set -gx PATH "%s" $PATH`, path)
}

func (s *FishShell) SetupCommands(binPath string) []string {
	return []string{
		"# GOVMAN - Go Version Manager",
		"set -gx PATH \"$HOME/.govman/bin\" $PATH",
		"",
		"function govman_auto_switch",
		"    # Check if a .govman-version file exists in the current directory",
		"    if test -f .govman-version",
		"        # Read the version and remove any carriage return characters",
		"        set version (cat .govman-version | tr -d '\\n\\r')",
		"        ",
		"        # If a version is specified, try to switch to it",
		"        if test -n \"$version\"",
		"            set output (command govman use \"$version\" 2>/dev/null)",
		"            # Extract only the first line if it starts with 'set -gx PATH'",
		"            if test $status -eq 0; and string match -q 'set -gx PATH*' \"$output\"",
		"                set export_cmd (echo \"$output\" | head -n 1)",
		"                eval $export_cmd",
		"            end",
		"        end",
		"    else",
		"        # If no .govman-version file found, switch to the default version",
		"        set output (command govman use default 2>/dev/null)",
		"        if test $status -eq 0; and string match -q 'set -gx PATH*' \"$output\"",
		"            set export_cmd (echo \"$output\" | head -n 1)",
		"            eval $export_cmd",
		"        end",
		"    end",
		"end",
		"",
		"# Override cd to trigger auto-switching",
		"function cd",
		"    builtin cd $argv",
		"    and govman_auto_switch",
		"end",
		"",
		"# Hook into pwd change for auto-switching",
		"function __govman_cd_hook --on-variable PWD",
		"    govman_auto_switch",
		"end",
		"",
		"# Production-safe govman wrapper - handles use and refresh commands",
		"function govman",
		"    set cmd $argv[1]",
		"    ",
		"    # Handle commands that output shell commands",
		"    if test \"$cmd\" = \"use\"; or test \"$cmd\" = \"refresh\"",
		"        set output (command govman $argv)",
		"        set ret $status",
		"        ",
		"        if test $ret -eq 0",
		"            # Check if output contains a set command",
		"            if string match -q 'set -gx PATH*' \"$output\"",
		"                # Extract and execute the export command safely",
		"                set export_cmd (echo \"$output\" | head -n 1)",
		"                eval $export_cmd 2>/dev/null; or begin",
		"                    echo \"Warning: Failed to update PATH\" >&2",
		"                end",
		"                # Show the rest of the output (status messages)",
		"                echo \"$output\" | tail -n +2",
		"            else",
		"                # No export command, just show the output",
		"                echo \"$output\"",
		"            end",
		"        else",
		"            # Command failed, show error output",
		"            echo \"$output\" >&2",
		"            return $ret",
		"        end",
		"    else",
		"        # For other govman commands, just run them normally",
		"        command govman $argv",
		"    end",
		"end",
		"",
		"# Run once at shell startup",
		"govman_auto_switch",
		"# END GOVMAN",
	}
}

// PowerShell implementation
func (s *PowerShell) Name() string {
	return "powershell"
}

func (s *PowerShell) DisplayName() string {
	return "PowerShell"
}

func (s *PowerShell) IsAvailable() bool {
	return isCommandAvailable("pwsh") || isCommandAvailable("powershell")
}

func (s *PowerShell) ConfigFile() string {
	// Return the profile variable reference
	return "$PROFILE"
}

func (s *PowerShell) PathCommand(path string) string {
	return fmt.Sprintf(`$env:PATH = "%s;" + $env:PATH`, path)
}

func (s *PowerShell) SetupCommands(binPath string) []string {
	return []string{
		"# GOVMAN - Go Version Manager",
		"$env:PATH = \"$env:USERPROFILE\\.govman\\bin;\" + $env:PATH",
		"",
		"function Invoke-GovmanAutoSwitch {",
		"    # Check if a .govman-version file exists in the current directory",
		"    if (Test-Path \".govman-version\") {",
		"        # Read the version and trim whitespace",
		"        $version = (Get-Content \".govman-version\").Trim()",
		"        ",
		"        # If a version is specified, try to switch to it",
		"        if ($version) {",
		"            $output = govman use $version 2>$null",
		"            # Check if the output starts with '$env:PATH ='",
		"            if ($LASTEXITCODE -eq 0 -and $output -like '$env:PATH = *') {",
		"                Invoke-Expression $output",
		"            }",
		"        }",
		"    } else {",
		"        # If no .govman-version file found, switch to the default version",
		"        $output = govman use default 2>$null",
		"        if ($LASTEXITCODE -eq 0 -and $output -like '$env:PATH = *') {",
		"            Invoke-Expression $output",
		"        }",
		"    }",
		"}",
		"",
		"# Override Set-Location to trigger auto-switching",
		"function Set-Location {",
		"    [CmdletBinding()]",
		"    param([string]$Path)",
		"    Microsoft.PowerShell.Management\\Set-Location $Path",
		"    if ($?) {",
		"        Invoke-GovmanAutoSwitch",
		"    }",
		"}",
		"Set-Alias -Name cd -Value Set-Location -Force -Option AllScope",
		"",
		"# Production-safe govman wrapper - handles use and refresh commands",
		"function govman {",
		"    param(",
		"        [Parameter(ValueFromRemainingArguments=$true)]",
		"        [string[]]$Arguments",
		"    )",
		"    ",
		"    $cmd = $Arguments[0]",
		"    ",
		"    # Handle commands that output shell commands",
		"    if ($cmd -eq 'use' -or $cmd -eq 'refresh') {",
		"        try {",
		"            $output = & govman.exe $Arguments",
		"            $ret = $LASTEXITCODE",
		"            ",
		"            if ($ret -eq 0) {",
		"                # Check if output contains an environment command",
		"                if ($output -like '$env:PATH = *') {",
		"                    # Extract and execute the environment command safely",
		"                    $export_cmd = ($output -split \"`n\")[0]",
		"                    try {",
		"                        Invoke-Expression $export_cmd",
		"                    } catch {",
		"                        Write-Warning \"Failed to update PATH: $_\"",
		"                    }",
		"                    # Show the rest of the output (status messages)",
		"                    $remaining = ($output -split \"`n\") | Select-Object -Skip 1",
		"                    if ($remaining) {",
		"                        $remaining | Write-Output",
		"                    }",
		"                } else {",
		"                    # No environment command, just show the output",
		"                    Write-Output $output",
		"                }",
		"            } else {",
		"                # Command failed, show error output",
		"                Write-Error $output",
		"                exit $ret",
		"            }",
		"        } catch {",
		"            Write-Error \"Failed to execute govman: $_\"",
		"            exit 1",
		"        }",
		"    } else {",
		"        # For other govman commands, just run them normally",
		"        & govman.exe $Arguments",
		"    }",
		"}",
		"",
		"# Run once at shell startup",
		"Invoke-GovmanAutoSwitch",
		"# END GOVMAN",
	}
}

// CmdShell implementation (Windows Command Prompt)
func (s *CmdShell) Name() string {
	return "cmd"
}

func (s *CmdShell) DisplayName() string {
	return "Command Prompt"
}

func (s *CmdShell) IsAvailable() bool {
	return runtime.GOOS == "windows"
}

func (s *CmdShell) ConfigFile() string {
	return "Registry/Environment Variables"
}

func (s *CmdShell) PathCommand(path string) string {
	return fmt.Sprintf(`set PATH=%s;%%PATH%%`, path)
}

func (s *CmdShell) SetupCommands(binPath string) []string {
	return []string{
		"REM GOVMAN - Go Version Manager",
		"REM Added by govman installer",
		s.PathCommand(binPath),
		"",
		"REM Basic govman wrapper for Windows Command Prompt",
		"REM Note: Auto-switching requires manual refresh due to CMD limitations",
		"@echo off",
		"doskey govman=govman_wrapper.bat $*",
		"",
		"REM Create govman_wrapper.bat in the same directory as govman.exe",
		"REM This provides basic use and refresh command support",
		"REM Content should be:",
		"REM @echo off",
		"REM set \"cmd=%1\"",
		"REM if \"%cmd%\"==\"use\" (",
		"REM   for /f \"tokens=*\" %%i in ('govman.exe %*') do (",
		"REM     if \"%%i:~0,4\"==\"set \" (",
		"REM       %%i",
		"REM     ) else (",
		"REM       echo %%i",
		"REM     )",
		"REM   )",
		"REM ) else if \"%cmd%\"==\"refresh\" (",
		"REM   for /f \"tokens=*\" %%i in ('govman.exe %*') do (",
		"REM     if \"%%i:~0,4\"==\"set \" (",
		"REM       %%i",
		"REM     ) else (",
		"REM       echo %%i",
		"REM     )",
		"REM   )",
		"REM ) else (",
		"REM   govman.exe %*",
		"REM )",
		"",
		"REM For full functionality, consider using PowerShell instead",
		"REM Run: powershell -Command \"govman init --force\"",
	}
}

// InitializeShell sets up shell integration
func InitializeShell(shell Shell, binPath string, force bool) error {
	switch shell.Name() {
	case "powershell":
		return initializePowerShell(shell, binPath, force)
	case "cmd":
		return initializeCmdShell(shell, binPath, force)
	default:
		return initializeUnixShell(shell, binPath, force)
	}
}

// Initialize Unix-like shells (bash, zsh, fish)
func initializeUnixShell(shell Shell, binPath string, force bool) error {
	configFile := shell.ConfigFile()

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	// Read existing config
	var existingContent string
	if content, err := os.ReadFile(configFile); err == nil {
		existingContent = string(content)
	}

	// Check if govman is already configured
	govmanExists := containsGovmanConfig(existingContent)
	if !force && govmanExists {
		return fmt.Errorf("govman is already configured in %s (use --force to override)", configFile)
	}

	// Remove existing govman configuration if it exists (either forcing or updating)
	if govmanExists {
		existingContent = removeExistingConfig(existingContent)
	}

	// Add govman configuration
	setupCommands := shell.SetupCommands(binPath)
	newConfig := strings.Join(setupCommands, "\n") + "\n"

	// Prepare final content
	finalContent := existingContent
	if existingContent != "" && !strings.HasSuffix(existingContent, "\n") {
		finalContent += "\n"
	}
	finalContent += newConfig

	// Write to config file
	if err := os.WriteFile(configFile, []byte(finalContent), 0644); err != nil {
		return fmt.Errorf("failed to write config to %s: %w", configFile, err)
	}

	return nil
}

// Initialize PowerShell
func initializePowerShell(shell Shell, binPath string, force bool) error {
	fmt.Printf("Setting up %s integration:\n\n", shell.DisplayName())
	fmt.Println("Please add the following to your PowerShell profile:")
	fmt.Printf("(Edit profile: notepad $PROFILE)\n\n")

	commands := shell.SetupCommands(binPath)
	for _, cmd := range commands {
		fmt.Println(cmd)
	}

	fmt.Println("\nAlternatively, run this command in PowerShell as Administrator:")
	fmt.Printf("echo '%s' | Out-File -Append -Encoding UTF8 $PROFILE\n", strings.Join(commands, "`n"))

	return nil
}

// Initialize Windows Command Prompt
func initializeCmdShell(shell Shell, binPath string, force bool) error {
	fmt.Printf("Setting up %s integration:\n\n", shell.DisplayName())

	// Create a batch wrapper for CMD
	wrapperPath := filepath.Join(binPath, "govman_wrapper.bat")
	wrapperContent := `@echo off
set "cmd=%1"
if "%cmd%"=="use" (
  for /f "tokens=*" %%i in ('govman.exe %*') do (
    if "%%i:~0,4"=="set " (
      %%i
    ) else (
      echo %%i
    )
  )
) else if "%cmd%"=="refresh" (
  for /f "tokens=*" %%i in ('govman.exe %*') do (
    if "%%i:~0,4"=="set " (
      %%i
    ) else (
      echo %%i
    )
  )
) else (
  govman.exe %*
)
`

	// Write the wrapper batch file
	if err := os.WriteFile(wrapperPath, []byte(wrapperContent), 0644); err != nil {
		fmt.Printf("Warning: Could not create wrapper batch file: %v\n", err)
	} else {
		fmt.Printf("✅ Created govman wrapper: %s\n", wrapperPath)
	}

	fmt.Println("\nTo complete setup:")
	fmt.Println("1. Add this to your system PATH:")
	fmt.Printf("   %s\n", binPath)
	fmt.Println("2. Open System Properties -> Advanced -> Environment Variables")
	fmt.Printf("3. Add '%s' to your PATH variable\n", binPath)
	fmt.Println("4. Use 'govman_wrapper' instead of 'govman' for full functionality")
	fmt.Println("\nAlternatively, for better experience, use PowerShell:")
	fmt.Println("- Auto-completion")
	fmt.Println("- Auto-switching based on .govman-version files")
	fmt.Println("- Better error handling")
	fmt.Printf("- Run: powershell -Command \"govman init --force\"\n")

	return nil
}

// Helper function to check if govman config exists
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

// removeExistingConfig removes the existing govman configuration from the content.
func removeExistingConfig(content string) string {
	// More precise regex to find the govman block
	// This matches the exact pattern used by the setup commands:
	// - Line with "GOVMAN - Go Version Manager" (with any comment prefix)
	// - All content until the line with "END GOVMAN"
	// - Includes the END GOVMAN line
	regex := `(?m)^[#\s]*GOVMAN - Go Version Manager[^\n]*\n(?s:.*?)^[#\s]*END GOVMAN[^\n]*(?:\n|$)`
	r := regexp.MustCompile(regex)

	// Replace the found block with an empty string
	cleanedContent := r.ReplaceAllString(content, "")

	// Clean up multiple consecutive newlines that might be left
	cleanedContent = regexp.MustCompile(`\n{3,}`).ReplaceAllString(cleanedContent, "\n\n")

	return strings.TrimSpace(cleanedContent)
}

// GetShellInstructions returns manual setup instructions for a shell
func GetShellInstructions(shell Shell, binPath string) string {
	var instructions strings.Builder

	instructions.WriteString(fmt.Sprintf("Manual setup instructions for %s:\n\n", shell.DisplayName()))
	instructions.WriteString(fmt.Sprintf("1. Edit your %s configuration file:\n", shell.Name()))
	instructions.WriteString(fmt.Sprintf("   %s\n\n", shell.ConfigFile()))
	instructions.WriteString("2. Add the following lines:\n\n")

	commands := shell.SetupCommands(binPath)
	for _, cmd := range commands {
		instructions.WriteString(fmt.Sprintf("   %s\n", cmd))
	}

	instructions.WriteString("\n3. Reload your shell configuration:\n")
	switch shell.Name() {
	case "fish":
		instructions.WriteString("   source ~/.config/fish/config.fish\n")
	case "powershell":
		instructions.WriteString("   . $PROFILE\n")
	default:
		instructions.WriteString(fmt.Sprintf("   source %s\n", shell.ConfigFile()))
	}

	return instructions.String()
}
