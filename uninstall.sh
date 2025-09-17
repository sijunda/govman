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
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m' # No Color

# Constants
REPO="sijunda/govman"
INSTALL_DIR="$HOME/.govman"
BIN_DIR="$INSTALL_DIR/bin"

# Global variables
FORCE_UNINSTALL=false
KEEP_CONFIG=false
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

# Interactive confirmation
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

# Show banner
show_banner() {
    local version="$1"
    
    echo -e "${BOLD}${RED}"
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                    GOVMAN - Go Version Manager Uninstaller                   ║"
    echo "╚══════════════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo -e "Uninstalling ${BOLD}govman $version${NC}..."
    echo
}

# Get current govman version
get_current_version() {
    if command -v govman >/dev/null 2>&1; then
        govman version 2>/dev/null | head -1 | awk '{print $3}' || echo "unknown"
    else
        echo "not installed"
    fi
}

# Find shell config files that contain govman
find_shell_configs() {
    local configs=()
    local possible_configs=(
        "$HOME/.bashrc"
        "$HOME/.bash_profile" 
        "$HOME/.zshrc"
        "$HOME/.profile"
        "$HOME/.config/fish/config.fish"
    )
    
    for config in "${possible_configs[@]}"; do
        if [[ -f "$config" ]] && grep -q "govman" "$config" 2>/dev/null; then
            configs+=("$config")
        fi
    done
    
    printf '%s\n' "${configs[@]}"
}

