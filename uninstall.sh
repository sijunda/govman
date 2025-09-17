#!/bin/bash

# GOVMAN Uninstallation Script for Unix-like systems
# Usage: curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/uninstall.sh | bash

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m' # No Color

# Constants
INSTALL_DIR="$HOME/.govman"
BIN_DIR="$INSTALL_DIR/bin"

# Global variables
FORCE_UNINSTALL=false
VERBOSE=false
DRY_RUN=false
INTERACTIVE=true

# Check if terminal supports interactive features
if [[ ! -t 0 ]] || [[ ! -t 1 ]]; then
    INTERACTIVE=false
fi

# Utility functions
log_debug() {
    if [ "$VERBOSE" = true ]; then
        echo "DEBUG: $1" >&2
    fi
}

log_info() { 
    echo -e "${CYAN}ℹ${NC} $1" >&2; 
}

log_success() { 
    echo -e "${GREEN}✓${NC} $1" >&2; 
}

log_warning() { 
    echo -e "${YELLOW}⚠${NC} $1" >&2; 
}

log_error() { 
    echo -e "${RED}✗${NC} $1" >&2; 
}

# Interactive confirmation with enhanced styling
confirm() {
    local message="$1"
    local default="${2:-n}"
    
    if [[ "$INTERACTIVE" != true ]] || [[ "$FORCE_UNINSTALL" == true ]]; then
        return 0
    fi
    
    local prompt
    if [[ "$default" == "y" ]]; then
        prompt="[${GREEN}Y${NC}/${DIM}n${NC}]"
    else
        prompt="[${DIM}y${NC}/${RED}N${NC}]"
    fi
    
    echo -e "\n${BOLD}❓ $message${NC} $prompt"
    read -rp "   → " REPLY
    echo
    
    if [[ -z "$REPLY" ]]; then
        REPLY="$default"
    fi
    
    [[ $REPLY =~ ^[Yy]$ ]]
}

# Enhanced banner
show_banner() {
    echo -e "${BOLD}${MAGENTA}"
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                    GOVMAN - Go Version Manager Uninstaller                   ║"
    echo "╚══════════════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo -e "Removing ${BOLD}govman${NC} from your system..."
    echo
}

# Cleanup function
cleanup() {
    local exit_code=$?
    if [[ "$INTERACTIVE" == true ]]; then
        printf "\r\033[K" # Clear current line
    fi
    
    if [[ $exit_code -ne 0 ]]; then
        echo
        log_error "Uninstallation failed"
        echo -e "${DIM}For support, visit: https://github.com/sijunda/govman/issues${NC}"
    fi
    
    exit $exit_code
}

trap cleanup EXIT INT TERM

# Detect OS
detect_os() {
    log_debug "Detecting system platform..."
    
    case "$(uname -s)" in
        Linux*)   echo "linux" ;;
        Darwin*)  echo "darwin" ;;
        FreeBSD*) echo "freebsd" ;;
        *)        echo "unknown" ;;
    esac
}

