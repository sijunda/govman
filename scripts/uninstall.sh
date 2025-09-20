#!/bin/bash

# govman uninstallation script
# This script removes govman from $HOME/.govman/bin and removes it from PATH

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Remove binary
remove_binary() {
    local install_dir="$HOME/.govman/bin"
    
    if [[ -d "$install_dir" ]]; then
        print_info "Removing govman binary from $install_dir"
        rm -rf "$install_dir"
        print_success "Removed govman binary"
    else
        print_warning "govman binary directory not found at $install_dir"
    fi
}

# Remove from PATH in shell configuration files
remove_from_path() {
    local shell_configs=("$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zshrc")
    
    # Add fish config if it exists
    if [[ -f "$HOME/.config/fish/config.fish" ]]; then
        shell_configs+=("$HOME/.config/fish/config.fish")
    fi
    
    for shell_config in "${shell_configs[@]}"; do
        if [[ -f "$shell_config" ]]; then
            # Check if govman is configured in this config
            if grep -q "# GOVMAN - Go Version Manager" "$shell_config" 2>/dev/null; then
                print_info "Removing govman from PATH in $shell_config"
                
                # Use sed to remove the block between the start and end markers
                # This is more robust than the previous awk implementation
                sed -i.bak '/# GOVMAN - Go Version Manager/,/# END GOVMAN/d' "$shell_config"
                
                # Clean up extra blank lines that might be left
                # This awk script removes consecutive blank lines, leaving only one
                awk 'NF || prev_blank {print} {prev_blank = !NF}' "$shell_config" > "${shell_config}.tmp" && mv "${shell_config}.tmp" "$shell_config"

                print_success "Removed govman from PATH in $shell_config"
                rm -f "${shell_config}.bak" # Clean up backup file
            fi
        fi
    done
}

# Remove entire govman directory
remove_govman_dir() {
    local govman_dir="$HOME/.govman"
    
    if [[ -d "$govman_dir" ]]; then
        print_info "Removing entire govman directory"
        rm -rf "$govman_dir"
        print_success "Removed govman directory"
    else
        print_warning "govman directory not found at $govman_dir"
    fi
}

# Main uninstallation function
main() {
    print_info "Starting govman uninstallation..."
    echo ""

    # First, confirm the user wants to uninstall
    print_warning "This will remove the govman binary and its configuration from your shell."
    read -p "Are you sure you want to uninstall govman? (y/N): " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Uninstallation cancelled."
        exit 0
    fi

    # Proceed with uninstallation
    remove_binary
    remove_from_path

    # Second, ask if they want to remove the data directory
    echo ""
    print_warning "Do you also want to remove the entire govman data directory?"
    print_warning "This will delete all downloaded Go versions and cannot be undone."
    read -p "Remove data directory ($HOME/.govman)? (y/N): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        remove_govman_dir
        print_success "govman was completely uninstalled."
    else
        print_success "Uninstalled govman, but the data directory was kept."
    fi

    echo ""
    print_info "Please restart your terminal to complete the uninstallation process."
}

# Run main function
main "$@"