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
		"      output=$(govman use \"$version\" 2>/dev/null)",
		"      # If the command succeeded and output starts with 'export PATH=', evaluate it to update the environment",
		"      if [[ $? -eq 0 && \"$output\" == export\\ PATH=* ]]; then",
		"        eval \"$output\"",
		"      fi",
		"    fi",
		"  else",
		"    # If no .govman-version file found, switch to the default version",
		"    local output",
		"    output=$(govman use default 2>/dev/null)",
		"    if [[ $? -eq 0 && \"$output\" == export\\ PATH=* ]]; then",
		"      eval \"$output\"",
		"    fi",
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
		"# Wrap the 'govman use' command so that it automatically evaluates the PATH export output",
		"govman() {",
		"  if [ \"$1\" = \"use\" ]; then",
		"    shift",
		"    local out",
		"    out=$(command govman use \"$@\")",
		"    local ret=$?",
		"    if [ $ret -eq 0 ]; then",
		"      if [[ \"$out\" == export\\ PATH=* ]]; then",
		"        eval \"$out\"",
		"      else",
		"        echo \"$out\"",
		"      fi",
		"    else",
		"      echo \"$out\" >&2",
		"      return $ret",
		"    fi",
		"  else",
		"    # For other govman commands, just run them normally",
		"    command govman \"$@\"",
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
		"      output=$(govman use \"$version\" 2>/dev/null)",
		"      # If the command succeeded and output starts with 'export PATH=', evaluate it to update the environment",
		"      if [[ $? -eq 0 && \"$output\" == export\\ PATH=* ]]; then",
		"        eval \"$output\"",
		"      fi",
		"    fi",
		"  else",
		"    # If no .govman-version file found, switch to the default version",
		"    local output",
		"    output=$(govman use default 2>/dev/null)",
		"    if [[ $? -eq 0 && \"$output\" == export\\ PATH=* ]]; then",
		"      eval \"$output\"",
		"    fi",
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
		"# Wrap the 'govman use' command so that it automatically evaluates the PATH export output",
		"govman() {",
		"  if [ \"$1\" = \"use\" ]; then",
		"    shift",
		"    local out",
		"    out=$(command govman use \"$@\")",
		"    local ret=$?",
		"    if [ $ret -eq 0 ]; then",
		"      if [[ \"$out\" == export\\ PATH=* ]]; then",
		"        eval \"$out\"",
		"      else",
		"        echo \"$out\"",
		"      fi",
		"    else",
		"      echo \"$out\" >&2",
		"      return $ret",
		"    fi",
		"  else",
		"    # For other govman commands, just run them normally",
		"    command govman \"$@\"",
		"  fi",
		"}",
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
		"            set output (govman use \"$version\" 2>/dev/null)",
		"            # Check if the output starts with 'set -gx PATH'",
		"            if test $status -eq 0; and string match -q 'set -gx PATH*' \"$output\"",
		"                eval $output",
		"            end",
		"        end",
		"    else",
		"        # If no .govman-version file found, switch to the default version",
		"        set output (govman use default 2>/dev/null)",
		"        if test $status -eq 0; and string match -q 'set -gx PATH*' \"$output\"",
		"            eval $output",
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
		"# Wrap the 'govman use' command",
		"function govman",
		"    if test \"$argv[1]\" = \"use\"",
		"        set -e argv[1]",
		"        set out (command govman use $argv)",
		"        set ret $status",
		"        if test $ret -eq 0",
		"            if string match -q 'set -gx PATH*' \"$out\"",
		"                eval $out",
		"            else",
		"                echo $out",
		"            end",
		"        else",
		"            echo $out >&2",
		"            return $ret",
		"        end",
		"    else",
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
		"# Wrap the 'govman' command",
		"function govman {",
		"    if ($args[0] -eq 'use') {",
		"        $output = & govman.exe $args",
		"        if ($LASTEXITCODE -eq 0) {",
		"            if ($output -like '$env:PATH = *') {",
		"                Invoke-Expression $output",
		"            } else {",
		"                Write-Output $output",
		"            }",
		"        } else {",
		"            Write-Error $output",
		"            exit $LASTEXITCODE",
		"        }",
		"    } else {",
		"        & govman.exe $args",
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
		"REM Note: Auto-switching and completion not supported in Command Prompt",
		"REM Consider using PowerShell for full govman features",
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
	fmt.Println("Command Prompt has limited integration support.")
	fmt.Println("You need to manually add the following to your system PATH:")
	fmt.Printf("  %s\n\n", binPath)
	fmt.Println("To add to PATH:")
	fmt.Println("1. Open System Properties -> Advanced -> Environment Variables")
	fmt.Printf("2. Add '%s' to your PATH variable\n\n", binPath)
	fmt.Println("For better experience, consider using PowerShell which supports:")
	fmt.Println("- Auto-completion")
	fmt.Println("- Auto-switching based on .govman-version files")
	fmt.Println("- Better error handling")

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
