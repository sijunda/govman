#!/bin/bash

# govman uninstallation script
# This script removes govman from $HOME/.govman/bin and removes it from PATH

set -e

# Enhanced colors and styles
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
GRAY='\033[0;90m'
NC='\033[0m'

# Style effects
BOLD='\033[1m'
DIM='\033[2m'
UNDERLINE='\033[4m'
BLINK='\033[5m'

# Unicode characters for better UI
CHECKMARK="✓"
CROSSMARK="✗"
ARROW="→"
TRASH="🗑"
WARNING="⚠"
QUESTION="❓"
STOP="🛑"
CLEAN="🧹"
SHIELD="🛡"
INFO="ℹ"

# Terminal width detection
TERM_WIDTH=$(tput cols 2>/dev/null || echo 80)

# Print separator line
print_separator() {
    local char="${1:-─}"
    printf "${GRAY}%*s${NC}\n" $TERM_WIDTH | tr ' ' "$char"
}

# Print fancy header
print_header() {
    clear
    print_separator "═"
    echo
    echo -e "${BOLD}${RED}    ██╗   ██╗███╗   ██╗██╗███╗   ██╗███████╗████████╗ █████╗ ██╗     ██╗     ${NC}"
    echo -e "${BOLD}${RED}    ██║   ██║████╗  ██║██║████╗  ██║██╔════╝╚══██╔══╝██╔══██╗██║     ██║     ${NC}"
    echo -e "${BOLD}${RED}    ██║   ██║██╔██╗ ██║██║██╔██╗ ██║███████╗   ██║   ███████║██║     ██║     ${NC}"
    echo -e "${BOLD}${RED}    ██║   ██║██║╚██╗██║██║██║╚██╗██║╚════██║   ██║   ██╔══██║██║     ██║     ${NC}"
    echo -e "${BOLD}${RED}    ╚██████╔╝██║ ╚████║██║██║ ╚████║███████║   ██║   ██║  ██║███████╗███████╗${NC}"
    echo -e "${BOLD}${RED}     ╚═════╝ ╚═╝  ╚═══╝╚═╝╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚══════╝${NC}"
    echo
    echo -e "${BOLD}${WHITE}                        Go Version Manager Uninstaller${NC}"
    echo -e "${DIM}${GRAY}                         Clean removal with enhanced UI${NC}"
    echo
    print_separator "═"
    echo
}

# Enhanced print functions with icons and styling
print_info() {
    echo -e "${BLUE}${BOLD} ${INFO}  INFO${NC} ${GRAY}│${NC} $1"
}

print_success() {
    echo -e "${GREEN}${BOLD} ${CHECKMARK}  SUCCESS${NC} ${GRAY}│${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}${BOLD} ${WARNING}  WARNING${NC} ${GRAY}│${NC} $1"
}

print_error() {
    echo -e "${RED}${BOLD} ${CROSSMARK}  ERROR${NC} ${GRAY}│${NC} $1"
}

print_step() {
    echo -e "${PURPLE}${BOLD} ${ARROW}  STEP${NC} ${GRAY}│${NC} $1"
}

print_clean() {
    echo -e "${CYAN}${BOLD} ${CLEAN}  CLEANING${NC} ${GRAY}│${NC} $1"
}

print_question() {
    echo -e "${YELLOW}${BOLD} ${QUESTION}  QUESTION${NC} ${GRAY}│${NC} $1"
}

# Animated confirmation prompt
print_confirmation() {
    local message="$1"
    local default="$2"
    
    echo
    print_separator "┄"
    echo -e "${BOLD}${RED} ${STOP}  CONFIRMATION REQUIRED${NC}"
    print_separator "┄"
    echo -e "${YELLOW}$message${NC}"
    print_separator "┄"
}

# Enhanced user input function
get_user_input() {
    local prompt="$1"
    local response=""
    
    echo -e "${BOLD}${WHITE}$prompt${NC}"
    echo -n "   "
    
    # Try to read from /dev/tty if available (works when script is piped)
    if [[ -t 0 ]] || [[ -r /dev/tty ]]; then
        if [[ -r /dev/tty ]]; then
            read -n 1 -r response < /dev/tty
        else
            read -p "" -n 1 -r response
        fi
        echo "" # New line after input
    else
        # Fallback: assume 'no' if no interactive terminal available
        print_warning "No interactive terminal detected. Assuming 'N' (no)."
        response="N"
    fi
    
    echo "$response"
}

