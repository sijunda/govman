# 🚀 GOVMAN Documentation

**Complete documentation for the fastest, most reliable Go version manager with zero-configuration setup and cross-platform support.**

## 📖 Table of Contents

- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [Shell Support](#-shell-support)
- [Core Features](#-core-features)
- [Production Features](#-production-features)
- [Documentation Index](#-documentation-index)
- [Recent Updates](#-recent-updates)

## ⚡ Quick Start

### One-Line Installation
```bash
# Unix/macOS/Linux
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

# Windows PowerShell
irm https://raw.githubusercontent.com/sijunda/govman/main/install.ps1 | iex
```

### Essential Commands
```bash
# Install and use latest Go
govman install latest
govman use latest --default

# Project-specific version
govman use 1.21.1 --local

# Check what's active
govman current
go version  # Should match govman current
```

## 🌍 Cross-Platform Support

GOVMAN works seamlessly across all major platforms:

### **Operating Systems**
- ✅ **Linux** (All distributions)
- ✅ **macOS** (Intel & Apple Silicon)
- ✅ **Windows** (10, 11, Server)

### **Architecture Support**
- ✅ **x86_64** (Intel/AMD 64-bit)
- ✅ **ARM64** (Apple Silicon, ARM servers)
- ✅ **x86** (32-bit legacy systems)

## 🐚 Shell Support

### **Full Production Support** ✅

| Shell | Auto-Switch | use/refresh | Error Handling | Platform |
|-------|-------------|-------------|----------------|----------|
| **Bash** | ✅ | ✅ | ✅ | Linux, macOS, WSL |
| **Zsh** | ✅ | ✅ | ✅ | macOS, Linux |
| **Fish** | ✅ | ✅ | ✅ | Linux, macOS |
| **PowerShell** | ✅ | ✅ | ✅ | Windows, Cross-platform |

### **Basic Support** ⚠️

| Shell | use/refresh | Manual Refresh | Platform |
|-------|-------------|----------------|----------|
| **Cmd** | ✅ | ✅ | Windows |

### **What Each Shell Gets**

#### **Bash & Zsh** (Full Featured)
```bash
# Auto-switch on directory change
cd my-project                # Automatically switches to project's Go version
govman use 1.21.1 --local   # Immediately updates both go version and govman current
govman refresh               # Manual refresh when needed
```

#### **Fish** (Full Featured)
```fish
# Same features as Bash/Zsh with Fish syntax
cd my-project
govman use 1.21.1 --local
govman refresh
```

#### **PowerShell** (Full Featured)
```powershell
# Native PowerShell integration
cd my-project
govman use 1.21.1 --local
govman refresh
```

#### **Command Prompt** (Basic)
```cmd
REM Basic functionality with wrapper
govman_wrapper use 1.21.1 --local
govman_wrapper refresh
```

## 🎯 Core Features

### **Zero Configuration**
- **No setup required** - works immediately after installation
- **Automatic shell detection** - finds and configures your shell
- **Smart PATH management** - handles environment automatically

### **Lightning Fast Performance**
- **Instant version switching** - changes take effect immediately
- **Parallel downloads** with resume capability
- **Intelligent caching** for offline usage
- **Production-safe integration** - no performance overhead

### **Project Management**
```bash
# Set project-specific version
govman use 1.19.5 --local
echo "1.19.5" > .govman-version  # Creates project file

# Automatic switching
cd my-old-project  # Uses Go 1.19.5
cd my-new-project  # Uses Go 1.21.1
cd ..              # Back to default version
```

### **Version Management**
```bash
# Install multiple versions
govman install 1.19.5 1.20.5 1.21.1

# List what's available
govman list                    # Installed versions
govman list --remote          # Available for download

# Detailed information
govman info 1.21.1           # Version details
govman current                # Currently active
```

## 🏭 Production Features

### **Safe Shell Integration**
- ✅ **No dangerous hooks** - no DEBUG traps or PROMPT_COMMAND override
- ✅ **No user conflicts** - preserves existing shell configuration
- ✅ **Error recovery** - graceful fallbacks on failures
- ✅ **Performance optimized** - minimal overhead

### **The Fix That Changed Everything**

**Problem Solved**: The major inconsistency issue where `govman current` and `go version` showed different versions after using `--local` flag.

**Before Fix** ❌:
```bash
govman use 1.14.1 --local
govman current    # Showed: Go 1.14.1
go version        # Showed: go1.25.1  # INCONSISTENT!
```

**After Fix** ✅:
```bash
govman use 1.14.1 --local
govman current    # Shows: Go 1.14.1
go version        # Shows: go1.14.1   # CONSISTENT!
```

### **What Was Fixed**
1. **Immediate PATH Updates** - Local versions now update current session immediately
2. **Cross-Shell Consistency** - All shells behave the same way
3. **Production Safety** - Removed risky shell hooks and global variables
4. **Error Handling** - Added comprehensive error recovery

### **Enterprise Ready**
- ✅ **Team consistency** with `.govman-version` files
- ✅ **CI/CD integration** - reliable in automated environments
- ✅ **Multiple environments** - session, project, and global scopes
- ✅ **Audit trail** - clear logging and status reporting

## 📖 Documentation Index

### **Getting Started**
- **[Installation Guide](installation.md)** - Complete setup for all platforms
- **[Quick Start](quick-start.md)** - Get up and running in 5 minutes
- **[Shell Integration](shell-integration.md)** - Platform-specific setup

### **Reference**
- **[Command Reference](commands.md)** - Complete API documentation
- **[Configuration](configuration.md)** - Advanced configuration options
- **[Shell Support](shell-support.md)** - Shell-specific features

### **Guides**
- **[Best Practices](best-practices.md)** - Production usage patterns
- **[Team Workflows](team-workflows.md)** - Multi-developer setups
- **[CI/CD Integration](cicd-integration.md)** - Automated environments

### **Help**
- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions
- **[FAQ](faq.md)** - Frequently asked questions
- **[Migration Guide](migration.md)** - Moving from other version managers

### **Examples**
- **[Real-World Examples](examples.md)** - Practical usage scenarios
- **[Project Templates](project-templates.md)** - Ready-to-use configurations
- **[Advanced Usage](advanced-usage.md)** - Power user features

### **Development**
- **[Contributing](contributing.md)** - How to contribute to GOVMAN
- **[Architecture](architecture.md)** - Technical design and implementation
- **[Testing](testing.md)** - Test suite and quality assurance

## 🆕 Recent Updates

### **v1.1.0 - Production-Safe Cross-Platform Release**

#### **🔧 Major Fixes**
- ✅ **Fixed local version inconsistency** - `go version` now matches `govman current`
- ✅ **Production-safe shell integration** - removed dangerous hooks and traps
- ✅ **Cross-platform consistency** - all shells behave identically

#### **🚀 New Features**
- ✅ **`govman refresh` command** - manual refresh for edge cases
- ✅ **Enhanced Fish shell support** - full feature parity
- ✅ **Advanced PowerShell integration** - native Windows experience
- ✅ **Command Prompt support** - basic functionality for legacy systems

#### **🛡️ Security & Stability**
- ✅ **Safe shell wrappers** - no conflicts with user configurations
- ✅ **Error recovery** - graceful handling of edge cases
- ✅ **Performance optimization** - zero overhead in daily usage

#### **🌐 Cross-Platform**
- ✅ **Universal shell support** - works on all major shells
- ✅ **ARM64 compatibility** - native Apple Silicon support
- ✅ **Windows improvements** - better PowerShell and CMD support

### **Upgrade Instructions**
```bash
# Automatic update
govman selfupdate

# Manual update
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

# Update shell integration
govman init --force
source ~/.bashrc  # or your shell's config file
```

## 🔗 Quick Links

- **GitHub Repository**: https://github.com/sijunda/govman
- **Release Notes**: https://github.com/sijunda/govman/releases
- **Issue Tracker**: https://github.com/sijunda/govman/issues
- **Discussions**: https://github.com/sijunda/govman/discussions

## 💬 Community & Support

- **🐛 Bug Reports**: [GitHub Issues](https://github.com/sijunda/govman/issues)
- **💡 Feature Requests**: [GitHub Discussions](https://github.com/sijunda/govman/discussions)
- **❓ Questions**: [GitHub Discussions Q&A](https://github.com/sijunda/govman/discussions/categories/q-a)
- **📖 Documentation**: [Complete Docs](https://govman.dev/docs)

---

**Made with ❤️ for the Go community**

*GOVMAN - Making Go version management simple, fast, and reliable across all platforms.*