# 🤝 Contributing to GOVMAN

Welcome to GOVMAN! We're excited that you want to contribute to making Go version management better for everyone.

## 📋 Table of Contents

- [Getting Started](#-getting-started)
- [Development Setup](#-development-setup)
- [Code Contribution Guidelines](#-code-contribution-guidelines)
- [Testing](#-testing)
- [Documentation](#-documentation)
- [Release Process](#-release-process)
- [Community Guidelines](#-community-guidelines)

## 🚀 Getting Started

### **What Can You Contribute?**

- 🐛 **Bug fixes** - Help us squash bugs
- ✨ **New features** - Add functionality that users need
- 📖 **Documentation** - Improve guides and examples
- 🧪 **Tests** - Increase code coverage and reliability
- 🔧 **Shell support** - Add support for new shells
- 🌐 **Platform support** - Help with Windows, ARM64, etc.
- 🎨 **UX improvements** - Make the CLI more user-friendly

### **Before You Start**

1. **Check existing issues**: Look at [GitHub Issues](https://github.com/sijunda/govman/issues) to see what needs work
2. **Read the code**: Familiarize yourself with the codebase structure
3. **Join discussions**: Participate in [GitHub Discussions](https://github.com/sijunda/govman/discussions)
4. **Start small**: Begin with documentation or small bug fixes

## 🛠️ Development Setup

### **Prerequisites**

- **Go 1.21+** (use GOVMAN to manage this! 😉)
- **Git**
- **Make** (optional, for convenience)

### **Fork and Clone**

```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/YOUR-USERNAME/govman.git
cd govman

# 3. Add upstream remote
git remote add upstream https://github.com/sijunda/govman.git

# 4. Install GOVMAN for development
govman install 1.21.1
govman use 1.21.1 --local
```

### **Build from Source**

```bash
# Build the binary
go build -o govman cmd/govman/main.go

# Or use make (if available)
make build

# Test the binary
./govman --version
```

### **Development Dependencies**

```bash
# Install development tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# Install testing tools
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install github.com/onsi/gomega@latest
```

### **Project Structure Overview**

```
govman/
├── 📂 cmd/govman/           # Main entry point
├── 📂 internal/             # Internal packages
│   ├── 📂 cli/             # CLI commands
│   ├── 📂 config/          # Configuration management
│   ├── 📂 downloader/      # Download functionality
│   ├── 📂 golang/          # Go releases API
│   ├── 📂 logger/          # Logging utilities
│   ├── 📂 manager/         # Core version management
│   ├── 📂 progress/        # Progress indicators
│   ├── 📂 shell/           # Shell integration
│   ├── 📂 symlink/         # Symlink management
│   ├── 📂 util/            # Utilities
│   └── 📂 version/         # Version information
├── 📂 scripts/             # Installation scripts
├── 📂 docs/                # Documentation
└── 📂 tests/               # Test files
```

## 📝 Code Contribution Guidelines

### **Coding Standards**

#### **Go Code Style**
```bash
# Format code
gofmt -w .
goimports -w .

# Lint code
golangci-lint run
staticcheck ./...

# Vet code
go vet ./...
```

#### **Code Quality Requirements**
- ✅ **Formatted**: Use `gofmt` and `goimports`
- ✅ **Linted**: Pass `golangci-lint` without warnings
- ✅ **Tested**: Include unit tests for new functionality
- ✅ **Documented**: Add Go doc comments for public functions
- ✅ **Error handling**: Proper error handling and logging

#### **Commit Message Format**
```
type(scope): brief description

Longer description if needed.

Fixes #123
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `style`: Code style changes
- `chore`: Maintenance tasks

**Examples:**
```bash
feat(shell): add Fish shell support for auto-switching

Add complete Fish shell integration with auto-switching
functionality, matching the feature set of Bash and Zsh.

- Add Fish-specific shell integration
- Update init command to detect Fish
- Add Fish syntax for wrapper functions
- Update documentation

Fixes #45

fix(manager): resolve local version inconsistency

Fix critical bug where 'govman current' and 'go version'
showed different versions after using --local flag.

The issue was that the Use() function didn't update the
current session's PATH when using local versions.

Fixes #123
```

### **Pull Request Process**

#### **1. Create Feature Branch**
```bash
# Create and switch to feature branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/issue-description
```

#### **2. Make Changes**
```bash
# Make your changes
# Follow coding standards
# Add tests
# Update documentation

# Test your changes
go test ./...
./govman --version
```

#### **3. Commit Changes**
```bash
# Stage changes
git add .

# Commit with good message
git commit -m "feat(cli): add new refresh command

Add manual refresh command for edge cases where
auto-switching doesn't trigger properly.

- Add refresh command implementation
- Update command registration
- Add tests and documentation

Fixes #67"
```

#### **4. Push and Create PR**
```bash
# Push to your fork
git push origin feature/your-feature-name

# Create pull request on GitHub
# Include:
# - Clear description of changes
# - Reference to related issues
# - Screenshots if UI changes
# - Testing instructions
```

#### **5. PR Review Process**
- ✅ **Automated checks** must pass (CI/CD)
- ✅ **Code review** by maintainers
- ✅ **Testing** in different environments
- ✅ **Documentation** review if needed

### **Code Review Checklist**

**For Authors:**
- [ ] Code follows Go conventions
- [ ] All tests pass
- [ ] Documentation updated
- [ ] No breaking changes (or clearly documented)
- [ ] Performance impact considered
- [ ] Cross-platform compatibility verified

**For Reviewers:**
- [ ] Code is readable and maintainable
- [ ] Logic is correct and efficient
- [ ] Error handling is appropriate
- [ ] Tests cover new functionality
- [ ] Documentation is accurate
- [ ] No security vulnerabilities

## 🧪 Testing

### **Running Tests**

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/manager

# Run specific test
go test -run TestManagerUse ./internal/manager
```

### **Test Structure**

#### **Unit Tests**
```go
// internal/manager/manager_test.go
package manager

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestUseLocal(t *testing.T) {
    // Setup
    m := NewManager(&Config{})

    // Test
    err := m.Use("1.21.1", true, false)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "1.21.1", m.GetCurrentVersion())
}
```

#### **Integration Tests**
```go
// tests/integration/install_test.go
package integration

func TestInstallAndUse(t *testing.T) {
    // Test full workflow
    // 1. Install version
    // 2. Use version
    // 3. Verify consistency
}
```

#### **Shell Integration Tests**
```bash
# tests/shell/bash_test.sh
#!/bin/bash
# Test shell integration

test_bash_integration() {
    # Setup test environment
    export GOVMAN_ROOT="$(mktemp -d)"

    # Test auto-switching
    mkdir test-project
    echo "1.21.1" > test-project/.govman-version
    cd test-project

    # Verify version switched
    current_version=$(govman current --quiet)
    if [[ "$current_version" != *"1.21.1"* ]]; then
        echo "FAIL: Auto-switching not working"
        exit 1
    fi

    echo "PASS: Bash integration working"
}
```

### **Test Coverage Requirements**

- **Minimum coverage**: 80%
- **Critical paths**: 95%+ (installation, version switching)
- **New features**: Must include tests
- **Bug fixes**: Must include regression tests

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 📖 Documentation

### **Documentation Types**

#### **Code Documentation**
```go
// Manager handles Go version management operations.
// It provides functionality to install, uninstall, and switch
// between different Go versions.
type Manager struct {
    config *Config
    logger *Logger
}

// Use switches to the specified Go version.
// If local is true, creates a .govman-version file.
// If makeDefault is true, sets as system default.
func (m *Manager) Use(version string, local bool, makeDefault bool) error {
    // Implementation...
}
```

#### **User Documentation**
- Update relevant `.md` files in `docs/`
- Include examples and use cases
- Add troubleshooting information
- Update command reference if needed

#### **API Documentation**
```bash
# Generate Go docs
godoc -http=:6060
# Visit http://localhost:6060/pkg/github.com/sijunda/govman/
```

### **Documentation Standards**

- ✅ **Clear and concise** writing
- ✅ **Examples included** for all features
- ✅ **Cross-platform** considerations
- ✅ **Up-to-date** information
- ✅ **Searchable** content

## 🚀 Release Process

### **Version Strategy**

GOVMAN follows [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH` (e.g., `1.2.3`)
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

### **Release Checklist**

#### **Preparation**
- [ ] All tests pass
- [ ] Documentation updated
- [ ] Changelog updated
- [ ] Version bumped in code
- [ ] Cross-platform testing completed

#### **Release Creation**
```bash
# 1. Create release branch
git checkout -b release/v1.2.0

# 2. Update version
echo "v1.2.0" > VERSION

# 3. Update changelog
cat >> CHANGELOG.md << EOF
## [v1.2.0] - $(date +%Y-%m-%d)

### Added
- New refresh command
- Enhanced shell integration

### Fixed
- Local version consistency issue

### Changed
- Improved error messages
EOF

# 4. Commit and tag
git add VERSION CHANGELOG.md
git commit -m "Release v1.2.0"
git tag -a v1.2.0 -m "Release v1.2.0"

# 5. Push
git push origin release/v1.2.0
git push origin v1.2.0
```

#### **Post-Release**
- [ ] GitHub release created
- [ ] Binaries uploaded
- [ ] Installation scripts updated
- [ ] Community notified
- [ ] Documentation site updated

### **Hotfix Process**

```bash
# For critical bugs in production
git checkout -b hotfix/v1.2.1 v1.2.0
# Fix the issue
# Test thoroughly
git commit -m "fix: critical security issue"
git tag -a v1.2.1 -m "Hotfix v1.2.1"
git push origin hotfix/v1.2.1
git push origin v1.2.1
```

## 👥 Community Guidelines

### **Code of Conduct**

We are committed to providing a welcoming and inclusive environment:

- ✅ **Be respectful** and considerate
- ✅ **Be collaborative** and helpful
- ✅ **Be patient** with newcomers
- ✅ **Be constructive** in feedback
- ❌ **No harassment** or discrimination
- ❌ **No offensive** language or behavior

### **Getting Help**

#### **For Contributors**
- **Documentation**: Start with this guide and `docs/`
- **Code questions**: Ask in [GitHub Discussions](https://github.com/sijunda/govman/discussions)
- **Bug reports**: Create [GitHub Issues](https://github.com/sijunda/govman/issues)
- **Feature ideas**: Discuss in [GitHub Discussions](https://github.com/sijunda/govman/discussions)

#### **Communication Channels**
- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and community chat
- **Pull Requests**: Code review and collaboration

### **Recognition**

Contributors will be recognized:
- 🏆 **Contributors list** in README
- 🎉 **Release notes** mention significant contributions
- ⭐ **GitHub achievements** and badges
- 📣 **Social media** shout-outs for major contributions

## 🎯 Specific Contribution Areas

### **Shell Support**

Adding support for new shells:

```bash
# 1. Add shell detection in shell.go
func detectShell() string {
    shell := os.Getenv("SHELL")
    switch {
    case strings.Contains(shell, "fish"):
        return "fish"
    case strings.Contains(shell, "zsh"):
        return "zsh"
    case strings.Contains(shell, "bash"):
        return "bash"
    case strings.Contains(shell, "newshell"): // Add new shell
        return "newshell"
    default:
        return "unknown"
    }
}

# 2. Add shell-specific integration
func generateNewShellIntegration() string {
    return `
# GOVMAN integration for NewShell
function govman() {
    // NewShell-specific syntax
}
# END GOVMAN
`
}

# 3. Add tests
func TestNewShellIntegration(t *testing.T) {
    // Test NewShell integration
}

# 4. Update documentation
```

### **Platform Support**

Adding support for new platforms:

```go
// internal/downloader/platform.go
func getPlatform() (string, string) {
    goos := runtime.GOOS
    goarch := runtime.GOARCH

    switch goos {
    case "linux":
        return "linux", goarch
    case "darwin":
        return "darwin", goarch
    case "windows":
        return "windows", goarch
    case "newos": // Add new OS
        return "newos", goarch
    default:
        return "", ""
    }
}
```

### **Performance Improvements**

```go
// Example: Optimize download speed
func (d *Downloader) downloadParallel(url string, dest string) error {
    // Implement parallel chunk downloading
    // Add progress tracking
    // Implement resume capability
}
```

## 📞 Questions?

- 💬 **General questions**: [GitHub Discussions](https://github.com/sijunda/govman/discussions)
- 🐛 **Bug reports**: [GitHub Issues](https://github.com/sijunda/govman/issues)
- 💡 **Feature requests**: [GitHub Discussions - Ideas](https://github.com/sijunda/govman/discussions/categories/ideas)
- 📖 **Documentation**: Check `docs/` folder

**Happy contributing! 🎉**

---

*Thank you for making GOVMAN better for the entire Go community! Every contribution, no matter how small, makes a difference.*