clean_shell_config() {
    local config_file="$1"
    log_debug "Cleaning shell config: $config_file"

    if [[ ! -f "$config_file" ]]; then
        log_debug "Config file not found: $config_file"
        return
    fi

    cp "$config_file" "${config_file}.govman-uninstall-backup"
    log_debug "Backup created: ${config_file}.govman-uninstall-backup"

    # Use a more robust approach - process the file in multiple passes
    local temp_file="${config_file}.tmp.$"
    
    # First pass: Remove obvious govman-related content including multi-line comments
    awk '
    BEGIN {
        in_govman_block = 0
        skip_line = 0
        govman_comment_block = 0
    }
    {
        line = $0
        skip_line = 0

        # Detect start of GOVMAN comment block
        if (line ~ /^[[:space:]]*#.*GOVMAN.*Go Version Manager/) {
            govman_comment_block = 1
            skip_line = 1
        }
        
        # Continue skipping lines in GOVMAN comment block
        if (govman_comment_block) {
            skip_line = 1
            # End of comment block when we hit a non-comment line or empty line
            if (line !~ /^[[:space:]]*#/ && line ~ /^[[:space:]]*$/) {
                govman_comment_block = 0
            } else if (line ~ /^[[:space:]]*[^#]/ && line !~ /^[[:space:]]*$/) {
                govman_comment_block = 0
                # Check this line again since it might be govman-related
                if (line !~ /govman|\.govman/) {
                    skip_line = 0
                }
            }
        }
        
        # Skip GOVMAN header comments (fallback)
        if (line ~ /^[[:space:]]*#.*GOVMAN/) {
            in_govman_block = 1
            skip_line = 1
        }
        
        # Skip govman-related exports and PATH modifications
        if (line ~ /export.*govman|PATH.*govman|\.govman/) {
            skip_line = 1
        }
        
        # Skip govman commands and completion
        if (line ~ /govman |eval.*govman|command.*govman/) {
            skip_line = 1
        }
        
        # Skip function definitions that contain govman
        if (line ~ /govman.*\(\).*{/) {
            skip_line = 1
            in_govman_block = 1
        }
        
        # If we are in a govman section, keep skipping until we find the end
        if (in_govman_block) {
            if (line ~ /^[[:space:]]*}[[:space:]]*$/) {
                in_govman_block = 0
            }
            skip_line = 1
        }
        
        # Skip empty comment lines that might be left over
        if (line ~ /^[[:space:]]*#[[:space:]]*$/) {
            skip_line = 1
        }
        
        # If this is not a line to skip, print it
        if (!skip_line) {
            print line
        }
    }
    ' "$config_file" > "$temp_file"
    
    # Second pass: Clean up remaining artifacts and orphaned blocks
    awk '
    BEGIN {
        in_if_block = 0
        if_buffer = ""
        if_depth = 0
        in_function = 0
        func_buffer = ""
        func_depth = 0
        prev_line = ""
    }
    {
        current_line = $0
        
        # Track if/fi blocks
        if (current_line ~ /^[[:space:]]*if[[:space:]]/ && !in_if_block && !in_function) {
            in_if_block = 1
            if_buffer = current_line "\n"
            if_depth = 1
            next
        }
        
        if (in_if_block) {
            if_buffer = if_buffer current_line "\n"
            
            if (current_line ~ /^[[:space:]]*if[[:space:]]/) {
                if_depth++
            } else if (current_line ~ /^[[:space:]]*fi[[:space:]]*$/) {
                if_depth--
                if (if_depth <= 0) {
                    # Check if this if block contains anything useful
                    if (if_buffer !~ /govman/ && if_buffer ~ /[a-zA-Z]/) {
                        printf "%s", if_buffer
                    }
                    in_if_block = 0
                    if_buffer = ""
                    if_depth = 0
                }
            }
            next
        }
        
        # Track function blocks
        if (current_line ~ /^[[:space:]]*[a-zA-Z_][a-zA-Z0-9_]*[[:space:]]*\(\)/ && !in_function) {
            in_function = 1
            func_buffer = current_line "\n"
            func_depth = 0
            next
        }
        
        if (in_function) {
            func_buffer = func_buffer current_line "\n"
            
            # Count braces
            brace_open = gsub(/{/, "&", current_line)
            brace_close = gsub(/}/, "&", current_line)
            func_depth += brace_open - brace_close
            
            if (func_depth <= 0 && current_line ~ /^[[:space:]]*}/) {
                # Check if this function contains anything useful (not govman-related)
                if (func_buffer !~ /govman/ && func_buffer ~ /[a-zA-Z].*[a-zA-Z]/) {
                    printf "%s", func_buffer
                }
                in_function = 0
                func_buffer = ""
                func_depth = 0
            }
            next
        }
        
        # Skip orphaned fi statements
        if (current_line ~ /^[[:space:]]*fi[[:space:]]*$/) {
            next
        }
        
        # Skip orphaned closing braces at the start of line
        if (current_line ~ /^[[:space:]]*}[[:space:]]*$/ && prev_line ~ /^[[:space:]]*$/) {
            next
        }
        
        # Skip obviously orphaned comments
        if (current_line ~ /^[[:space:]]*#.*[Hh]ook.*function.*directories[[:space:]]*$/) {
            next
        }
        
        if (current_line ~ /^[[:space:]]*#.*[Ii]nitialize.*completion.*system[[:space:]]*$/) {
            next
        }
        
        if (current_line ~ /^[[:space:]]*#.*[Aa]uto-switch.*Go.*version[[:space:]]*$/) {
            next
        }
        
        if (current_line ~ /^[[:space:]]*#.*[Aa]dded by govman installer/) {
            next
        }
        
        # Print the line if it passed all filters
        print current_line
        prev_line = current_line
    }
    ' "$temp_file" > "${temp_file}.2"
    
    # Third pass: Clean up excessive blank lines and normalize spacing
    awk '
    BEGIN {
        blank_count = 0
        buffer = ""
        line_count = 0
    }
    {
        line_count++
        
        if ($0 ~ /^[[:space:]]*$/) {
            blank_count++
            buffer = buffer $0 "\n"
        } else {
            # If we had blank lines before this content line
            if (blank_count > 0 && line_count > 1) {
                # Only add one blank line between content blocks
                print ""
            }
            blank_count = 0
            buffer = ""
            print $0
        }
    }
    END {
        # Don not add trailing blank lines
    }
    ' "${temp_file}.2" > "${temp_file}.3"
    
    # Fourth pass: Remove any remaining trailing empty lines and normalize file ending
    awk '
    {
        lines[NR] = $0
        if ($0 !~ /^[[:space:]]*$/) {
            last_content_line = NR
        }
    }
    END {
        for (i = 1; i <= last_content_line; i++) {
            print lines[i]
        }
    }
    ' "${temp_file}.3" > "$config_file"
    
    # Clean up temp files
    rm -f "$temp_file" "${temp_file}.2" "${temp_file}.3" 2>/dev/null || true
}

# Remove govman installation
remove_govman_files() {
    log_info "Removing govman installation..."
    
    # Remove binary and installation directory
    if [[ -d "$INSTALL_DIR" ]]; then
        local size="unknown"
        if command -v du >/dev/null 2>&1; then
            size=$(du -sh "$INSTALL_DIR" 2>/dev/null | cut -f1 || echo "unknown")
        fi
        
        rm -rf "$INSTALL_DIR"
        log_success "Removed govman installation directory ($size)"
    else
        log_info "Installation directory not found"
    fi
    
    # Remove any govman binaries in common locations
    local common_paths=("/usr/local/bin/govman" "/usr/bin/govman")
    for path in "${common_paths[@]}"; do
        if [[ -f "$path" ]]; then
            log_debug "Found govman at: $path"
            if [[ -w "$(dirname "$path")" ]]; then
                rm -f "$path"
                log_success "Removed govman from $path"
            else
                log_warning "Cannot remove $path (permission denied)"
                log_info "You may need to run: sudo rm $path"
            fi
        fi
    done
}

