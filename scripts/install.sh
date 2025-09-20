#!/bin/bash

# govman installation script
# This script installs govman to $HOME/.govman/bin and adds it to PATH

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
STAR="★"
ROCKET="🚀"
GEAR="⚙"
DOWNLOAD="📦"
SUCCESS="🎉"
WARNING="⚠"
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
    echo -e "${BOLD}${CYAN}    ██████╗  ██████╗ ██╗   ██╗███╗   ███╗ █████╗ ███╗   ██╗${NC}"
    echo -e "${BOLD}${CYAN}   ██╔════╝ ██╔═══██╗██║   ██║████╗ ████║██╔══██╗████╗  ██║${NC}"
    echo -e "${BOLD}${CYAN}   ██║  ███╗██║   ██║██║   ██║██╔████╔██║███████║██╔██╗ ██║${NC}"
    echo -e "${BOLD}${CYAN}   ██║   ██║██║   ██║╚██╗ ██╔╝██║╚██╔╝██║██╔══██║██║╚██╗██║${NC}"
    echo -e "${BOLD}${CYAN}   ╚██████╔╝╚██████╔╝ ╚████╔╝ ██║ ╚═╝ ██║██║  ██║██║ ╚████║${NC}"
    echo -e "${BOLD}${CYAN}    ╚═════╝  ╚═════╝   ╚═══╝  ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝${NC}"
    echo
    echo -e "${BOLD}${WHITE}                     Go Version Manager Installer${NC}"
    echo -e "${DIM}${GRAY}                   Enhanced with modern UI/UX${NC}"
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

print_download() {
    echo -e "${CYAN}${BOLD} ${DOWNLOAD}  DOWNLOAD${NC} ${GRAY}│${NC} $1"
}

