# 🆘 Troubleshooting Guide

Complete guide to diagnosing and fixing common GOVMAN issues across all platforms and shells.

## 📋 Table of Contents

- [Quick Diagnostics](#-quick-diagnostics)
- [Installation Issues](#-installation-issues)
- [Version Switching Problems](#-version-switching-problems)
- [Shell Integration Issues](#-shell-integration-issues)
- [Platform-Specific Problems](#-platform-specific-problems)
- [Performance Issues](#-performance-issues)
- [Network and Download Problems](#-network-and-download-problems)
- [Configuration Issues](#-configuration-issues)
- [Getting Help](#-getting-help)

## 🔍 Quick Diagnostics

Before diving into specific issues, run these commands to gather diagnostic information:

### **System Information**
```bash
# GOVMAN version and status
govman --version
govman current
govman list

# System information
echo "OS: $(uname -s)"
echo "Architecture: $(uname -m)"
echo "Shell: $SHELL"
echo "PATH: $PATH"

# Go installation check
which go
go version 2>/dev/null || echo "Go not found in PATH"
```

### **Configuration Check**
```bash
# Check GOVMAN installation
ls -la ~/.govman/
ls -la ~/.govman/bin/
ls -la ~/.govman/versions/

# Check configuration
cat ~/.govman/config.yaml 2>/dev/null || echo "No config file found"

# Check shell integration
grep -n "GOVMAN" ~/.bashrc ~/.zshrc ~/.config/fish/config.fish 2>/dev/null
```

## 🚨 Installation Issues

### **"Command not found: govman"**

#### **Diagnosis**
```bash
# Check if binary exists
ls -la /usr/local/bin/govman
ls -la ~/.govman/bin/govman

# Check PATH
echo $PATH | grep -o '[^:]*govman[^:]*'
```

#### **Solutions**

**Option 1: Reinstall GOVMAN**
```bash
# Re-run installation script
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

**Option 2: Fix PATH manually**
```bash
# Add to shell configuration
echo 'export PATH="$HOME/.govman/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

**Option 3: Manual binary placement**
```bash
# Download and place manually
wget https://github.com/sijunda/govman/releases/latest/download/govman-linux-amd64.tar.gz
tar -xzf govman-linux-amd64.tar.gz
sudo mv govman /usr/local/bin/
```

### **"Permission denied" during installation**

#### **Diagnosis**
```bash
# Check permissions
ls -la /usr/local/bin/
ls -la ~/.govman/
```

#### **Solutions**

**Option 1: Install to user directory**
```bash
# Install to ~/.local/bin instead
mkdir -p ~/.local/bin
export PATH="$HOME/.local/bin:$PATH"
# Re-run installation
```

**Option 2: Fix permissions**
```bash
# Create directories with correct permissions
mkdir -p ~/.govman/bin
chmod 755 ~/.govman/bin
```

**Option 3: Use sudo (if appropriate)**
```bash
# Only if installing system-wide
sudo curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

### **Installation script fails to download**

#### **Diagnosis**
```bash
# Test connectivity
curl -I https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh
wget --spider https://github.com/sijunda/govman/releases/latest
```

#### **Solutions**

**Option 1: Use alternative download method**
```bash
# Try wget instead of curl
wget -qO- https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

**Option 2: Manual download**
```bash
# Download manually and execute
curl -o install.sh https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh
chmod +x install.sh
./install.sh
```

**Option 3: Use proxy or VPN**
```bash
# If behind corporate firewall
export http_proxy=http://proxy.company.com:8080
export https_proxy=http://proxy.company.com:8080
```

## 🔄 Version Switching Problems

### **🎯 The Major Fix: Local Version Inconsistency**

**Problem**: `govman current` and `go version` showing different versions after `govman use --local`.

#### **This is FIXED in the latest version!**

**Before (Broken)** ❌:
```bash
govman use 1.14.1 --local
govman current    # Showed: Go 1.14.1
go version        # Showed: go1.25.1  # INCONSISTENT!
```

**After (Fixed)** ✅:
```bash
govman use 1.14.1 --local
govman current    # Shows: Go 1.14.1
go version        # Shows: go1.14.1   # CONSISTENT!
```

#### **Update to Latest Version**
```bash
# Update GOVMAN
govman selfupdate

# Or reinstall
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

# Update shell integration
govman init --force
source ~/.bashrc  # or your shell config
```

### **Version switch doesn't take effect**

#### **Diagnosis**
```bash
# Check current state
govman current
go version
echo $PATH | tr ':' '\n' | grep go

# Check if shell integration is working
type govman
```

#### **Solutions**

**Option 1: Use govman refresh** ⭐ *NEW*
```bash
# Manual refresh for immediate effect
govman refresh
```

**Option 2: Reload shell configuration**
```bash
# Reload shell
source ~/.bashrc    # Bash
source ~/.zshrc     # Zsh
source ~/.config/fish/config.fish  # Fish
. $PROFILE          # PowerShell
```

**Option 3: Reinitialize shell integration**
```bash
# Force reinitialize
govman init --force
# Restart terminal or reload config
```

### **Auto-switching not working on directory change**

#### **Diagnosis**
```bash
# Check if shell integration is installed
grep -A 10 -B 2 "GOVMAN" ~/.bashrc ~/.zshrc ~/.config/fish/config.fish

# Test .govman-version file
echo "1.21.1" > test-project/.govman-version
cd test-project
govman current
```

#### **Solutions**

**Option 1: Enable shell integration**
```bash
govman init
source ~/.bashrc  # or appropriate config file
```

**Option 2: Check shell compatibility**
```bash
# GOVMAN supports these shells fully:
echo $SHELL
# Should be: bash, zsh, fish, or powershell
```

**Option 3: Manual workaround**
```bash
# Create shell function manually
auto_govman() {
    if [ -f .govman-version ]; then
        govman use $(cat .govman-version)
    fi
}

# Add to shell config and use manually
auto_govman
```

## 🐚 Shell Integration Issues

### **Shell integration overwrites existing configuration**

#### **This is FIXED! GOVMAN now uses production-safe integration.**

**Fixed Issues**:
- ✅ No more PROMPT_COMMAND override
- ✅ No more DEBUG trap conflicts
- ✅ No more global variable pollution
- ✅ Safe error handling

#### **Update to Latest Version**
```bash
govman selfupdate
govman init --force
```

### **Shell functions not working properly**

#### **Diagnosis**
```bash
# Check function definition
type govman
type govman_auto_switch

# Check for errors
govman use 1.21.1 2>&1
```

#### **Solutions**

**Option 1: Reinstall shell integration**
```bash
govman init --force
source ~/.bashrc
```

**Option 2: Manual function check**
```bash
# Test if wrapper is working
command govman use 1.21.1
# vs
govman use 1.21.1
```

**Option 3: Shell-specific fixes**

**Bash/Zsh**:
```bash
# Check for syntax errors
bash -n ~/.bashrc
zsh -n ~/.zshrc
```

**Fish**:
```bash
# Check Fish configuration
fish -n -c 'source ~/.config/fish/config.fish'
```

**PowerShell**:
```powershell
# Check PowerShell profile
Test-Path $PROFILE
Get-Content $PROFILE | Select-String "govman"
```

## 🖥️ Platform-Specific Problems

### **macOS Issues**

#### **"govman" cannot be opened because the developer cannot be verified**

**Solution**:
```bash
# Remove quarantine attribute
sudo xattr -r -d com.apple.quarantine /usr/local/bin/govman

# Or allow in System Preferences > Security & Privacy
```

#### **Homebrew conflicts**

**Solution**:
```bash
# Check for conflicts
brew list | grep go
which go

# Use GOVMAN instead of Homebrew Go
govman use latest --default
```

### **Windows Issues**

#### **PowerShell execution policy errors**

**Solution**:
```powershell
# Allow script execution (run as Administrator)
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser

# Or bypass for specific command
powershell -ExecutionPolicy Bypass -Command "govman use 1.21.1"
```

#### **Windows Defender blocking downloads**

**Solution**:
1. Add GOVMAN directory to Windows Defender exclusions
2. Use alternative download method:
```powershell
# Download manually from GitHub releases
Invoke-WebRequest -Uri "https://github.com/sijunda/govman/releases/latest/download/govman-windows-amd64.zip" -OutFile "govman.zip"
```

#### **Command Prompt limitations**

**Solution**:
```cmd
REM Use PowerShell instead for full functionality
powershell -Command "govman use 1.21.1 --local"

REM Or use the wrapper (after govman init)
govman_wrapper use 1.21.1 --local
```

### **Linux Issues**

#### **Missing dependencies**

**Solution**:
```bash
# Install required packages
# Ubuntu/Debian
sudo apt update && sudo apt install curl wget tar

# CentOS/RHEL/Fedora
sudo dnf install curl wget tar

# Alpine
sudo apk add curl wget tar
```

#### **SELinux blocking execution**

**Solution**:
```bash
# Check SELinux status
getenforce

# Temporarily disable (if appropriate)
sudo setenforce 0

# Or set proper context
sudo restorecon -v /usr/local/bin/govman
```

## ⚡ Performance Issues

### **Slow download speeds**

#### **Diagnosis**
```bash
# Test download speed
time govman install 1.21.1

# Check connection to Go servers
curl -w "@curl-format.txt" -o /dev/null -s https://go.dev/dl/
```

#### **Solutions**

**Option 1: Use mirror**
```bash
# Set mirror in config
echo "mirror:
  enabled: true
  url: https://golang.google.cn/dl/" >> ~/.govman/config.yaml
```

**Option 2: Adjust download settings**
```bash
# Reduce parallel connections
echo "download:
  parallel: true
  max_connections: 2
  timeout: 600s" >> ~/.govman/config.yaml
```

### **High memory usage during installation**

#### **Solutions**
```bash
# Disable parallel downloads
echo "download:
  parallel: false" >> ~/.govman/config.yaml

# Install versions one at a time
govman install 1.21.1
govman install 1.20.5
```

## 🌐 Network and Download Problems

### **Certificate verification errors**

#### **Solutions**

**Option 1: Update certificates**
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install ca-certificates

# CentOS/RHEL
sudo yum update ca-certificates

# macOS
brew install ca-certificates
```

**Option 2: Bypass SSL (not recommended for production)**
```bash
# Temporary workaround
export GOVMAN_INSECURE_SKIP_VERIFY=true
govman install 1.21.1
```

### **Proxy configuration**

#### **Solutions**
```bash
# Set proxy environment variables
export http_proxy=http://proxy.company.com:8080
export https_proxy=http://proxy.company.com:8080
export no_proxy=localhost,127.0.0.1

# Or configure in GOVMAN config
echo "download:
  proxy: http://proxy.company.com:8080" >> ~/.govman/config.yaml
```

### **Download corruption or checksum failures**

#### **Solutions**
```bash
# Clear cache and retry
govman clean
govman install 1.21.1

# Verify download manually
govman --verbose install 1.21.1
```

## ⚙️ Configuration Issues

### **Configuration file not loading**

#### **Diagnosis**
```bash
# Check config file location and syntax
ls -la ~/.govman/config.yaml
govman --config ~/.govman/config.yaml list
```

#### **Solutions**

**Option 1: Recreate configuration**
```bash
# Remove and recreate
rm ~/.govman/config.yaml
govman list  # This will create default config
```

**Option 2: Validate YAML syntax**
```bash
# Check YAML syntax (if you have python/yq)
python -c "import yaml; yaml.safe_load(open('~/.govman/config.yaml'))"
```

### **Custom installation directory issues**

#### **Solutions**
```bash
# Set GOVMAN_ROOT environment variable
export GOVMAN_ROOT="/custom/path"
echo 'export GOVMAN_ROOT="/custom/path"' >> ~/.bashrc

# Ensure directory exists and has correct permissions
mkdir -p "$GOVMAN_ROOT"
chmod 755 "$GOVMAN_ROOT"
```

## 📞 Getting Help

### **Before Reporting Issues**

1. **Update to latest version**:
```bash
govman selfupdate
govman init --force
```

2. **Collect diagnostic information**:
```bash
# System info
uname -a
echo $SHELL
govman --version

# GOVMAN state
govman current
govman list
ls -la ~/.govman/

# Recent error logs
govman --verbose list 2>&1
```

3. **Try safe mode**:
```bash
# Minimal configuration test
govman --config /dev/null list
```

### **Where to Get Help**

1. **GitHub Issues**: [Report bugs and request features](https://github.com/sijunda/govman/issues)
   - Include diagnostic information
   - Describe steps to reproduce
   - Mention your OS and shell

2. **GitHub Discussions**: [Ask questions and share tips](https://github.com/sijunda/govman/discussions)
   - General usage questions
   - Best practices
   - Community support

3. **Documentation**:
   - [Command Reference](commands.md)
   - [Installation Guide](installation.md)
   - [FAQ](faq.md)

### **Creating Effective Bug Reports**

Include this information:

```bash
# System Information
OS: $(uname -s) $(uname -r)
Architecture: $(uname -m)
Shell: $SHELL
GOVMAN Version: $(govman --version)

# Problem Description
What you expected to happen:
What actually happened:
Steps to reproduce:

# Diagnostic Output
$(govman --verbose current 2>&1)
$(govman list 2>&1)

# Configuration
$(cat ~/.govman/config.yaml 2>/dev/null || echo "No config file")
```

## 🔧 Emergency Recovery

### **Complete Reset**

If GOVMAN is completely broken:

```bash
# 1. Backup any important data
cp -r ~/.govman ~/.govman.backup

# 2. Complete removal
rm -rf ~/.govman
# Remove shell integration manually from config files

# 3. Fresh installation
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

# 4. Reinitialize
govman init
```

### **Minimal Working Setup**

If shell integration is causing issues:

```bash
# 1. Install without shell integration
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
# Skip govman init

# 2. Manual PATH management
export PATH="$HOME/.govman/bin:$PATH"
export PATH="$HOME/.govman/versions/go1.21.1/bin:$PATH"

# 3. Use basic commands
govman install 1.21.1
go version
```

---

**Still having issues?** Create an issue on [GitHub](https://github.com/sijunda/govman/issues) with detailed diagnostic information.