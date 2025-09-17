#!/bin/bash
# GOVMAN Installation Script for Unix-like systems
# Usage: curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/install.sh | bash

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
REPO="sijunda/govman"
GITHUB_API_URL="https://api.github.com/repos/$REPO/releases/latest"
GITHUB_DOWNLOAD_URL="https://github.com/$REPO/releases/download"
INSTALL_DIR="$HOME/.govman"
BIN_DIR="$INSTALL_DIR/bin"

# Global variables
FORCE_INSTALL=false
SKIP_PATH_UPDATE=false
SKIP_SHELL_INTEGRATION=false
CUSTOM_VERSION=""
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

# Enhanced progress bar function
show_progress() {
    local current=$1
    local total=$2
    local prefix="$3"
    local suffix="$4"
    
    if [[ "$INTERACTIVE" != true ]]; then
        return
    fi
    
    local bar_length=50
    local filled_length=$((current * bar_length / total))
    local empty_length=$((bar_length - filled_length))
    
    local bar=$(printf "%*s" "$filled_length" | tr ' ' '█')
    bar+=$(printf "%*s" "$empty_length" | tr ' ' '░')
    
    local percentage=$((current * 100 / total))
    printf "\r${CYAN}%s${NC} [%s] %3d%% %s" "$prefix" "$bar" "$percentage" "$suffix"
    
    if [[ $current -eq $total ]]; then
        echo
    fi
}

# Spinner function for long running tasks
show_spinner() {
    local pid=$1
    local message="$2"
    local delay=0.1
    local spinstr='⣾⣽⣻⢿⡿⣟⣯⣷'
    
    if [[ "$INTERACTIVE" != true ]]; then
        wait $pid
        return $?
    fi
    
    local i=0
    while kill -0 $pid 2>/dev/null; do
        local char=${spinstr:$((i % ${#spinstr})):1}
        printf "\r${BLUE}%s${NC} %s" "$char" "$message"
        sleep $delay
        ((i++))
    done
    
    wait $pid
    local exit_code=$?
    printf "\r"
    return $exit_code
}

# Interactive confirmation with enhanced styling
confirm() {
    local message="$1"
    local default="${2:-n}"
    
    if [[ "$INTERACTIVE" != true ]] || [[ "$FORCE_INSTALL" == true ]]; then
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

# Enhanced banner with system info
show_banner() {
    local platform="$1"
    local version="$2"
    
    echo -e "${BOLD}${CYAN}"
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                    🚀 GOVMAN - Go Version Manager Installer                  ║"
    echo "╚══════════════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo -e "Installing ${BOLD}govman v$version${NC} for ${CYAN}$platform${NC}..."
    echo
}

# Cleanup function
cleanup() {
    local exit_code=$?
    if [[ "$INTERACTIVE" == true ]]; then
        printf "\r\033[K" # Clear current line
    fi
    [[ -n "${TEMP_DIR:-}" ]] && [[ -d "$TEMP_DIR" ]] && rm -rf "$TEMP_DIR"
    
    if [[ $exit_code -ne 0 ]]; then
        echo
        log_error "Installation failed"
        echo -e "${DIM}For support, visit: https://github.com/$REPO/issues${NC}"
    fi
    
    exit $exit_code
}

trap cleanup EXIT INT TERM

# Detect OS and architecture
detect_platform() {
    log_debug "Detecting system platform..."
    
    local os arch
    case "$(uname -s)" in
        Linux*)   os="linux" ;;
        Darwin*)  os="darwin" ;;
        FreeBSD*) os="freebsd" ;;
        *)        os="unknown" ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        armv7l) arch="arm" ;;
        *) arch="unknown" ;;
    esac

    if [[ "$os" == "unknown" || "$arch" == "unknown" ]]; then
        log_error "Unsupported platform: $(uname -s)/$(uname -m)"
        log_info "Supported platforms: linux/darwin/freebsd with amd64/arm64/arm"
        exit 1
    fi
    
    local platform="${os}-${arch}"
    log_debug "Detected platform: $platform"
    echo "$platform"
}

