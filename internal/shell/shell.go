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
		"",
		"# GOVMAN - Go Version Manager",
		"export PATH=\"" + binPath + ":$PATH\"",
		"",
			"# Auto-switch Go version based on .govman-version file",
			"govman_auto_switch() {",
			"    if [ -f \".govman-version\" ]; then",
			"        local version=$(cat .govman-version | tr -d '\\n\\r')",
			"        if [ -n \"$version\" ]; then",
			"            eval \"$(govman use \"$version\" 2>/dev/null)\"",
			"        fi",
			"    else",
			"        # No .govman-version file, revert to default version",
			"        eval \"$(govman use default 2>/dev/null)\"",
			"    fi",
			"}",
			"",
			"# Hook into cd command for auto-switching",
			"cd() {",
			"    builtin cd \"$@\" && govman_auto_switch",
			"}",
		"",
		// "# Enable bash completion for govman",
		// "if command -v govman >/dev/null 2>&1; then",
		// "    eval \"$(govman completion bash)\"",
		// "fi",
		// "",
		"# Initial auto-switch on shell startup",
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
		"",
		"# GOVMAN - Go Version Manager",
		"export PATH=\"" + binPath + ":$PATH\"",
		"",
		"# Auto-switch Go version based on .govman-version file",
		"govman_auto_switch() {",
		"    if [ -f \".govman-version\" ]; then",
		"        local version=$(cat .govman-version | tr -d '\\n\\r')",
		"        if [ -n \"$version\" ]; then",
		"            eval \"$(govman use \"$version\" 2>/dev/null)\"",
		"        fi",
		"    else",
		"        # No .govman-version file, revert to default version",
		"        eval \"$(govman use default 2>/dev/null)\"",
		"    fi",
		"}",
		"",
		"# Hook into cd command for auto-switching",
		"cd() {",
		"    builtin cd \"$@\" && govman_auto_switch",
		"}",
		"",
		"# Hook into chpwd for auto-switching (zsh-specific)",
		"autoload -U add-zsh-hook 2>/dev/null",
		"if command -v add-zsh-hook >/dev/null 2>&1; then",
		"    add-zsh-hook chpwd govman_auto_switch",
		"fi",
		"",
		// "# Enable zsh completion for govman",
		// "if command -v govman >/dev/null 2>&1; then",
		// "    eval \"$(govman completion zsh)\"",
		// "fi",
		// "",
		"# Initial auto-switch on shell startup",
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
		"",
		"# GOVMAN - Go Version Manager",
		"set -gx PATH \"" + binPath + "\" $PATH",
		"",
		"# Auto-switch Go version based on .govman-version file",
		"function govman_auto_switch",
		"    if test -f .govman-version",
		"        set version (cat .govman-version | tr -d '\\n\\r')",
		"        if test -n \"$version\"",
		"            eval (govman use \"$version\" 2>/dev/null)",
		"        end",
		"    else",
		"        # No .govman-version file, revert to default version",
		"        eval (govman use default 2>/dev/null)",
		"    end",
		"end",
		"",
		"# Hook into pwd change for auto-switching",
		"function __govman_cd_hook --on-variable PWD",
		"    govman_auto_switch",
		"end",
		"",
		// "# Enable fish completion for govman",
		// "if command -v govman >/dev/null 2>&1",
		// "    govman completion fish | source",
		// "end",
		// "",
		"# Initial auto-switch on shell startup",
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
		"",
		"# GOVMAN - Go Version Manager",
		"$env:PATH = \"" + binPath + ";\" + $env:PATH",
		"",
		"# Auto-switch Go version based on .govman-version file",
		"function Invoke-GovmanAutoSwitch {",
		"    if (Test-Path \".govman-version\") {",
		"        $version = (Get-Content \".govman-version\").Trim()",
		"        if ($version) {",
		"            Invoke-Expression (govman use $version 2>$null)",
		"        }",
		"    } else {",
		"        # No .govman-version file, revert to default version",
		"        Invoke-Expression (govman use default 2>$null)",
		"    }",
		"}",
		"",
		"# Override Set-Location to trigger auto-switching",
		"function Set-Location {",
		"    [CmdletBinding()]",
		"    param([string]$Path)",
		"    Microsoft.PowerShell.Management\\Set-Location $Path",
		"    Invoke-GovmanAutoSwitch",
		"}",
		"Set-Alias -Name cd -Value Set-Location -Force",
		"",
		// "# Enable PowerShell completion for govman",
		// "if (Get-Command govman -ErrorAction SilentlyContinue) {",
		// "    Invoke-Expression (govman completion powershell | Out-String)",
		// "}",
		// "",
		"# Initial auto-switch on shell startup",
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
	if !force && containsGovmanConfig(existingContent) {
		return fmt.Errorf("govman is already configured in %s (use --force to override)", configFile)
	}

	// Remove existing govman configuration if forcing
	if force {
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
		"govman completion",
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
	// Regex to find the govman block, accommodating different comment styles
	// It looks for a block starting with a line containing "GOVMAN - Go Version Manager"
	// and ending with a line containing "END GOVMAN".
	// The (?s) flag allows . to match newlines.
	// The non-greedy .*? ensures it only matches one block if multiple exist.
	regex := `(?s)(?m)^.*GOVMAN - Go Version Manager.*?\n(?:.*?\n)*?# END GOVMAN\s*\n?`
	r := regexp.MustCompile(regex)

	// Replace the found block with an empty string
	cleanedContent := r.ReplaceAllString(content, "")

	// Trim any leading/trailing whitespace that might be left
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
	if shell.Name() == "fish" {
		instructions.WriteString("   source ~/.config/fish/config.fish\n")
	} else {
		instructions.WriteString(fmt.Sprintf("   source %s\n", shell.ConfigFile()))
	}

	return instructions.String()
}
