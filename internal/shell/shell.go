package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	CompletionCommand() string
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
	home, _ := os.UserHomeDir()
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

func (s *BashShell) CompletionCommand() string {
	return `eval "$(govman completion bash)"`
}

func (s *BashShell) SetupCommands(binPath string) []string {
	return []string{
		"",
		"# GOVMAN - Go Version Manager",
		"# Added by govman installer",
		s.PathCommand(binPath),
		"",
		"# Initialize completion if govman is available",
		"if command -v govman >/dev/null 2>&1; then",
		"  " + s.CompletionCommand(),
		"fi",
		"",
		"# Auto-switch Go version based on .govman-version file",
		"govman_auto_switch() {",
		"  if [ -f .govman-version ]; then",
		"    local version=$(cat .govman-version 2>/dev/null | tr -d '\\n\\r')",
		"    if [ -n \"$version\" ]; then",
		"      local current=$(govman current 2>/dev/null || echo \"none\")",
		"      if [ \"$current\" != \"$version\" ]; then",
		"        echo \"🐹 Switching to Go $version (from .govman-version)\"",
		"        govman use \"$version\"",
		"      fi",
		"    fi",
		"  fi",
		"}",
		"",
		"# Override cd command to trigger auto-switch",
		"govman_original_cd=$(declare -f cd | head -1 | grep -q 'cd is a function' && echo 'function' || echo 'builtin')",
		"cd() {",
		"  if [ \"$govman_original_cd\" = \"function\" ]; then",
		"    command cd \"$@\"",
		"  else",
		"    builtin cd \"$@\"",
		"  fi",
		"  govman_auto_switch",
		"}",
		"",
		"# Trigger auto-switch for current directory on shell startup",
		"govman_auto_switch",
		"# GOVMAN - Go Version Manager",
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
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".zshrc")
}

func (s *ZshShell) PathCommand(path string) string {
	return fmt.Sprintf(`export PATH="%s:$PATH"`, path)
}

func (s *ZshShell) CompletionCommand() string {
	return `eval "$(govman completion zsh)"`
}