# Show what will be removed
show_removal_preview() {
    echo -e "${BOLD}${WHITE}Removal Preview:${NC}"
    print_separator "┄"
    
    local install_dir="$HOME/.govman/bin"
    local govman_dir="$HOME/.govman"
    local shell_configs=("$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zshrc")
    
    # Check binary
    if [[ -d "$install_dir" ]]; then
        echo -e "${RED} ${TRASH}${NC} Binary directory: ${BOLD}$install_dir${NC}"
    else
        echo -e "${GRAY} ${CROSSMARK}${NC} Binary directory: ${DIM}$install_dir (not found)${NC}"
    fi
    
    # Check shell configurations
    local config_found=false
    for shell_config in "${shell_configs[@]}"; do
        if [[ -f "$shell_config" ]] && grep -q "# GOVMAN - Go Version Manager" "$shell_config" 2>/dev/null; then
            echo -e "${RED} ${TRASH}${NC} Shell config: ${BOLD}$shell_config${NC}"
            config_found=true
        fi
    done
    
    if [[ "$config_found" == false ]]; then
        echo -e "${GRAY} ${CROSSMARK}${NC} Shell configs: ${DIM}No govman configuration found${NC}"
    fi
    
    # Check data directory
    if [[ -d "$govman_dir" ]]; then
        local dir_size=$(du -sh "$govman_dir" 2>/dev/null | cut -f1 || echo "unknown")
        echo -e "${YELLOW} ${SHIELD}${NC} Data directory: ${BOLD}$govman_dir${NC} ${DIM}($dir_size)${NC}"
    else
        echo -e "${GRAY} ${CROSSMARK}${NC} Data directory: ${DIM}$govman_dir (not found)${NC}"
    fi
    
    print_separator "┄"
    echo
}

