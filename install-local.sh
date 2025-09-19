#!/bin/bash
# GOVMAN Installation Script for Unix-like systems
# Usage: curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/install.sh | bash

set -e

# Get the directory of the script
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Constants
REPO="sijunda/govman"
GITHUB_API_URL="https://api.github.com/repos/$REPO/releases/latest"
GITHUB_DOWNLOAD_URL="https://github.com/$REPO/releases/download"
INSTALL_DIR="$HOME/.govman"
BIN_DIR="$INSTALL_DIR/bin"

# Utility functions
log_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

log_success() {
    echo -e "${GREEN}✅${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}⚠️${NC} $1"
}

log_error() {
    echo -e "${RED}❌${NC} $1" >&2
}

# Detect OS and architecture
detect_platform() {
    local os arch
    
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        FreeBSD*)   os="freebsd" ;;
        *)          os="unknown" ;;
    esac
    
    case "$(uname -m)" in
        x86_64)     arch="amd64" ;;
        amd64)      arch="amd64" ;;
        arm64)      arch="arm64" ;;
        aarch64)    arch="arm64" ;;
        armv7l)     arch="arm" ;;
        *)          arch="unknown" ;;
    esac
    
    if [[ "$os" == "unknown" || "$arch" == "unknown" ]]; then
        log_error "Unsupported platform: $(uname -s)/$(uname -m)"
        exit 1
    fi
    
    echo "${os}-${arch}"
}

# Get latest release version
get_latest_version() {
    echo "0.1.0"
}

# Install govman from local build
install_govman() {
    local platform version
    
    platform=$(detect_platform)
    version=$(get_latest_version)
    
    log_info "Installing govman v$version for $platform from local build..."
    
    local local_binary="$SCRIPT_DIR/dist/govman-darwin-arm64"

    if [ ! -f "$local_binary" ]; then
        log_error "Local binary not found at $local_binary"
        exit 1
    fi
    
    # Create installation directory
    mkdir -p "$BIN_DIR"
    
    # Copy binary to final location
    cp "$local_binary" "$BIN_DIR/govman"
    
    # Make executable
    chmod +x "$BIN_DIR/govman"
    
    log_success "govman v$version installed to $BIN_DIR/govman"
}

# Add to PATH
update_path() {
    local shell_rc
    
    # Determine shell configuration file
    if [[ -n "$ZSH_VERSION" ]]; then
        shell_rc="$HOME/.zshrc"
    elif [[ -n "$BASH_VERSION" ]]; then
        if [[ -f "$HOME/.bashrc" ]]; then
            shell_rc="$HOME/.bashrc"
        else
            shell_rc="$HOME/.bash_profile"
        fi
    elif [[ "$SHELL" == */fish ]]; then
        shell_rc="$HOME/.config/fish/config.fish"
        mkdir -p "$(dirname "$shell_rc")"
    else
        shell_rc="$HOME/.profile"
    fi
    
    # Add to PATH if not already there
    if [[ -f "$shell_rc" ]] && grep -q "\.govman/bin" "$shell_rc"; then
        log_info "PATH already configured in $shell_rc"
    else
        log_info "Adding govman to PATH in $shell_rc..."
        
        if [[ "$shell_rc" == *"config.fish" ]]; then
            echo "" >> "$shell_rc"
            echo "# GOVMAN - Go Version Manager" >> "$shell_rc"
            echo 'set -gx PATH $HOME/.govman/bin $PATH' >> "$shell_rc"
        else
            echo "" >> "$shell_rc"
            echo "# GOVMAN - Go Version Manager" >> "$shell_rc"
            echo 'export PATH="$HOME/.govman/bin:$PATH"' >> "$shell_rc"
        fi
        
        log_success "Added govman to PATH in $shell_rc"
    fi
}

# Setup shell integration
setup_shell_integration() {
    log_info "Setting up shell integration..."
    
    # Export PATH for current session
    export PATH="$BIN_DIR:$PATH"
    
    # Initialize shell integration
    if command -v govman >/dev/null 2>&1; then
        govman init
    else
        log_warning "Could not run 'govman init'. Please run it manually after restarting your shell."
    fi
}

# Main installation
main() {
    echo "🚀 GOVMAN - Go Version Manager Installer"
    echo "========================================"
    echo ""
    
    # Check if already installed
    if command -v govman >/dev/null 2>&1; then
        local current_version
        current_version=$(govman version 2>/dev/null | head -1 | cut -d' ' -f3 || echo "unknown")
        log_warning "govman is already installed (version: $current_version)"
        echo ""
        read -p "Do you want to reinstall? [y/N]: " -n 1 -r
        echo ""
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Installation cancelled"
            exit 0
        fi
        echo ""
    fi
    
    # Check dependencies
    if ! command -v curl >/dev/null 2>&1; then
        log_error "curl is required but not installed"
        exit 1
    fi
    
    # Install
    install_govman
    update_path
    setup_shell_integration
    
    echo ""
    log_success "Installation completed successfully!"
    echo ""
    echo "💡 Get help: govman --help"
    echo "🌐 Documentation: https://github.com/$REPO"
}

# Run installation
main "$@"