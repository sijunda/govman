#!/bin/bash

# govman installation script
# This script installs govman to $HOME/.govman/bin and adds it to PATH

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

# Check if running on Windows (Git Bash)
is_windows() {
    [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]
}

# Detect current shell and return appropriate config file
detect_shell_config() {
    local shell_name=""
    local config_file=""
    
    # Get the current shell
    if [ -n "$BASH_VERSION" ]; then
        shell_name="bash"
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS uses .bash_profile by default
            config_file="~/.bash_profile"
        else
            config_file="~/.bashrc"
        fi
    elif [ -n "$ZSH_VERSION" ]; then
        shell_name="zsh"
        config_file="~/.zshrc"
    elif [ -n "$FISH_VERSION" ]; then
        shell_name="fish"
        config_file="~/.config/fish/config.fish"
    else
        # Try to detect from $SHELL environment variable
        case "$(basename "$SHELL")" in
            bash)
                shell_name="bash"
                if [[ "$OSTYPE" == "darwin"* ]]; then
                    config_file="~/.bash_profile"
                else
                    config_file="~/.bashrc"
                fi
                ;;
            zsh)
                shell_name="zsh"
                config_file="~/.zshrc"
                ;;
            fish)
                shell_name="fish"
                config_file="~/.config/fish/config.fish"
                ;;
            *)
                shell_name="shell"
                config_file="your shell's configuration file"
                ;;
        esac
    fi
    
    echo "${shell_name}:${config_file}"
}

# Get restart instruction based on OS and shell
get_restart_instruction() {
    local shell_info
    shell_info=$(detect_shell_config)
    local shell_name=$(echo "$shell_info" | cut -d':' -f1)
    local config_file=$(echo "$shell_info" | cut -d':' -f2)
    
    if is_windows; then
        echo "Please restart your terminal or PowerShell window"
    else
        case "$shell_name" in
            bash|zsh)
                echo "Please restart your terminal or run 'source $config_file'"
                ;;
            fish)
                echo "Please restart your terminal or run 'source $config_file'"
                ;;
            *)
                echo "Please restart your terminal or reload your shell configuration"
                ;;
        esac
    fi
}

# Detect OS and architecture
detect_platform() {
    local os=""
    local arch=""
    
    case "$(uname -s)" in
        Linux*)     os=linux;;
        Darwin*)    os=darwin;;
        MINGW*)     os=windows;;
        MSYS*)      os=windows;;
        *)          print_error "Unsupported operating system"; exit 1;;
    esac
    
    case "$(uname -m)" in
        x86_64)     arch=amd64;;
        aarch64)    arch=arm64;;
        arm64)      arch=arm64;;
        armv7l)     arch=arm;;
        i386|i686) arch=386;;
        *)          print_error "Unsupported architecture"; exit 1;;
    esac
    
    # Special case for Windows
    if [[ "$os" == "windows" ]]; then
        arch=amd64  # Default to amd64 for Windows
    fi
    
    echo "${os}/${arch}"
}

# Get the latest release version from GitHub
get_latest_version() {
    local version=""
    if command -v curl >/dev/null 2>&1; then
        version=$(curl -s https://api.github.com/repos/sijunda/govman/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        version=$(wget -qO- https://api.github.com/repos/sijunda/govman/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        print_error "Either curl or wget is required to download govman"
        exit 1
    fi
    
    if [[ -z "$version" ]]; then
        print_error "Failed to get latest version information"
        exit 1
    fi
    
    echo "$version"
}

# Download the binary
download_binary() {
    local version="$1"
    local platform="$2"
    local install_dir="$3"
    
    local os=$(echo "$platform" | cut -d'/' -f1)
    local arch=$(echo "$platform" | cut -d'/' -f2)
    
    # Construct binary name
    local binary_name="govman"
    if [[ "$os" == "windows" ]]; then
        binary_name="govman.exe"
    fi
    
    # Construct download URL
    local download_url="https://github.com/sijunda/govman/releases/download/${version}/govman-${os}-${arch}"
    if [[ "$os" == "windows" ]]; then
        download_url="${download_url}.exe"
    fi
    
    print_info "Downloading govman ${version} for ${platform}..."
    print_info "Download URL: $download_url"
    
    # Create install directory
    mkdir -p "$install_dir"
    
    # Download binary
    if command -v curl >/dev/null 2>&1; then
        curl -sSL -o "${install_dir}/${binary_name}" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "${install_dir}/${binary_name}" "$download_url"
    else
        print_error "Either curl or wget is required to download govman"
        exit 1
    fi
    
    # Check if download was successful
    if [[ ! -f "${install_dir}/${binary_name}" ]]; then
        print_error "Failed to download govman binary"
        exit 1
    fi
    
    # Make binary executable
    chmod +x "${install_dir}/${binary_name}"
}

# Add to PATH and initialize shell configuration
add_to_path() {
    local install_dir="$1"
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
    print_info "Starting govman installation..."
    
    # Detect platform
    local platform
    platform=$(detect_platform)
    print_info "Detected platform: $platform"
    
    # Get latest version
    local version
    version=$(get_latest_version)
    print_info "Latest version: $version"
    
    # Set installation directory
    local install_dir="$HOME/.govman/bin"
    print_info "Installation directory: $install_dir"
    
    # Download binary
    download_binary "$version" "$platform" "$install_dir"
    
    # Add to PATH
    add_to_path "$install_dir"
    
    # Verify installation
    print_info "Verifying installation..."
    if "$install_dir/govman" --version >/dev/null 2>&1; then
        print_success "govman installed successfully!"
        
        # Dynamic restart instruction
        local restart_instruction
        restart_instruction=$(get_restart_instruction)
        print_info "$restart_instruction"
        
        print_info "Then you can use 'govman --help' to get started"
    else
        print_warning "Installation completed, but verification failed"
        print_info "Please restart your terminal and try running 'govman --version'"
    fi
}

# Run main function
main "$@"