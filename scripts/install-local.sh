#!/bin/bash

# govman local installation script
# This script builds and installs govman from local source code to $HOME/.govman/bin and adds it to PATH

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

# Check if Go is installed
check_go() {
    if ! command -v go >/dev/null 2>&1; then
        print_error "Go is not installed or not in PATH"
        print_info "Please install Go from https://golang.org/dl/"
        exit 1
    fi
    
    local go_version
    go_version=$(go version | cut -d' ' -f3)
    print_info "Found Go version: $go_version"
}

# Build the binary
build_binary() {
    local install_dir="$HOME/.govman/bin"
    local binary_name="govman"
    
    # For Windows
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
        binary_name="govman.exe"
    fi
    
    print_info "Building govman binary..."
    
    # Create install directory
    mkdir -p "$install_dir"
    
    # Build binary
    go build -o "${install_dir}/${binary_name}" ./cmd/govman
    
    # Check if build was successful
    if [[ ! -f "${install_dir}/${binary_name}" ]]; then
        print_error "Failed to build govman binary"
        exit 1
    fi
    
    # Make binary executable (not needed on Windows)
    if [[ "$OSTYPE" != "msys" && "$OSTYPE" != "win32" ]]; then
        chmod +x "${install_dir}/${binary_name}"
    fi
    
    print_success "Built govman binary successfully"
}

# Add to PATH in shell configuration files
# Add to PATH and initialize shell configuration
add_to_path() {
    local install_dir="$HOME/.govman/bin"
    local govman_binary="${install_dir}/govman"

    # Ensure the govman binary is executable
    if [ ! -x "$govman_binary" ]; then
        print_error "govman binary not found or not executable at $govman_binary"
        exit 1
    fi

    print_info "Running 'govman init' to configure your shell..."

    # Run `govman init` and capture its output
    # The `init` command will automatically detect the shell and provide setup instructions
    # We use `--force` to ensure it overwrites any existing configuration
    local init_output
    if init_output=$("$govman_binary" init --force 2>&1); then
        print_success "govman init completed successfully."
        # The output of `govman init` will guide the user if manual steps are needed
        # For Unix-like shells, it will automatically update the config file
        # For PowerShell/cmd, it will print instructions
        echo "$init_output"
    else
        print_error "govman init failed. Please check the output below for details:"
        echo "$init_output"
        print_warning "You may need to run 'govman init' manually to complete the setup."
    fi
}

# Main installation function
main() {
    print_info "Starting govman local installation..."
    
    # Check if Go is installed
    check_go
    
    # Build binary
    build_binary
    
    # Add to PATH
        add_to_path
    
    # Verify installation
    print_info "Verifying installation..."
    local install_dir="$HOME/.govman/bin"
    if "$install_dir/govman" --version >/dev/null 2>&1; then
        print_success "govman installed successfully from local source!"
        print_info "Please restart your terminal or run 'source ~/.bashrc' (or the appropriate config file for your shell)"
        print_info "Then you can use 'govman --help' to get started"
    else
        print_warning "Installation completed, but verification failed"
        print_info "Please restart your terminal and try running 'govman --version'"
    fi
}

# Run main function
main "$@"