# Get latest release version
get_latest_version() {
    if [[ -n "$CUSTOM_VERSION" ]]; then
        log_debug "Using custom version: $CUSTOM_VERSION"
        echo "$CUSTOM_VERSION"
        return
    fi

    log_debug "Fetching latest version from GitHub..."
    
    local version
    local max_retries=3
    local retry_count=0
    
    while [[ $retry_count -lt $max_retries ]]; do
        if [[ $retry_count -gt 0 ]]; then
            log_warning "Retry attempt $retry_count/$max_retries..."
            sleep 2
        fi
        
        if command -v jq >/dev/null 2>&1; then
            version=$(curl -s --connect-timeout 10 --max-time 30 "$GITHUB_API_URL" | jq -r '.tag_name' 2>/dev/null | sed 's/^v//')
        else
            version=$(curl -s --connect-timeout 10 --max-time 30 "$GITHUB_API_URL" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/^v//' 2>/dev/null)
        fi

        if [[ -n "$version" && "$version" != "null" && "$version" != "" ]]; then
            log_debug "Latest version: $version"
            echo "$version"
            return
        fi
        
        ((retry_count++))
        log_debug "Failed to fetch version information"
    done
    
    log_error "Failed to get latest version after $max_retries attempts"
    log_info "You can specify a version manually with --version=X.Y.Z"
    exit 1
}

# Download binary
download_binary() {
    local url=$1
    local output_path=$2
    local max_retries=3
    local retry_count=0
    
    log_debug "Download URL: $url"
    
    while [[ $retry_count -lt $max_retries ]]; do
        if [[ $retry_count -gt 0 ]]; then
            log_warning "Download failed, retry attempt $retry_count/$max_retries..."
            sleep $((retry_count * 2))
        fi
        
        local curl_opts=(
            --location
            --fail
            --connect-timeout 30
            --max-time 300
            --retry 2
            --retry-delay 1
            --retry-max-time 600
        )
        
        # Add progress bar for interactive terminals
        if [[ "$INTERACTIVE" == true ]]; then
            curl_opts+=(--progress-bar)
        else
            curl_opts+=(--silent --show-error)
        fi
        
        if curl "${curl_opts[@]}" -o "$output_path" "$url"; then
            # Verify file was downloaded and has content
            if [[ -f "$output_path" && -s "$output_path" ]]; then
                local file_size=$(du -h "$output_path" | cut -f1)
                log_debug "Download completed successfully ($file_size)"
                return 0
            else
                log_warning "Downloaded file is empty or missing"
                rm -f "$output_path"
            fi
        fi
        
        ((retry_count++))
    done
    
    log_error "Failed to download binary after $max_retries attempts"
    log_info "Please check your internet connection and try again"
    return 1
}

# Download and install govman
install_govman() {
    local platform version binary_name download_url temp_dir
    platform=$(detect_platform)
    version=$(get_latest_version)
    binary_name="govman-${platform}"
    download_url="${GITHUB_DOWNLOAD_URL}/v${version}/${binary_name}"

    show_banner "$platform" "$version"

    # Create temporary directory
    if temp_dir=$(mktemp -d 2>/dev/null); then
        :
    else
        temp_dir=$(mktemp -d -t govman)
    fi
    TEMP_DIR="$temp_dir"
    log_debug "Created temporary directory: $temp_dir"

    # Create install directory
    mkdir -p "$BIN_DIR"
    log_debug "Created installation directory: $BIN_DIR"

    # Download binary
    local temp_binary="$temp_dir/govman"
    if ! download_binary "$download_url" "$temp_binary"; then
        return 1
    fi

    # Verify and install binary
    log_info "Installing govman..."
    
    # Check if binary is executable
    if ! file "$temp_binary" | grep -q "executable"; then
        log_warning "Downloaded file may not be a valid executable"
    fi
    
    chmod +x "$temp_binary"
    
    # Backup existing installation if present
    if [[ -f "$BIN_DIR/govman" ]]; then
        log_debug "Backing up existing installation"
        cp "$BIN_DIR/govman" "$BIN_DIR/govman.backup"
    fi
    
    mv "$temp_binary" "$BIN_DIR/govman"

    # Verify installation
    if [[ -x "$BIN_DIR/govman" ]]; then
        log_success "govman v$version installed successfully"
    else
        log_error "Installation verification failed"
        return 1
    fi
}

# Add to PATH
update_path() {
    [[ "$SKIP_PATH_UPDATE" == true ]] && return 0
    
    log_debug "Configuring PATH environment..."
    
    local shell_rc

    # Detect shell and appropriate RC file
    if [[ -n "${ZSH_VERSION:-}" ]]; then
        shell_rc="$HOME/.zshrc"
        log_debug "Detected Zsh shell"
    elif [[ -n "${BASH_VERSION:-}" ]]; then
        if [[ -f "$HOME/.bashrc" ]]; then
            shell_rc="$HOME/.bashrc"
        else
            shell_rc="$HOME/.bash_profile"
        fi
        log_debug "Detected Bash shell"
    elif [[ "${SHELL:-}" == */fish ]]; then
        shell_rc="$HOME/.config/fish/config.fish"
        mkdir -p "$(dirname "$shell_rc")"
        log_debug "Detected Fish shell"
    else
        shell_rc="$HOME/.profile"
        log_debug "Using default shell profile"
    fi

    log_debug "Shell configuration file: $shell_rc"

    # Check if PATH is already configured
    if [[ -f "$shell_rc" ]] && grep -q 'govman/bin' "$shell_rc"; then
        log_debug "PATH already configured in $shell_rc"
    else
        log_info "Configuring shell environment..."
        
        # Create backup of RC file
        if [[ -f "$shell_rc" ]]; then
            cp "$shell_rc" "${shell_rc}.govman-backup"
            log_debug "Created backup: ${shell_rc}.govman-backup"
        fi
        
        # Add PATH configuration
        {
            echo ""
            echo "# GOVMAN - Go Version Manager"
            echo "# Added by govman installer on $(date)"
            if [[ "$shell_rc" == *"config.fish" ]]; then
                echo 'set -gx PATH $HOME/.govman/bin $PATH'
            else
                echo 'export PATH="$HOME/.govman/bin:$PATH"'
            fi
            echo "# GOVMAN - Go Version Manager"
        } >> "$shell_rc"
        log_success "Shell environment configured"
    fi
}

# Setup shell integration
setup_shell_integration() {
    [[ "$SKIP_SHELL_INTEGRATION" == true ]] && return 0
    
    export PATH="$BIN_DIR:$PATH"

    if command -v govman >/dev/null 2>&1; then
        log_info "Setting up shell integration..."
        
        # Run govman init in background to show spinner
        {
            if govman init >/dev/null 2>&1; then
                echo "success"
            else
                echo "failed"
            fi
        } &
        
        local pid=$!
        if show_spinner $pid "Configuring shell integration"; then
            local result
            wait $pid
            result=$(wait $pid 2>/dev/null && echo "success" || echo "failed")
            
            if [[ "$result" == "success" ]] || [[ $? -eq 0 ]]; then
                log_success "Shell integration configured"
            else
                log_warning "Shell integration setup encountered issues"
                log_info "You can run 'govman init' manually after restarting your shell"
            fi
        fi
    else
        log_warning "govman not found in current PATH"
        log_info "Please restart your shell or run: source ~/.$(basename $SHELL)rc"
    fi
}

# Parse command line arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --force)
                FORCE_INSTALL=true
                shift
                ;;
            --version=*)
                CUSTOM_VERSION="${1#*=}"
                CUSTOM_VERSION="${CUSTOM_VERSION#v}"
                shift
                ;;
            --skip-path)
                SKIP_PATH_UPDATE=true
                shift
                ;;
            --skip-shell-integration)
                SKIP_SHELL_INTEGRATION=true
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
    echo -e "${BOLD}${CYAN}GOVMAN Installation Script${NC}"
    echo -e "${DIM}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo
    echo -e "${BOLD}USAGE:${NC}"
    echo -e "    ${GREEN}curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/install.sh | bash${NC}"
    echo 
    echo -e "    ${DIM}# Or with options:${NC}"
    echo -e "    ${GREEN}curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/install.sh | bash -s -- [OPTIONS]${NC}"
    echo
    echo -e "${BOLD}OPTIONS:${NC}"
    echo -e "    ${CYAN}--force${NC}                     Force reinstallation even if already installed"
    echo -e "    ${CYAN}--version=VERSION${NC}           Install specific version instead of latest"
    echo -e "    ${CYAN}--skip-path${NC}                 Skip PATH configuration"
    echo -e "    ${CYAN}--skip-shell-integration${NC}    Skip shell integration setup"
    echo -e "    ${CYAN}--verbose${NC}                   Enable verbose output"
    echo -e "    ${CYAN}--dry-run${NC}                   Show what would be done without making changes"
    echo -e "    ${CYAN}--help, -h${NC}                  Show this help message"
    echo
    echo -e "${BOLD}EXAMPLES:${NC}"
    echo -e "    ${DIM}# Install latest version with verbose output${NC}"
    echo -e "    curl -sSL ... | bash -s -- ${YELLOW}--verbose${NC}"
    echo 
    echo -e "    ${DIM}# Install specific version${NC}"
    echo -e "    curl -sSL ... | bash -s -- ${YELLOW}--version=1.2.3${NC}"
    echo 
    echo -e "    ${DIM}# Reinstall forcefully${NC}"
    echo -e "    curl -sSL ... | bash -s -- ${YELLOW}--force${NC}"
    echo
    echo -e "${DIM}For more information, visit: ${NC}${BLUE}https://github.com/sijunda/govman${NC}"
}