# Check if govman is installed
is_govman_installed() {
    if [[ -d "$INSTALL_DIR" ]] || command -v govman >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Get shell configuration files that might contain govman PATH
get_shell_config_files() {
    local files=()
    
    # Check common shell configuration files
    if [[ -f "$HOME/.zshrc" ]]; then
        files+=("$HOME/.zshrc")
    fi
    
    if [[ -f "$HOME/.bashrc" ]]; then
        files+=("$HOME/.bashrc")
    fi
    
    if [[ -f "$HOME/.bash_profile" ]]; then
        files+=("$HOME/.bash_profile")
    fi
    
    if [[ -f "$HOME/.profile" ]]; then
        files+=("$HOME/.profile")
    fi
    
    # Fish shell
    if [[ -f "$HOME/.config/fish/config.fish" ]]; then
        files+=("$HOME/.config/fish/config.fish")
    fi
    
    # For other shells, check profile
    if [[ ! -f "$HOME/.zshrc" && ! -f "$HOME/.bashrc" && ! -f "$HOME/.bash_profile" ]]; then
        if [[ -f "$HOME/.profile" ]]; then
            files+=("$HOME/.profile")
        fi
    fi
    
    printf '%s\n' "${files[@]}"
}

# Remove PATH modifications from shell configs
remove_path_from_configs() {
    local shell_files
    shell_files=$(get_shell_config_files)

    if [[ -z "$shell_files" ]]; then
        log_debug "No shell configuration files found"
        return 0
    fi

    local modified=false

    while IFS= read -r shell_rc; do
        if [[ ! -f "$shell_rc" ]]; then
            continue
        fi

        log_debug "Checking $shell_rc for govman modifications..."

        # Check if govman is configured
        if grep -q '^# GOVMAN - Go Version Manager' "$shell_rc" 2>/dev/null; then
            log_info "Removing govman configuration block from $shell_rc..."

            if [[ "$DRY_RUN" == true ]]; then
                log_debug "Would remove govman configuration block from $shell_rc"
            else
                # Create backup
                cp "$shell_rc" "${shell_rc}.govman-uninstall-backup"
                log_debug "Created backup: ${shell_rc}.govman-uninstall-backup"

                local temp_file="${shell_rc}.tmp.$$"

                # Get first and last GOVMAN marker lines
                local first_marker_line
                local last_marker_line
                first_marker_line=$(grep -n '^# GOVMAN - Go Version Manager' "$shell_rc" | head -1 | cut -d: -f1)
                last_marker_line=$(grep -n '^# GOVMAN - Go Version Manager' "$shell_rc" | tail -1 | cut -d: -f1)

                # Adjust start line to 1 line before
                local start_line=$((first_marker_line - 1))
                [[ "$start_line" -lt 1 ]] && start_line=1

                log_debug "Removing lines $start_line to $last_marker_line from $shell_rc"

                # Remove GOVMAN block
                sed "${start_line},${last_marker_line}d" "$shell_rc" > "$temp_file"

                # Remove trailing multiple empty lines using awk
                awk 'NF{blank=0} !NF{blank++} blank<2' "$temp_file" > "${temp_file}.cleaned"
                mv "${temp_file}.cleaned" "$temp_file"

                # Ensure ends with single newline
                perl -i -pe 'chomp if eof' "$temp_file" 2>/dev/null || true
                echo "" >> "$temp_file"

                # Replace original
                mv "$temp_file" "$shell_rc"

                # Clean up
                rm -f "${shell_rc}.tmp."* 2>/dev/null || true
            fi

            modified=true
        else
            log_debug "No govman configuration block found in $shell_rc"
        fi
    done <<< "$shell_files"

    if [[ "$modified" == true ]]; then
        log_success "Removed govman configuration from shell files"
    else
        log_info "No govman modifications found in shell configurations"
    fi
}

# Remove govman installation directory
remove_installation() {
    if [[ ! -d "$INSTALL_DIR" ]]; then
        log_info "govman installation directory not found at $INSTALL_DIR"
        return 0
    fi
    
    log_info "Removing govman installation directory..."
    
    if [[ "$DRY_RUN" == true ]]; then
        log_debug "Would remove $INSTALL_DIR"
        log_debug "Contents would include:"
        if command -v find >/dev/null 2>&1; then
            find "$INSTALL_DIR" -maxdepth 2 2>/dev/null || true
        fi
    else
        # Remove the entire installation directory
        rm -rf "$INSTALL_DIR"
        log_success "Removed govman installation directory"
    fi
}

# Remove govman binary from system PATH locations
remove_binary_from_system_paths() {
    local system_paths=(
        "/usr/local/bin/govman"
        "/usr/bin/govman"
        "/bin/govman"
        "$HOME/bin/govman"
        "$HOME/.local/bin/govman"
    )
    
    for path in "${system_paths[@]}"; do
        if [[ -f "$path" ]]; then
            log_info "Removing govman binary from $path..."
            if [[ "$DRY_RUN" == true ]]; then
                log_debug "Would remove $path"
            else
                rm -f "$path"
                log_success "Removed govman binary from $path"
            fi
        fi
    done
}

# Parse command line arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --force)
                FORCE_UNINSTALL=true
                shift
                ;;
            --verbose)
                VERBOSE=true
                shift
                ;;
            --dry-run)
                DRY_RUN=true
                VERBOSE=true
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# Show help information
show_help() {
    echo -e "${BOLD}${MAGENTA}GOVMAN Uninstallation Script${NC}"
    echo -e "${DIM}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo
    echo -e "${BOLD}USAGE:${NC}"
    echo -e "    ${GREEN}curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/uninstall.sh | bash${NC}"
    echo 
    echo -e "    ${DIM}# Or with options:${NC}"
    echo -e "    ${GREEN}curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/uninstall.sh | bash -s -- [OPTIONS]${NC}"
    echo
    echo -e "${BOLD}OPTIONS:${NC}"
    echo -e "    ${CYAN}--force${NC}      Force uninstallation without confirmation"
    echo -e "    ${CYAN}--verbose${NC}    Enable verbose output"
    echo -e "    ${CYAN}--dry-run${NC}    Show what would be done without making changes"
    echo -e "    ${CYAN}--help, -h${NC}   Show this help message"
    echo
    echo -e "${BOLD}EXAMPLES:${NC}"
    echo -e "    ${DIM}# Uninstall with verbose output${NC}"
    echo -e "    curl -sSL ... | bash -s -- ${YELLOW}--verbose${NC}"
    echo 
    echo -e "    ${DIM}# Force uninstallation${NC}"
    echo -e "    curl -sSL ... | bash -s -- ${YELLOW}--force${NC}"
    echo
    echo -e "${DIM}For more information, visit: ${NC}${BLUE}https://github.com/sijunda/govman${NC}"
}

