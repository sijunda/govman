# ❓ Frequently Asked Questions

Complete answers to the most common questions about GOVMAN usage, features, and troubleshooting.

## 📋 Table of Contents

- [General Questions](#-general-questions)
- [Installation Questions](#-installation-questions)
- [Usage Questions](#-usage-questions)
- [Shell Integration Questions](#-shell-integration-questions)
- [Version Management Questions](#-version-management-questions)
- [Project Management Questions](#-project-management-questions)
- [Troubleshooting Questions](#-troubleshooting-questions)
- [Platform-Specific Questions](#-platform-specific-questions)

## 🌟 General Questions

### **Q: What is GOVMAN?**

**A:** GOVMAN is a fast, secure, and cross-platform Go version manager that allows you to install, manage, and switch between multiple Go versions effortlessly. It supports session-only, system-wide, and project-specific version management.

### **Q: How is GOVMAN different from other Go version managers?**

**A:** GOVMAN offers several unique advantages:
- ✅ **Production-safe shell integration** - no dangerous hooks or performance overhead
- ✅ **Immediate consistency** - `go version` and `govman current` always match
- ✅ **Cross-platform support** - works on Linux, macOS, Windows with all major shells
- ✅ **Zero configuration** - works immediately after installation
- ✅ **Parallel downloads** with resume capability
- ✅ **Built-in security** - automatic checksum verification

### **Q: Is GOVMAN free and open source?**

**A:** Yes! GOVMAN is completely free and open source under the MIT license. You can find the source code at [GitHub](https://github.com/sijunda/govman).

### **Q: Which platforms does GOVMAN support?**

**A:** GOVMAN supports:
- **Operating Systems**: Linux, macOS, Windows
- **Architectures**: x86_64 (Intel/AMD), ARM64 (Apple Silicon), x86 (32-bit)
- **Shells**: Bash, Zsh, Fish, PowerShell, Command Prompt

## 📦 Installation Questions

### **Q: How do I install GOVMAN?**

**A:** Use the one-line installation command:

**Unix/macOS/Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

**Windows PowerShell:**
```powershell
irm https://raw.githubusercontent.com/sijunda/govman/main/install.ps1 | iex
```

### **Q: Do I need administrator privileges to install GOVMAN?**

**A:** No, GOVMAN installs to your home directory (`~/.govman`) by default and doesn't require administrator privileges. If you want to install system-wide, you can use `sudo` (not recommended).

### **Q: Can I install GOVMAN offline?**

**A:** Yes, you can download the binary from the [GitHub releases page](https://github.com/sijunda/govman/releases), extract it, and place it in your PATH. Then run `govman init` to set up shell integration.

### **Q: How do I update GOVMAN?**

**A:** Use the built-in self-update command:
```bash
govman selfupdate
```

Or re-run the installation script:
```bash
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

### **Q: How do I uninstall GOVMAN completely?**

**A:** Use the uninstall script:
```bash
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/uninstall.sh | bash
```

Or manually:
```bash
# Remove binary
rm /usr/local/bin/govman  # or wherever installed

# Remove data directory
rm -rf ~/.govman

# Remove shell integration (edit your shell config files manually)
```

## 🔧 Usage Questions

### **Q: How do I install a specific Go version?**

**A:**
```bash
# Install latest stable version
govman install latest

# Install specific version
govman install 1.21.1

# Install multiple versions
govman install 1.19.5 1.20.5 1.21.1
```

### **Q: How do I switch between Go versions?**

**A:**
```bash
# Session-only (temporary)
govman use 1.21.1

# Set as system default (permanent)
govman use 1.21.1 --default

# Set for current project
govman use 1.21.1 --local
```

### **Q: What's the difference between session-only, default, and local versions?**

**A:**
- **Session-only**: Active only for the current terminal session
- **System default**: Becomes the default for all new terminal sessions
- **Project local**: Creates `.govman-version` file, auto-switches when entering directory

### **Q: How do I check which Go version is currently active?**

**A:**
```bash
# GOVMAN's view
govman current

# Go's view (should match govman current)
go version

# List all installed versions
govman list
```

### **Q: How do I see which Go versions are available for download?**

**A:**
```bash
govman list --remote
```

## 🐚 Shell Integration Questions

### **Q: What is shell integration and do I need it?**

**A:** Shell integration enables automatic version switching when you enter directories with `.govman-version` files. It's optional but highly recommended for project-based development.

To enable:
```bash
govman init
source ~/.bashrc  # or your shell's config file
```

### **Q: Which shells are supported?**

**A:**
- **Full support**: Bash, Zsh, Fish, PowerShell
- **Basic support**: Command Prompt (Windows)

### **Q: Can I use GOVMAN without shell integration?**

**A:** Yes! You can use all GOVMAN commands without shell integration. You just won't get automatic version switching based on `.govman-version` files.

### **Q: Will shell integration conflict with my existing shell configuration?**

**A:** No, GOVMAN uses production-safe integration that won't interfere with your existing configuration. It doesn't override global variables or use dangerous shell hooks.

### **Q: GOVMAN is not switching versions automatically. What's wrong?**

**A:** Check the following:
1. Ensure shell integration is enabled: `govman init`
2. Restart your terminal or run: `source ~/.bashrc`
3. Check if `.govman-version` file exists in your project directory
4. Verify your shell is supported

## 🔄 Version Management Questions

### **Q: Can I have multiple Go versions installed simultaneously?**

**A:** Yes! That's the whole point of GOVMAN. You can install as many versions as you want:
```bash
govman install 1.19.5 1.20.5 1.21.1
govman list  # See all installed versions
```

### **Q: How much disk space do Go versions take?**

**A:** Each Go version typically takes 200-500MB. You can check disk usage with:
```bash
govman list  # Shows disk usage for each version
```

### **Q: How do I remove unused Go versions?**

**A:**
```bash
# Remove specific version
govman uninstall 1.19.5

# Clean download cache
govman clean
```

### **Q: Can I use pre-release or beta versions?**

**A:** Yes! GOVMAN supports pre-release versions:
```bash
govman install 1.22rc1
govman install 1.22beta1
```

### **Q: What happens if I try to use a version that's not installed?**

**A:** GOVMAN will show an error and list your installed versions. You'll need to install the version first:
```bash
govman install 1.21.1
govman use 1.21.1
```

## 📁 Project Management Questions

### **Q: How do I set a specific Go version for my project?**

**A:** Use the `--local` flag:
```bash
cd my-project
govman use 1.21.1 --local
```

This creates a `.govman-version` file with the version number.

### **Q: What is a `.govman-version` file?**

**A:** It's a simple text file containing a Go version number. When you enter a directory with this file, GOVMAN automatically switches to that version (if shell integration is enabled).

Example `.govman-version` file:
```
1.21.1
```

### **Q: Should I commit `.govman-version` files to version control?**

**A:** Yes! This ensures all team members use the same Go version:
```bash
git add .govman-version
git commit -m "Set Go version for project"
```

### **Q: Can I have different Go versions for different projects?**

**A:** Absolutely! Each project can have its own `.govman-version` file:
```bash
# Project A uses Go 1.21.1
cd project-a
govman use 1.21.1 --local

# Project B uses Go 1.20.5
cd ../project-b
govman use 1.20.5 --local
```

### **Q: What happens when I leave a project directory?**

**A:** GOVMAN reverts to your system default version when you leave a directory with a `.govman-version` file.

## 🔧 Troubleshooting Questions

### **Q: "Command not found: govman" - What's wrong?**

**A:** GOVMAN is not in your PATH. Try:
1. Restart your terminal
2. Check if the binary exists: `ls -la ~/.govman/bin/govman`
3. Add to PATH manually: `export PATH="$HOME/.govman/bin:$PATH"`
4. Re-run installation script

### **Q: `govman current` and `go version` show different versions. Why?**

**A:** This was a major bug that has been **FIXED** in the latest version! Update GOVMAN:
```bash
govman selfupdate
# Or reinstall
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

### **Q: GOVMAN says a version is installed but Go commands don't work?**

**A:** The version might be corrupted. Try:
1. Check current version: `govman current`
2. Reinstall the version: `govman uninstall 1.21.1 && govman install 1.21.1`
3. Use refresh command: `govman refresh`

### **Q: Downloads are very slow. Can I speed them up?**

**A:** Yes, you can:
1. Use a mirror (for China users):
```bash
echo "mirror:
  enabled: true
  url: https://golang.google.cn/dl/" >> ~/.govman/config.yaml
```

2. Check your internet connection
3. Try the installation again

### **Q: I'm getting permission errors during installation?**

**A:** Try:
1. Install to user directory (default behavior)
2. Check directory permissions: `ls -la ~/.govman/`
3. Create directory manually: `mkdir -p ~/.govman/bin`

### **Q: Shell integration stopped working after updating my shell configuration?**

**A:** Reinitialize shell integration:
```bash
govman init --force
source ~/.bashrc  # or your shell's config file
```

## 🖥️ Platform-Specific Questions

### **Q: Does GOVMAN work on Windows?**

**A:** Yes! GOVMAN has full Windows support:
- **PowerShell**: Full featured support
- **Command Prompt**: Basic support with wrapper commands

### **Q: Does GOVMAN work on Apple Silicon Macs?**

**A:** Yes! GOVMAN natively supports ARM64 architecture for Apple Silicon Macs.

### **Q: Does GOVMAN work in WSL (Windows Subsystem for Linux)?**

**A:** Yes! Use the Linux installation method in WSL:
```bash
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

### **Q: Can I use GOVMAN in Docker containers?**

**A:** Yes! Add this to your Dockerfile:
```dockerfile
RUN curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
RUN govman install 1.21.1 && govman use 1.21.1 --default
```

### **Q: Does GOVMAN work with GitHub Actions/CI?**

**A:** Yes! Example GitHub Actions setup:
```yaml
- name: Install GOVMAN
  run: curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

- name: Setup Go
  run: |
    govman install 1.21.1
    govman use 1.21.1
```

### **Q: I'm behind a corporate firewall. Can I still use GOVMAN?**

**A:** Yes, configure proxy settings:
```bash
export http_proxy=http://proxy.company.com:8080
export https_proxy=http://proxy.company.com:8080
```

Or download manually from the [releases page](https://github.com/sijunda/govman/releases).

## 🚀 Advanced Questions

### **Q: Can I customize the installation directory?**

**A:** Yes, set the `GOVMAN_ROOT` environment variable:
```bash
export GOVMAN_ROOT="/custom/path"
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
```

### **Q: Can I configure GOVMAN behavior?**

**A:** Yes, edit `~/.govman/config.yaml`:
```yaml
install_dir: ~/.govman/versions
cache_dir: ~/.govman/cache
download:
  parallel: true
  max_connections: 4
auto_switch:
  enabled: true
```

### **Q: How do I use GOVMAN in scripts?**

**A:** GOVMAN is script-friendly:
```bash
#!/bin/bash
# Ensure specific Go version
if [ -f .govman-version ]; then
    govman use $(cat .govman-version)
fi

# Use Go
go build ./...
```

### **Q: Can I use GOVMAN with other development tools?**

**A:** Yes! GOVMAN works with:
- **IDEs**: VSCode, GoLand, vim-go
- **Build tools**: Make, Bazel, Mage
- **CI/CD**: GitHub Actions, GitLab CI, Jenkins
- **Containers**: Docker, Podman

### **Q: How does GOVMAN compare to gvm or g?**

**A:** GOVMAN offers:
- ✅ **Better cross-platform support** (Windows, ARM64)
- ✅ **Production-safe shell integration**
- ✅ **Immediate version consistency**
- ✅ **Built-in security and checksum verification**
- ✅ **Zero configuration setup**
- ✅ **Modern CLI with rich output**

## 🤝 Community Questions

### **Q: How do I report bugs or request features?**

**A:** Use GitHub:
- **Bug reports**: [GitHub Issues](https://github.com/sijunda/govman/issues)
- **Feature requests**: [GitHub Discussions](https://github.com/sijunda/govman/discussions)

### **Q: How do I contribute to GOVMAN?**

**A:** We welcome contributions! See our [Contributing Guide](contributing.md) for details.

### **Q: Where can I get help?**

**A:**
1. **Documentation**: [Complete docs](https://github.com/sijunda/govman/tree/main/docs)
2. **GitHub Discussions**: [Community Q&A](https://github.com/sijunda/govman/discussions)
3. **GitHub Issues**: [Bug reports](https://github.com/sijunda/govman/issues)

### **Q: Is there a roadmap for GOVMAN?**

**A:** Yes! Check our [GitHub Issues](https://github.com/sijunda/govman/issues) for planned features and [Discussions](https://github.com/sijunda/govman/discussions) for community input.

---

**Still have questions?** Ask in our [GitHub Discussions](https://github.com/sijunda/govman/discussions) or check our [Troubleshooting Guide](troubleshooting.md).