show_completion_summary() {
    echo
    echo -e "${BOLD}${GREEN}🎉 Installation Complete!${NC}"
    echo -e "${DIM}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    echo -e "\n${BOLD}Next Steps:${NC}"
    echo -e "1️⃣  Restart your shell or run: ${DIM}source ~/.$(basename $SHELL)rc${NC}"
    echo -e "2️⃣  Verify installation: ${DIM}govman --version${NC}"
    echo -e "3️⃣  Get started: ${DIM}govman --help${NC}"
    
    echo -e "\n${BOLD}Resources:${NC}"
    echo -e "📖 Documentation: ${BLUE}https://github.com/$REPO${NC}"
    echo -e "🐛 Report issues: ${BLUE}https://github.com/$REPO/issues${NC}"
}

# Main installation
main() {
    parse_arguments "$@"
    
    # Pre-installation checks
    if command -v govman >/dev/null 2>&1; then
        local current_version
        current_version=$(govman version 2>/dev/null | head -1 | awk '{print $3}' || echo "unknown")
        
        echo -e "\n${YELLOW}⚠️  EXISTING INSTALLATION DETECTED${NC}"
        echo -e "Current version: ${CYAN}$current_version${NC}"
        echo
        
        if ! confirm "Do you want to reinstall govman?" "n"; then
            log_info "Installation cancelled"
            exit 0
        fi
    fi

    # Check dependencies
    if ! command -v curl >/dev/null 2>&1; then
        log_error "'curl' is required but not installed"
        log_info "Please install curl and try again"
        exit 1
    fi

    # Run installation or dry run
    if [[ "$DRY_RUN" == true ]]; then
        echo -e "\n${BOLD}${YELLOW}🔍 DRY RUN MODE${NC} - No changes will be made\n"
        
        log_info "Would detect platform and fetch version"
        log_info "Would download and install govman binary"
        log_info "Would update PATH configuration"
        log_info "Would setup shell integration"
        
        echo
        log_success "Dry run completed successfully"
        log_info "Run without --dry-run to perform actual installation"
    else
        # Actual installation
        if install_govman && update_path && setup_shell_integration; then
            show_completion_summary
        else
            log_error "Installation failed"
            exit 1
        fi
    fi
}

# Run installation
main "$@"