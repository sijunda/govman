# 📦 Installation Guide

Complete installation instructions for GOVMAN across all supported platforms and shells.

## 🚀 Quick Installation

### **One-Line Install (Recommended)**

#### **Unix/macOS/Linux**
```bash
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

#### **Windows PowerShell**
```powershell
irm https://raw.githubusercontent.com/sijunda/govman/main/install.ps1 | iex
```

### **Alternative Download Methods**

#### **Using wget (Linux)**
```bash
wget -qO- https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

#### **Using Homebrew (macOS/Linux)**
```bash
# Coming soon
brew install govman
```

## 🖥️ Platform-Specific Installation

### **🐧 Linux**

#### **Automatic Installation**
```bash
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

#### **Manual Installation**
1. Download the latest release:
```bash
wget https://github.com/sijunda/govman/releases/latest/download/govman-linux-amd64.tar.gz
tar -xzf govman-linux-amd64.tar.gz
sudo mv govman /usr/local/bin/
```

2. Initialize shell integration:
```bash
govman init
source ~/.bashrc  # or ~/.zshrc, ~/.config/fish/config.fish
```

#### **Package Managers**
```bash
# Ubuntu/Debian (coming soon)
sudo apt install govman

# Arch Linux (AUR)
yay -S govman

# CentOS/RHEL/Fedora (coming soon)
sudo dnf install govman
```

### **🍎 macOS**

#### **Automatic Installation**
```bash
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

#### **Manual Installation**
1. Download for your architecture:
```bash
# Intel Macs
wget https://github.com/sijunda/govman/releases/latest/download/govman-darwin-amd64.tar.gz

# Apple Silicon Macs
wget https://github.com/sijunda/govman/releases/latest/download/govman-darwin-arm64.tar.gz

tar -xzf govman-darwin-*.tar.gz
sudo mv govman /usr/local/bin/
```

2. Initialize shell integration:
```bash
govman init
source ~/.zshrc  # macOS default shell
```

#### **Homebrew (Recommended for macOS)**
```bash
# Add tap (coming soon)
brew tap sijunda/govman
brew install govman
```

### **🪟 Windows**

#### **PowerShell (Recommended)**
```powershell
irm https://raw.githubusercontent.com/sijunda/govman/main/install.ps1 | iex
```