# Animated loading spinner
show_spinner() {
    local pid=$1
    local delay=0.1
    local spinstr='⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏'
    local temp
    echo -n " "
    while kill -0 $pid 2>/dev/null; do
        temp=${spinstr#?}
        printf "\r${CYAN}%c${NC}" "$spinstr"
        spinstr=$temp${spinstr%"$temp"}
        sleep $delay
    done
    printf "\r"
}

# Progress bar function
show_progress() {
    local current=$1
    local total=$2
    local width=40
    local percentage=$((current * 100 / total))
    local completed=$((current * width / total))
    local remaining=$((width - completed))
    
    printf "\r${BOLD}Progress: ${NC}["
    printf "${GREEN}%*s" $completed | tr ' ' '█'
    printf "${GRAY}%*s" $remaining | tr ' ' '░'
    printf "${NC}] ${BOLD}%d%%${NC}" $percentage
}

# System information display
show_system_info() {
    echo -e "${BOLD}${WHITE}System Information:${NC}"
    print_separator "┄"
    echo -e "${GRAY} OS:${NC}           $(uname -s)"
    echo -e "${GRAY} Architecture:${NC} $(uname -m)"
    echo -e "${GRAY} Shell:${NC}        $SHELL"
    echo -e "${GRAY} User:${NC}         $USER"
    echo -e "${GRAY} Home:${NC}         $HOME"
    echo
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

# Detect OS and architecture with enhanced display
detect_platform() {
    local os=""
    local arch=""
    
    print_step "Detecting platform..."
    
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

# Get the latest release version from GitHub with progress
get_latest_version() {
    local version=""
    print_step "Fetching latest version from GitHub..."
    
    # Show spinner for API call
    (
        if command -v curl >/dev/null 2>&1; then
            version=$(curl -s https://api.github.com/repos/sijunda/govman/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
        elif command -v wget >/dev/null 2>&1; then
            version=$(wget -qO- https://api.github.com/repos/sijunda/govman/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
        else
            print_error "Either curl or wget is required to download govman"
            exit 1
        fi
    ) & 
    
    show_spinner $!
    wait
    
    if [[ -z "$version" ]]; then
        print_error "Failed to get latest version information"
        exit 1
    fi
    
    echo "$version"
}

# Enhanced download with progress simulation
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
    
    print_download "Downloading govman ${version} for ${platform}..."
    echo -e "${DIM}${GRAY}   URL: $download_url${NC}"
    echo
    
    # Create install directory
    mkdir -p "$install_dir"
    
    # Simulate download progress (since we can't easily get real progress from curl/wget)
    echo -n "   "
    for i in {1..20}; do
        show_progress $i 20
        sleep 0.05
    done
    echo
    echo
    
    # Download binary
    (
        if command -v curl >/dev/null 2>&1; then
            curl -sSL -o "${install_dir}/${binary_name}" "$download_url"
        elif command -v wget >/dev/null 2>&1; then
            wget -qO "${install_dir}/${binary_name}" "$download_url"
        else
            print_error "Either curl or wget is required to download govman"
            exit 1
        fi
    ) &
    
    show_spinner $!
    wait
    
    # Check if download was successful
    if [[ ! -f "${install_dir}/${binary_name}" ]]; then
        print_error "Failed to download govman binary"
        exit 1
    fi
    
    # Make binary executable
    chmod +x "${install_dir}/${binary_name}"
    print_success "Binary downloaded and made executable"
}

# Enhanced PATH configuration
add_to_path() {
    local install_dir="$1"
    local govman_binary="${install_dir}/govman"

    # Ensure the govman binary is executable
    if [ ! -x "$govman_binary" ]; then
        print_error "govman binary not found or not executable at $govman_binary"
        exit 1
    fi

    print_step "Configuring shell environment..."
    echo -e "${DIM}${GRAY}   Running 'govman init --force'...${NC}"
    echo

    # Run `govman init` and capture its output
    local init_output
    (
        if init_output=$("$govman_binary" init --force 2>&1); then
            echo "$init_output" > /tmp/govman_init_output
        else
            echo "$init_output" > /tmp/govman_init_output
            exit 1
        fi
    ) &
    
    show_spinner $!
    wait
    
    if [ $? -eq 0 ]; then
        print_success "Shell configuration completed"
        local output=$(cat /tmp/govman_init_output)
        if [[ -n "$output" ]]; then
            echo -e "${DIM}${GRAY}$output${NC}"
        fi
    else
        print_error "Shell configuration failed"
        local output=$(cat /tmp/govman_init_output)
        echo "$output"
        print_warning "You may need to run 'govman init' manually to complete the setup"
    fi
    
    # Cleanup temp file
    rm -f /tmp/govman_init_output
}

# Enhanced verification with detailed output
verify_installation() {
    local install_dir="$1"
    
    print_step "Verifying installation..."
    
    if "$install_dir/govman" --version >/dev/null 2>&1; then
        local version_output=$("$install_dir/govman" --version)
        print_success "Installation verified successfully!"
        echo -e "${DIM}${GRAY}   Version: $version_output${NC}"
        return 0
    else
        print_warning "Installation completed, but verification failed"
        print_info "Please restart your terminal and try running 'govman --version'"
        return 1
    fi
}

# Final success message with instructions
show_completion() {
    local restart_instruction="$1"
    
    echo
    print_separator "═"
    echo
    echo -e "${GREEN}${BOLD} ${SUCCESS}  INSTALLATION COMPLETE!${NC}"
    echo
    print_separator "┄"
    echo -e "${BOLD}${WHITE}Next Steps:${NC}"
    echo -e "${GRAY} 1.${NC} $restart_instruction"
    echo -e "${GRAY} 2.${NC} Run ${CYAN}govman --help${NC} to get started"
    echo -e "${GRAY} 3.${NC} Use ${CYAN}govman list${NC} to see available Go versions"
    echo -e "${GRAY} 4.${NC} Use ${CYAN}govman install <version>${NC} to install a Go version"
    print_separator "┄"
    echo -e "${DIM}${GRAY}Thank you for using govman! ${STAR}${NC}"
    print_separator "═"
    echo
}

# Main installation function
main() {
    # Show header
    print_header
    
    # Show system information
    show_system_info
    
    print_info "Starting govman installation process..."
    echo
    
    # Detect platform
    local platform
    platform=$(detect_platform)
    print_success "Platform detected: ${BOLD}$platform${NC}"
    
    # Get latest version
    local version
    version=$(get_latest_version)
    print_success "Latest version found: ${BOLD}$version${NC}"
    
    # Set installation directory
    local install_dir="$HOME/.govman/bin"
    print_info "Installation directory: ${BOLD}$install_dir${NC}"
    echo
    
    # Download binary
    download_binary "$version" "$platform" "$install_dir"
    echo
    
    # Add to PATH
    add_to_path "$install_dir"
    echo
    
    # Verify installation
    if verify_installation "$install_dir"; then
        # Get restart instruction
        local restart_instruction
        restart_instruction=$(get_restart_instruction)
        
        # Show completion message
        show_completion "$restart_instruction"
    else
        echo
        print_separator "═"
        print_warning "Installation completed with warnings"
        print_info "Please restart your terminal and verify with 'govman --version'"
        print_separator "═"
        echo
    fi
}

# Trap to ensure clean exit
trap 'echo -e "\n${RED}Installation interrupted${NC}"; exit 1' INT TERM

# Run main function
main "$@"