# Animated loading for removal process
show_removal_progress() {
    local item="$1"
    local delay=0.1
    local spinstr='⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏'
    local temp
    
    echo -n "   ${DIM}Removing $item... ${NC}"
    for i in {1..10}; do
        temp=${spinstr#?}
        printf "\r   ${DIM}Removing $item... ${CYAN}%c${NC}" "$spinstr"
        spinstr=$temp${spinstr%"$temp"}
        sleep $delay
    done
    printf "\r   ${GREEN}${CHECKMARK}${NC} Removed $item successfully\n"
}

# Remove binary with enhanced feedback
remove_binary() {
    local install_dir="$HOME/.govman/bin"
    
    print_step "Removing govman binary..."
    
    if [[ -d "$install_dir" ]]; then
        show_removal_progress "binary directory"
        rm -rf "$install_dir"
        print_success "Removed govman binary from $install_dir"
    else
        print_warning "govman binary directory not found at $install_dir"
    fi
}

# Remove from PATH with enhanced feedback
remove_from_path() {
    local shell_configs=("$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zshrc")
    local configs_modified=0
    
    # Add fish config if it exists
    if [[ -f "$HOME/.config/fish/config.fish" ]]; then
        shell_configs+=("$HOME/.config/fish/config.fish")
    fi
    
    print_step "Cleaning shell configurations..."
    
    for shell_config in "${shell_configs[@]}"; do
        if [[ -f "$shell_config" ]]; then
            # Check if govman is configured in this config
            if grep -q "# GOVMAN - Go Version Manager" "$shell_config" 2>/dev/null; then
                show_removal_progress "$(basename "$shell_config") configuration"
                
                # Use sed to remove the block between the start and end markers
                sed -i.bak '/# GOVMAN - Go Version Manager/,/# END GOVMAN/d' "$shell_config"
                
                # Clean up extra blank lines that might be left
                awk 'NF || prev_blank {print} {prev_blank = !NF}' "$shell_config" > "${shell_config}.tmp" && mv "${shell_config}.tmp" "$shell_config"

                print_success "Cleaned PATH configuration in $(basename "$shell_config")"
                rm -f "${shell_config}.bak" # Clean up backup file
                ((configs_modified++))
            fi
        fi
    done
    
    if [[ $configs_modified -eq 0 ]]; then
        print_info "No shell configurations found with govman setup"
    else
        print_success "Cleaned $configs_modified shell configuration(s)"
    fi
}

# Remove entire govman directory with enhanced feedback
remove_govman_dir() {
    local govman_dir="$HOME/.govman"
    
    print_step "Removing govman data directory..."
    
    if [[ -d "$govman_dir" ]]; then
        # Show what's being removed
        local dir_size=$(du -sh "$govman_dir" 2>/dev/null | cut -f1 || echo "unknown size")
        print_info "Removing directory: $govman_dir ($dir_size)"
        
        show_removal_progress "data directory"
        rm -rf "$govman_dir"
        print_success "Removed govman data directory"
    else
        print_warning "govman directory not found at $govman_dir"
    fi
}

# Show completion message
show_completion() {
    local complete_removal="$1"
    
    echo
    print_separator "═"
    echo
    if [[ "$complete_removal" == "true" ]]; then
        echo -e "${GREEN}${BOLD} ${CHECKMARK}  COMPLETE UNINSTALLATION SUCCESSFUL!${NC}"
        echo
        print_separator "┄"
        echo -e "${BOLD}${WHITE}What was removed:${NC}"
        echo -e "${GRAY} •${NC} govman binary and executable"
        echo -e "${GRAY} •${NC} Shell PATH configurations"
        echo -e "${GRAY} •${NC} All downloaded Go versions"
        echo -e "${GRAY} •${NC} Complete .govman directory"
    else
        echo -e "${GREEN}${BOLD} ${CHECKMARK}  PARTIAL UNINSTALLATION COMPLETE!${NC}"
        echo
        print_separator "┄"
        echo -e "${BOLD}${WHITE}What was removed:${NC}"
        echo -e "${GRAY} •${NC} govman binary and executable"
        echo -e "${GRAY} •${NC} Shell PATH configurations"
        echo
        echo -e "${BOLD}${WHITE}What was kept:${NC}"
        echo -e "${GRAY} •${NC} Downloaded Go versions in ~/.govman"
    fi
    print_separator "┄"
    echo -e "${BOLD}${WHITE}Final Steps:${NC}"
    echo -e "${GRAY} 1.${NC} Restart your terminal to complete the process"
    echo -e "${GRAY} 2.${NC} Verify with ${RED}govman --version${NC} (should show 'command not found')"
    if [[ "$complete_removal" != "true" ]]; then
        echo -e "${GRAY} 3.${NC} Manually remove ${CYAN}~/.govman${NC} if you change your mind later"
    fi
    print_separator "┄"
    echo -e "${DIM}${GRAY}Thank you for using govman!${NC}"
    print_separator "═"
    echo
}

# Main uninstallation function
main() {
    # Show header
    print_header
    
    print_info "Starting govman uninstallation process..."
    echo
    
    # Show what will be removed
    show_removal_preview
    
    # First confirmation
    print_confirmation "This will remove the govman binary and its configuration from your shell." "N"
    local response
    response=$(get_user_input "Are you sure you want to uninstall govman? ${DIM}(y/N):${NC}")
    
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo
        print_info "Uninstallation cancelled by user"
        print_separator "═"
        echo -e "${DIM}${GRAY}No changes were made to your system.${NC}"
        print_separator "═"
        echo
        exit 0
    fi

    echo
    print_info "Proceeding with govman removal..."
    echo

    # Proceed with basic uninstallation
    remove_binary
    echo
    remove_from_path
    echo

    # Second confirmation for data directory
    print_confirmation "Do you want to remove ALL downloaded Go versions and data?" "N"
    print_warning "This will delete the entire ~/.govman directory and cannot be undone!"
    response=$(get_user_input "Remove data directory permanently? ${DIM}(y/N):${NC}")
    
    echo
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        remove_govman_dir
        show_completion "true"
    else
        print_info "Keeping govman data directory for future use"
        show_completion "false"
    fi
}

# Trap to ensure clean exit
trap 'echo -e "\n${RED}Uninstallation interrupted${NC}"; exit 1' INT TERM

# Run main function
main "$@"