#### **Manual Installation**
1. Download the Windows binary:
   - Visit [GitHub Releases](https://github.com/sijunda/govman/releases)
   - Download `govman-windows-amd64.zip`
   - Extract to a folder (e.g., `C:\tools\govman\`)

2. Add to PATH:
   - Open "Environment Variables" in System Properties
   - Add the govman folder to your PATH

3. Initialize shell integration:
```powershell
govman init
```

#### **Chocolatey (Coming Soon)**
```powershell
choco install govman
```

#### **Scoop**
```powershell
scoop bucket add sijunda https://github.com/sijunda/scoop-bucket
scoop install govman
```

#### **WinGet (Coming Soon)**
```powershell
winget install sijunda.govman
```

## 🐚 Shell-Specific Setup

### **Bash**
```bash
# Install GOVMAN
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

# Initialize shell integration
govman init

# Reload configuration
source ~/.bashrc
```

### **Zsh (macOS default)**
```bash
# Install GOVMAN
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

# Initialize shell integration
govman init

# Reload configuration
source ~/.zshrc
```

### **Fish**
```bash
# Install GOVMAN
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

# Initialize shell integration
govman init

# Reload configuration
source ~/.config/fish/config.fish
```

### **PowerShell**
```powershell
# Install GOVMAN
irm https://raw.githubusercontent.com/sijunda/govman/main/install.ps1 | iex

# Initialize shell integration
govman init

# Reload profile
. $PROFILE
```

### **Command Prompt (Windows)**
```cmd
REM Download and install manually (see Windows section above)
REM Initialize creates a wrapper batch file
govman init

REM Use govman_wrapper for full functionality
govman_wrapper use 1.21.1 --local
```

## 🔧 Advanced Installation Options

### **Custom Installation Directory**
```bash
# Set custom install directory
export GOVMAN_ROOT="$HOME/custom/govman"
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

### **Installation Without Shell Integration**
```bash
# Install binary only
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
# Skip govman init for manual configuration
```

### **Offline Installation**
1. Download release archive on a connected machine
2. Transfer to target machine
3. Extract and place binary in PATH
4. Run `govman init` to set up shell integration

### **Corporate/Enterprise Installation**
```bash
# Download verification
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh -o install.sh
# Review script before execution
cat install.sh
# Execute after review
bash install.sh
```

## 📋 Verification

### **Verify Installation**
```bash
# Check version
govman --version

# Check help
govman --help

# List available commands
govman
```

### **Test Basic Functionality**
```bash
# List remote versions
govman list --remote

# Install a Go version
govman install latest

# Check current version
govman current

# Verify Go is working
go version
```

### **Test Shell Integration**
```bash
# Create test project
mkdir test-project && cd test-project

# Set local version
govman use 1.21.1 --local

# Verify consistency
govman current  # Should show 1.21.1
go version      # Should show go1.21.1

# Test directory switching
cd .. && cd test-project  # Should auto-switch back to 1.21.1
```

## 🔄 Updating GOVMAN

### **Automatic Update**
```bash
govman selfupdate
```

### **Manual Update**
```bash
# Re-run installation script
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

# Update shell integration if needed
govman init --force
```

### **Update Shell Integration Only**
```bash
# Force update shell configuration
govman init --force

# Reload shell
source ~/.bashrc  # or your shell's config file
```

## 🗑️ Uninstallation

### **Complete Removal**
```bash
# Remove all Go versions and GOVMAN
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/uninstall.sh | bash
```

### **Manual Removal**
```bash
# Remove binary
rm /usr/local/bin/govman  # or wherever installed

# Remove data directory
rm -rf ~/.govman

# Remove shell integration (edit config files manually)
# Remove lines between "# GOVMAN" and "# END GOVMAN"
```

## 🆘 Installation Troubleshooting

### **Common Issues**

#### **Permission Denied**
```bash
# If you get permission errors, try:
sudo curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

#### **Shell Integration Not Working**
```bash
# Force reinitialize
govman init --force

# Check shell detection
echo $SHELL

# Manually source configuration
source ~/.bashrc  # or appropriate config file
```

#### **PATH Issues**
```bash
# Check if govman is in PATH
which govman

# Add to PATH manually if needed (add to shell config)
export PATH="$HOME/.govman/bin:$PATH"
```

#### **Download Failures**
```bash
# Try alternative download method
wget -qO- https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

# Or download manually from GitHub releases
```

### **Getting Help**

If you encounter issues:

1. **Check existing issues**: [GitHub Issues](https://github.com/sijunda/govman/issues)
2. **Create new issue**: Include your OS, shell, and error message
3. **Discussion forum**: [GitHub Discussions](https://github.com/sijunda/govman/discussions)

### **Debug Information**
```bash
# Collect debug info for issue reports
govman --version
echo $SHELL
echo $PATH
govman current
```

## 🏢 Enterprise Installation

### **Mass Deployment**
```bash
# Script for multiple machines
#!/bin/bash
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
govman init
govman install 1.21.1
govman use 1.21.1 --default
```

### **Docker Integration**
```dockerfile
# Add to Dockerfile
RUN curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
RUN govman install 1.21.1 && govman use 1.21.1 --default
```

### **CI/CD Integration**
```yaml
# GitHub Actions example
- name: Install GOVMAN
  run: curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

- name: Setup Go
  run: |
    govman install 1.21.1
    govman use 1.21.1
```

---

**Need help?** Check our [Troubleshooting Guide](troubleshooting.md) or [create an issue](https://github.com/sijunda/govman/issues).