# Clean shell configurations
clean_shell_configurations() {
    log_info "Cleaning shell configurations..."
    
    local configs
    configs=($(find_shell_configs))
    
    if [[ ${#configs[@]} -eq 0 ]]; then
        log_info "No shell configurations found with govman entries"
        return
    fi
    
    for config in "${configs[@]}"; do
        log_debug "Processing: $config"
        clean_shell_config "$config"
    done
    
    log_success "Cleaned ${#configs[@]} shell configuration file(s)"
}

# Remove Go versions installed by govman
remove_go_versions() {
    if [[ "$KEEP_CONFIG" == true ]]; then
        log_info "Keeping Go versions (--keep-config specified)"
        return
    fi
    
    local go_install_dir="$HOME/.govman/versions"
    
    if [[ -d "$go_install_dir" ]]; then
        local version_count
        version_count=$(find "$go_install_dir" -maxdepth 1 -type d | wc -l)
        version_count=$((version_count - 1)) # Subtract 1 for the parent directory
        
        if [[ $version_count -gt 0 ]]; then
            log_warning "Found $version_count installed Go version(s)"
            
            if confirm "Remove all installed Go versions?" "n"; then
                rm -rf "$go_install_dir"
                log_success "Removed all Go versions"
            else
                log_info "Go versions preserved"
            fi
        fi
    fi
}

# Parse command line arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --force)
                FORCE_UNINSTALL=true
                shift
                ;;
            --keep-config)
                KEEP_CONFIG=true
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
    echo -e "${BOLD}${CYAN}GOVMAN Uninstallation Script${NC}"
    echo -e "${DIM}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo
    echo -e "${BOLD}USAGE:${NC}"
    echo -e "    ${GREEN}curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/uninstall.sh | bash${NC}"
    echo 
    echo -e "    ${DIM}# Or with options:${NC}"
    echo -e "    ${GREEN}curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/uninstall.sh | bash -s -- [OPTIONS]${NC}"
    echo
    echo -e "${BOLD}OPTIONS:${NC}"
    echo -e "    ${CYAN}--force${NC}                     Force uninstallation without prompts"
    echo -e "    ${CYAN}--keep-config${NC}               Keep Go versions and configuration files"
    echo -e "    ${CYAN}--verbose${NC}                   Enable verbose output"
    echo -e "    ${CYAN}--dry-run${NC}                   Show what would be removed without making changes"
    echo -e "    ${CYAN}--help, -h${NC}                  Show this help message"
    echo
    echo -e "${BOLD}EXAMPLES:${NC}"
    echo -e "    ${DIM}# Uninstall with verbose output${NC}"
    echo -e "    curl -sSL ... | bash -s -- ${YELLOW}--verbose${NC}"
    echo 
    echo -e "    ${DIM}# Force uninstall without prompts${NC}"
    echo -e "    curl -sSL ... | bash -s -- ${YELLOW}--force${NC}"
    echo 
    echo -e "    ${DIM}# Keep installed Go versions${NC}"
    echo -e "    curl -sSL ... | bash -s -- ${YELLOW}--keep-config${NC}"
    echo
    echo -e "${DIM}For more information, visit: ${NC}${BLUE}https://github.com/sijunda/govman${NC}"
}

# Show what will be removed
show_removal_summary() {
    echo -e "${BOLD}${YELLOW}📋 Uninstallation Preview${NC}"
    echo -e "${DIM}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    # Check what exists
    local items_to_remove=()
    local has_items=false
    
    if [[ -d "$INSTALL_DIR" ]]; then
        local size="unknown"
        if command -v du >/dev/null 2>&1; then
            size=$(du -sh "$INSTALL_DIR" 2>/dev/null | cut -f1 || echo "unknown")
        fi
        items_to_remove+=("📁 Installation directory: ${CYAN}$INSTALL_DIR${NC} ($size)")
        has_items=true
    fi
    
    if [[ -f "$BIN_DIR/govman" ]]; then
        items_to_remove+=("🔧 Binary: ${CYAN}$BIN_DIR/govman${NC}")
        has_items=true
    fi
    
    local configs
    configs=($(find_shell_configs))
    if [[ ${#configs[@]} -gt 0 ]]; then
        items_to_remove+=("🐚 Shell configurations: ${CYAN}${#configs[@]} file(s) found${NC}")
        has_items=true
        
        # Show preview of what will be cleaned from each config
        for config in "${configs[@]}"; do
            local govman_lines
            govman_lines=$(grep -c -E "(govman|GOVMAN)" "$config" 2>/dev/null || echo "0")
            items_to_remove+=("   ↳ $config (${govman_lines} govman-related lines)")
        done
    fi
    
    local go_versions_dir="$HOME/.govman/versions"
    if [[ -d "$go_versions_dir" ]] && [[ "$KEEP_CONFIG" != true ]]; then
        local version_count
        version_count=$(find "$go_versions_dir" -maxdepth 1 -type d 2>/dev/null | wc -l)
        version_count=$((version_count - 1))
        if [[ $version_count -gt 0 ]]; then
            items_to_remove+=("🐹 Go versions: ${CYAN}$version_count version(s)${NC}")
            has_items=true
            
            # List the versions
            while IFS= read -r -d '' version_dir; do
                local version_name=$(basename "$version_dir")
                if [[ "$version_name" != "versions" ]]; then
                    items_to_remove+=("   ↳ Go $version_name")
                fi
            done < <(find "$go_versions_dir" -maxdepth 1 -type d -print0 2>/dev/null)
        fi
    fi
    
    if [[ "$has_items" != true ]]; then
        echo -e "${GREEN}ℹ${NC} No govman installation found to remove"
        return 1
    fi
    
    echo -e "\n${BOLD}The following will be removed:${NC}"
    for item in "${items_to_remove[@]}"; do
        echo -e "  $item"
    done
    
    if [[ "$KEEP_CONFIG" == true ]]; then
        echo -e "\n${YELLOW}ℹ${NC} Go versions will be preserved (--keep-config)"
    fi
    
    echo -e "\n${DIM}Note: Configuration files will be backed up before modification${NC}"
    
    return 0
}

# Show completion summary
show_completion_summary() {
    echo
    echo -e "${BOLD}${GREEN}🎉 Uninstallation Complete!${NC}"
    echo -e "${DIM}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    echo -e "\n${BOLD}Final Steps:${NC}"
    echo -e "1️⃣  Restart your shell or run: ${DIM}source ~/.$(basename "${SHELL:-bash}")rc${NC}"
    echo -e "2️⃣  Verify removal: ${DIM}command -v govman${NC} (should return nothing)"
    
    local config_backups
    config_backups=($(find "$HOME" -name "*.govman-uninstall-backup" 2>/dev/null))
    if [[ ${#config_backups[@]} -gt 0 ]]; then
        echo -e "\n${BOLD}Backup Files:${NC}"
        echo -e "📄 Configuration backups created:"
        for backup in "${config_backups[@]}"; do
            echo -e "   ↳ ${DIM}$backup${NC}"
        done
        echo -e "\n${DIM}You can safely remove these backup files once you've verified everything works.${NC}"
    fi
    
    echo -e "\n${BOLD}Resources:${NC}"
    echo -e "📖 Documentation: ${BLUE}https://github.com/$REPO${NC}"
    echo -e "💔 Sorry to see you go! Consider sharing feedback: ${BLUE}https://github.com/$REPO/discussions${NC}"
}

# Main uninstallation
main() {
    parse_arguments "$@"
    
    local current_version
    current_version=$(get_current_version)
    
    # Check if govman is installed
    if [[ "$current_version" == "not installed" ]]; then
        echo -e "\n${GREEN}ℹ${NC} govman is not currently installed"
        
        # Check if installation directory exists
        if [[ -d "$INSTALL_DIR" ]]; then
            log_warning "Found installation directory but govman command not available"
            if confirm "Remove installation directory anyway?" "y"; then
                rm -rf "$INSTALL_DIR"
                log_success "Removed orphaned installation directory"
            fi
        fi
        exit 0
    fi
    
    show_banner "$current_version"
    
    # Show what will be removed
    if ! show_removal_summary; then
        exit 0
    fi
    
    # Confirm uninstallation
    if [[ "$FORCE_UNINSTALL" != true ]]; then
        echo
        if ! confirm "Are you sure you want to uninstall govman?" "n"; then
            log_info "Uninstallation cancelled"
            exit 0
        fi
    fi
    
    # Run uninstallation or dry run
    if [[ "$DRY_RUN" == true ]]; then
        echo -e "\n${BOLD}${YELLOW}🔍 DRY RUN MODE${NC} - No changes will be made\n"
        
        log_info "Would remove govman installation files"
        log_info "Would clean shell configuration files"
        if [[ "$KEEP_CONFIG" != true ]]; then
            log_info "Would remove installed Go versions"
        fi
        
        echo
        log_success "Dry run completed successfully"
        log_info "Run without --dry-run to perform actual uninstallation"
    else
        echo
        # Actual uninstallation
        remove_govman_files
        clean_shell_configurations
        remove_go_versions
        
        show_completion_summary
    fi
}

# Run uninstallation
main "$@"