func (s *ZshShell) SetupCommands(binPath string) []string {
	return []string{
		"",
		"# GOVMAN - Go Version Manager",
		"# Added by govman installer",
		"",
		"# Initialize zsh completion system if not already done",
		"if [[ ! -f ~/.zcompdump || $(find ~/.zcompdump -mtime +1) ]]; then",
		"  autoload -Uz compinit",
		"  compinit",
		"fi",
		"",
		s.PathCommand(binPath),
		"",
		"# Initialize govman completion",
		"if command -v govman >/dev/null 2>&1; then",
		"  " + s.CompletionCommand(),
		"fi",
		"",
		"# Auto-switch Go version based on .govman-version file",
		"govman_auto_switch() {",
		"  if [[ -f .govman-version ]]; then",
		"    local version=$(cat .govman-version 2>/dev/null | tr -d '\\n\\r')",
		"    if [[ -n \"$version\" ]]; then",
		"      local current=$(govman current 2>/dev/null || echo \"none\")",
		"      if [[ \"$current\" != \"$version\" ]]; then",
		"        echo \"🐹 Switching to Go $version (from .govman-version)\"",
		"        govman use \"$version\"",
		"      fi",
		"    fi",
		"  fi",
		"}",
		"",
		"# Hook function for directory changes",
		"chpwd_functions+=(govman_auto_switch)",
		"",
		"# Trigger auto-switch for current directory on shell startup",
		"govman_auto_switch",
		"# GOVMAN - Go Version Manager",
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

func (s *FishShell) CompletionCommand() string {
	return `govman completion fish | source`
}

func (s *FishShell) SetupCommands(binPath string) []string {
	return []string{
		"",
		"# GOVMAN - Go Version Manager",
		"# Added by govman installer",
		s.PathCommand(binPath),
		"",
		"# Initialize completion",
		"if command -q govman",
		"  " + s.CompletionCommand(),
		"end",
		"",
		"# Auto-switch Go version based on .govman-version file",
		"function govman_auto_switch --on-variable PWD --description 'Auto-switch Go version'",
		"  if test -f .govman-version",
		"    set version (cat .govman-version 2>/dev/null | tr -d '\\n\\r')",
		"    if test -n \"$version\"",
		"      set current (govman current 2>/dev/null; or echo \"none\")",
		"      if test \"$current\" != \"$version\"",
		"        echo \"🐹 Switching to Go $version (from .govman-version)\"",
		"        govman use \"$version\"",
		"      end",
		"    end",
		"  end",
		"end",
		"",
		"# Trigger auto-switch for current directory on shell startup",
		"govman_auto_switch",
		"# GOVMAN - Go Version Manager",
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

func (s *PowerShell) CompletionCommand() string {
	return `govman completion powershell | Out-String | Invoke-Expression`
}

func (s *PowerShell) SetupCommands(binPath string) []string {
	return []string{
		"",
		"# GOVMAN - Go Version Manager",
		"# Added by govman installer",
		s.PathCommand(binPath),
		"",
		"# Initialize completion",
		"if (Get-Command govman -ErrorAction SilentlyContinue) {",
		"  " + s.CompletionCommand(),
		"}",
		"",
		"# Auto-switch Go version based on .govman-version file",
		"function Set-GovmanAutoSwitch {",
		"  if (Test-Path .govman-version) {",
		"    try {",
		"      $version = (Get-Content .govman-version -ErrorAction Stop).Trim()",
		"      if ($version) {",
		"        $current = try { govman current 2>$null } catch { \"none\" }",
		"        if ($current -ne $version) {",
		"          Write-Host \"🐹 Switching to Go $version (from .govman-version)\" -ForegroundColor Green",
		"          govman use $version",
		"        }",
		"      }",
		"    } catch {",
		"      # Silently ignore errors reading .govman-version",
		"    }",
		"  }",
		"}",
		"",
		"# Store original prompt function if it exists",
		"if (Get-Command prompt -ErrorAction SilentlyContinue) {",
		"  if (-not (Get-Variable -Name OriginalPrompt -ErrorAction SilentlyContinue)) {",
		"    $global:OriginalPrompt = Get-Command prompt | Select-Object -ExpandProperty Definition",
		"  }",
		"}",
		"",
		"# Override prompt to trigger auto-switch",
		"function prompt {",
		"  Set-GovmanAutoSwitch",
		"  if ($global:OriginalPrompt) {",
		"    Invoke-Expression $global:OriginalPrompt",
		"  } else {",
		"    \"PS $($executionContext.SessionState.Path.CurrentLocation)$('>' * ($nestedPromptLevel + 1)) \"",
		"  }",
		"}",
		"",
		"# Trigger auto-switch for current directory on startup",
		"Set-GovmanAutoSwitch",
		"# GOVMAN - Go Version Manager",
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

func (s *CmdShell) CompletionCommand() string {
	return "REM Completion not supported in CMD"
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

// Remove existing govman configuration
func removeExistingConfig(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inGovmanSection := false
	skipEmptyLines := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Start of govman section
		if strings.Contains(line, "GOVMAN - Go Version Manager") {
			inGovmanSection = true
			// Skip any preceding empty line
			if len(result) > 0 && strings.TrimSpace(result[len(result)-1]) == "" {
				result = result[:len(result)-1]
			}
			continue
		}

		if inGovmanSection {
			// Skip govman-related lines
			if strings.Contains(line, "govman") ||
				strings.Contains(line, "GOVMAN") ||
				strings.Contains(line, "# Added by govman") ||
				isGovmanFunction(line) {
				continue
			}

			// Skip empty lines immediately after govman section
			if trimmed == "" {
				skipEmptyLines++
				continue
			}

			// End of govman section when we hit non-empty, non-govman content
			if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
				inGovmanSection = false
				skipEmptyLines = 0
			}
		}

		// Add line if not in govman section
		if !inGovmanSection {
			// Add back skipped empty lines (but limit to 1)
			if skipEmptyLines > 0 && trimmed != "" {
				result = append(result, "")
				skipEmptyLines = 0
			}
			result = append(result, line)
		}
	}

	// Clean up trailing empty lines
	for len(result) > 0 && strings.TrimSpace(result[len(result)-1]) == "" {
		result = result[:len(result)-1]
	}

	return strings.Join(result, "\n")
}

// Check if line is part of a govman function
func isGovmanFunction(line string) bool {
	govmanPatterns := []string{
		"govman_auto_switch",
		"Set-GovmanAutoSwitch",
		"chpwd_functions+=(govman_auto_switch)",
		"--on-variable PWD",
		"OriginalPrompt",
		"builtin cd",
		"command cd",
		"# Override cd command",
		"# Hook function",
		"# Auto-switch Go version",
		"# Trigger auto-switch",
	}

	for _, pattern := range govmanPatterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}

	return false
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
