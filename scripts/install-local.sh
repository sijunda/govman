#!/usr/bin/env bash
 set -e
 QUIET_MODE=false
SPECIFIC_VERSION=""
LOCAL_BINARY_PATH=""
 # ANSI colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
GRAY='\033[0;90m'
NC='\033[0m'
BOLD='\033[1m'
DIM='\033[2m'
 CHECKMARK="✓"
INFO="ℹ"
WARNING="⚠"
CROSSMARK="✗"
INSTALL="📦"
ARROW="→"
 print_info() {
    [[ "$QUIET_MODE" == "true" ]] && return
    echo -e "${BLUE}${BOLD} ${INFO}  INFO${NC} ${GRAY}│${NC} $1"
}
 print_success() {
    [[ "$QUIET_MODE" == "true" ]] && return
    echo -e "${GREEN}${BOLD} ${CHECKMARK}  SUCCESS${NC} ${GRAY}│${NC} $1"
}
 print_error() {
    echo -e "${RED}${BOLD} ${CROSSMARK}  ERROR${NC} ${GRAY}│${NC} $1"
}
 print_step() {
    [[ "$QUIET_MODE" == "true" ]] && return
    echo -e "${PURPLE}${BOLD} ${ARROW}  STEP${NC} ${GRAY}│${NC} $1"
}
 show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --quiet, -q             Run in quiet mode"
    echo "  --version, -v <ver>     Specify version label (e.g., v1.0.0)"
    echo "  --local-binary <path>   Path to local govman binary"
    echo "  --help, -h              Show this help message"
    echo ""
    echo "Example:"
    echo "  $0 --version v1.0.0 --local-binary ./build/govman"
}
 parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --quiet|-q)
                QUIET_MODE=true
                shift
                ;;
            --version|-v)
                SPECIFIC_VERSION="$2"
                shift 2
                ;;
            --local-binary)
                LOCAL_BINARY_PATH="$2"
                shift 2
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}
 detect_shell_config() {
    local shell_name=""
    local config_file=""
    case "$(basename "$SHELL")" in
        zsh) shell_name="zsh"; config_file="$HOME/.zshrc" ;;
        bash)
            shell_name="bash"
            if [[ "$OSTYPE" == "darwin"* ]]; then
                config_file="$HOME/.bash_profile"
            else
                config_file="$HOME/.bashrc"
            fi
            ;;
        fish) shell_name="fish"; config_file="$HOME/.config/fish/config.fish" ;;
        *) shell_name="sh"; config_file="$HOME/.profile" ;;
    esac
    echo "${shell_name}:${config_file}"
}
 verify_binary() {
    local binary_path="$1"
    if [[ ! -f "$binary_path" ]]; then
        print_error "Binary does not exist at: $binary_path"
        return 1
    fi
    if [[ ! -x "$binary_path" ]]; then
        print_error "Binary is not executable: $binary_path"
        return 1
    fi
    if ! "$binary_path" --version >/dev/null 2>&1; then
        print_error "Binary failed to run: $binary_path"
        return 1
    fi
    print_success "Binary verification successful"
    return 0
}
 install_local_binary() {
    local binary_path="$1"
    local install_dir="$HOME/.govman/bin"
    local target_path="$install_dir/govman"
     print_step "Installing local binary..."
    mkdir -p "$install_dir"
     cp "$binary_path" "$target_path"
    chmod +x "$target_path"
     verify_binary "$target_path"
}
 add_to_path() {
    local shell_info
    shell_info=$(detect_shell_config)
    local shell_name=$(echo "$shell_info" | cut -d':' -f1)
    local config_file=$(echo "$shell_info" | cut -d':' -f2)
    local export_line='export PATH="$HOME/.govman/bin:$PATH"'
    local marker="# GOVMAN - Go Version Manager"
     print_step "Updating PATH in $config_file"
     if grep -q "$marker" "$config_file" 2>/dev/null; then
        print_info "govman PATH already set in $config_file"
    else
        {
            echo ""
            echo "$marker"
            echo "$export_line"
        } >> "$config_file"
        print_success "PATH updated in $config_file"
    fi
}
 show_completion() {
    local restart_instr
    shell_info=$(detect_shell_config)
    shell_config=$(echo "$shell_info" | cut -d':' -f2)
    restart_instr="Run: source $shell_config OR restart terminal"
     echo
    echo -e "${GREEN}${BOLD}🎉 GOVMAN Installed successfully!${NC}"
    echo -e "${WHITE}Binary path:${NC} $HOME/.govman/bin/govman"
    echo -e "${WHITE}Version label:${NC} ${BOLD}$SPECIFIC_VERSION${NC}"
    echo -e "${WHITE}Next steps:${NC}"
    echo -e "  $restart_instr"
    echo -e "  Run ${CYAN}govman --version${NC} to confirm"
    echo
}
 main() {
    parse_arguments "$@"
     # ⬇️ Set default version jika belum ditentukan
    if [[ -z "$SPECIFIC_VERSION" ]]; then
        SPECIFIC_VERSION="v1.0.0"
        print_info "No version specified, using default: $SPECIFIC_VERSION"
    fi
     # ⬇️ Set default path ke binary jika belum ditentukan
    if [[ -z "$LOCAL_BINARY_PATH" ]]; then
        LOCAL_BINARY_PATH="./build/govman"
        print_info "No binary path specified, using default: $LOCAL_BINARY_PATH"
    fi
     install_local_binary "$LOCAL_BINARY_PATH"
    add_to_path
    show_completion
}
 trap 'echo -e "\n${RED}Installation interrupted.${NC}"; exit 1' INT TERM
main "$@"