show_completion_summary() {
    echo
    echo -e "${BOLD}${GREEN}🎉 Uninstallation Complete!${NC}"
    echo -e "${DIM}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    echo -e "\n${BOLD}What was removed:${NC}"
    echo -e "🗑️  $INSTALL_DIR directory and all contents"
    echo -e "🗑️  govman PATH entries from shell configurations"
    echo -e "🗑️  govman binaries from system paths"
    
    echo -e "\n${BOLD}Next Steps:${NC}"
    echo -e "1️⃣  Restart your shell or run: ${DIM}source ~/.$(basename $SHELL)rc${NC}"
    echo -e "2️⃣  Verify removal: ${DIM}which govman${NC} (should return nothing)"
    
    echo -e "\n${BOLD}Note:${NC}"
    echo -e "• Backup files were created with .govman-uninstall-backup extension"
    echo -e "• You may want to remove these backup files manually after confirming everything works"
}

# Main uninstallation
main() {
    parse_arguments "$@"
    
    show_banner
    
    # Check if govman is installed
    if ! is_govman_installed; then
        log_warning "govman does not appear to be installed"
        if ! confirm "Do you want to proceed with cleanup anyway?" "n"; then
            log_info "Uninstallation cancelled"
            exit 0
        fi
    else
        echo -e "\n${YELLOW}⚠️  GOVMAN INSTALLATION DETECTED${NC}"
        echo -e "Installation directory: ${CYAN}$INSTALL_DIR${NC}"
        echo
        
        if ! confirm "Are you sure you want to completely remove govman?" "n"; then
            log_info "Uninstallation cancelled"
            exit 0
        fi
    fi
    
    # Run uninstallation or dry run
    if [[ "$DRY_RUN" == true ]]; then
        echo -e "\n${BOLD}${YELLOW}🔍 DRY RUN MODE${NC} - No changes will be made\n"
        
        log_info "Would remove installation directory: $INSTALL_DIR"
        log_info "Would remove PATH modifications from shell configs"
        log_info "Would remove govman binaries from system paths"
        
        echo
        log_success "Dry run completed successfully"
        log_info "Run without --dry-run to perform actual uninstallation"
    else
        # Actual uninstallation
        log_info "Starting govman uninstallation..."
        
        # Remove PATH modifications first (while govman might still be functional)
        remove_path_from_configs
        
        # Remove installation directory
        remove_installation
        
        # Remove binaries from system paths
        remove_binary_from_system_paths
        
        show_completion_summary
    fi
}

# Run uninstallation
main "$@"