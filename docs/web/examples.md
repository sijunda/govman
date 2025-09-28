# 🌟 Examples & Best Practices

Real-world examples and proven patterns for using GOVMAN effectively in development teams and production environments.

## 📋 Table of Contents

- [Quick Start Examples](#-quick-start-examples)
- [Project Workflows](#-project-workflows)
- [Team Collaboration](#-team-collaboration)
- [CI/CD Integration](#-cicd-integration)
- [Development Environments](#-development-environments)
- [Performance Optimization](#-performance-optimization)
- [Security Best Practices](#-security-best-practices)
- [Advanced Use Cases](#-advanced-use-cases)

## 🚀 Quick Start Examples

### **Getting Started in 2 Minutes**

```bash
# 1. Install GOVMAN
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

# 2. Initialize shell integration
govman init
source ~/.bashrc

# 3. Install and use latest Go
govman install latest
govman use latest --default

# 4. Verify everything works
govman current
go version  # Should match govman current
```

### **Basic Version Management**

```bash
# Install multiple versions for testing
govman install 1.19.5 1.20.5 1.21.1

# Quick switching between versions
govman use 1.21.1  # Latest features
govman use 1.20.5  # Stable version
govman use 1.19.5  # Legacy compatibility

# Check what's installed
govman list
```

### **Project Setup**

```bash
# Create new project with specific Go version
mkdir my-go-project && cd my-go-project
govman use 1.21.1 --local

# Initialize Go module
go mod init my-go-project

# Verify version consistency
govman current  # Should show: Go 1.21.1
go version      # Should show: go1.21.1
```

## 📁 Project Workflows

### **Single Project Development**

```bash
# Method 1: Using govman command
cd my-project
govman use 1.21.1 --local
# Creates .govman-version file automatically

# Method 2: Manual file creation
echo "1.21.1" > .govman-version
govman refresh  # Apply the version

# Working with the project
go mod tidy
go build ./...
go test ./...
```

### **Multi-Project Development**

```bash
# Project structure
~/projects/
├── legacy-app/     # Uses Go 1.19.5
├── current-app/    # Uses Go 1.20.5
└── new-app/        # Uses Go 1.21.1

# Setup each project
cd ~/projects/legacy-app
govman use 1.19.5 --local

cd ~/projects/current-app
govman use 1.20.5 --local

cd ~/projects/new-app
govman use 1.21.1 --local

# Automatic switching when moving between projects
cd ~/projects/legacy-app  # Auto-switches to Go 1.19.5
cd ~/projects/new-app     # Auto-switches to Go 1.21.1
```

### **Microservices Architecture**

```bash
# Different services can use different Go versions
~/microservices/
├── auth-service/      # Go 1.21.1 (latest features)
├── payment-service/   # Go 1.20.5 (stable)
├── legacy-service/    # Go 1.19.5 (compatibility)
└── docker-compose.yml

# Each service has its own .govman-version
cd auth-service && echo "1.21.1" > .govman-version
cd payment-service && echo "1.20.5" > .govman-version
cd legacy-service && echo "1.19.5" > .govman-version
```

## 👥 Team Collaboration

### **Onboarding New Team Members**

**Team Lead Setup:**
```bash
# 1. Set project Go version
cd company-project
govman use 1.21.1 --local
git add .govman-version
git commit -m "Set team Go version to 1.21.1"
git push
```

**New Developer Setup:**
```bash
# 1. Install GOVMAN
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

# 2. Clone project
git clone https://github.com/company/project.git
cd project

# 3. GOVMAN automatically detects and installs required version
govman install $(cat .govman-version)  # Install Go 1.21.1
# Auto-switches when entering directory
```

### **Team Standards Enforcement**

**`.govman-version` in Repository Root:**
```bash
# Create team standard
echo "1.21.1" > .govman-version
echo "# Go Version" >> README.md
echo "This project uses Go $(cat .govman-version)" >> README.md
echo "GOVMAN will automatically switch to this version" >> README.md

# Add to version control
git add .govman-version README.md
git commit -m "Standardize Go version across team"
```

**Pre-commit Hook for Version Consistency:**
```bash
#!/bin/bash
# .git/hooks/pre-commit
if [ -f .govman-version ]; then
    required_version=$(cat .govman-version)
    current_version=$(go version | cut -d' ' -f3 | sed 's/go//')

    if [ "$current_version" != "$required_version" ]; then
        echo "Error: Wrong Go version!"
        echo "Required: $required_version"
        echo "Current:  $current_version"
        echo "Run: govman use $required_version"
        exit 1
    fi
fi
```

### **Code Review Workflow**

```bash
# Reviewer checks out PR branch
git checkout feature-branch

# GOVMAN automatically switches to branch's Go version
# (if .govman-version was changed in the PR)

# Verify the build works with specified version
go mod tidy
go build ./...
go test ./...

# Test with other versions if needed
govman use 1.20.5  # Test backward compatibility
go test ./...
```

## 🔄 CI/CD Integration

### **GitHub Actions Workflow**

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.19.5, 1.20.5, 1.21.1]

    steps:
    - uses: actions/checkout@v3

    - name: Install GOVMAN
      run: curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

    - name: Setup Go version
      run: |
        govman install ${{ matrix.go-version }}
        govman use ${{ matrix.go-version }}

    - name: Verify Go version
      run: |
        govman current
        go version

    - name: Run tests
      run: |
        go mod download
        go test ./...
```

### **Docker Integration**

```dockerfile
# Dockerfile
FROM ubuntu:22.04

# Install GOVMAN
RUN apt-get update && apt-get install -y curl
RUN curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

# Install specific Go version
COPY .govman-version .
RUN govman install $(cat .govman-version)
RUN govman use $(cat .govman-version) --default

# Copy application
COPY . /app
WORKDIR /app

# Build application
RUN go mod download
RUN go build -o main .

CMD ["./main"]
```

### **Jenkins Pipeline**

```groovy
// Jenkinsfile
pipeline {
    agent any

    stages {
        stage('Setup') {
            steps {
                sh 'curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash'

                script {
                    if (fileExists('.govman-version')) {
                        def goVersion = readFile('.govman-version').trim()
                        sh "govman install ${goVersion}"
                        sh "govman use ${goVersion}"
                    } else {
                        sh 'govman install latest'
                        sh 'govman use latest'
                    }
                }
            }
        }

        stage('Test') {
            steps {
                sh 'go mod download'
                sh 'go test ./...'
            }
        }

        stage('Build') {
            steps {
                sh 'go build ./...'
            }
        }
    }
}
```

## 🖥️ Development Environments

### **VSCode Integration**

**`.vscode/settings.json`:**
```json
{
  "go.goroot": "${env:HOME}/.govman/versions/go1.21.1",
  "go.gopath": "${env:HOME}/go",
  "go.toolsGopath": "${env:HOME}/go",
  "go.useLanguageServer": true,
  "go.alternateTools": {
    "go": "${env:HOME}/.govman/versions/go1.21.1/bin/go"
  }
}
```

**Dynamic Go Root Script:**
```bash
#!/bin/bash
# .vscode/update-go-settings.sh
GOVMAN_VERSION=$(govman current --quiet | grep -o 'Go [0-9.]*' | cut -d' ' -f2)
GO_ROOT="$HOME/.govman/versions/go$GOVMAN_VERSION"

cat > .vscode/settings.json << EOF
{
  "go.goroot": "$GO_ROOT",
  "go.alternateTools": {
    "go": "$GO_ROOT/bin/go"
  }
}
EOF

echo "Updated VSCode settings for Go $GOVMAN_VERSION"
```

### **GoLand/IntelliJ Integration**

```bash
# Script to update GoLand Go SDK
#!/bin/bash
# update-goland-sdk.sh
GOVMAN_VERSION=$(govman current --quiet | grep -o 'Go [0-9.]*' | cut -d' ' -f2)
GO_ROOT="$HOME/.govman/versions/go$GOVMAN_VERSION"

echo "Current Go Root: $GO_ROOT"
echo "Update GoLand SDK to point to: $GO_ROOT"
echo "GoLand > Preferences > Go > GOROOT > $GO_ROOT"
```

### **Terminal Setup**

**Bash/Zsh Prompt with Go Version:**
```bash
# Add to ~/.bashrc or ~/.zshrc
show_go_version() {
    if command -v govman &> /dev/null; then
        local version=$(govman current --quiet 2>/dev/null | grep -o 'Go [0-9.]*' | cut -d' ' -f2)
        if [ -n "$version" ]; then
            echo " 🐹$version"
        fi
    fi
}

# Update PS1 to include Go version
PS1="${PS1}\$(show_go_version) $ "
```

## ⚡ Performance Optimization

### **Parallel Development**

```bash
# Install multiple versions in parallel
govman install 1.19.5 1.20.5 1.21.1 &

# Multiple terminal sessions for different projects
# Terminal 1: Legacy project (Go 1.19.5)
cd ~/projects/legacy && govman use 1.19.5 --local

# Terminal 2: Current project (Go 1.20.5)
cd ~/projects/current && govman use 1.20.5 --local

# Terminal 3: Experimental project (Go 1.21.1)
cd ~/projects/experimental && govman use 1.21.1 --local
```

### **Cache Optimization**

```bash
# Pre-download versions for team
govman install 1.19.5 1.20.5 1.21.1

# Clean old cache periodically
govman clean

# Check disk usage
govman list  # Shows disk usage per version
```

### **Build Optimization**

```bash
# Use version-specific build caches
export GOCACHE="$HOME/.cache/go-build-$(govman current --quiet | grep -o '[0-9.]*')"

# Version-specific module cache
export GOMODCACHE="$HOME/.cache/go-mod-$(govman current --quiet | grep -o '[0-9.]*')"
```

## 🔒 Security Best Practices

### **Secure Installation**

```bash
# Verify installation script before running
curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh -o install.sh
# Review the script
cat install.sh
# Run after verification
bash install.sh
```

### **Checksum Verification**

```bash
# GOVMAN automatically verifies checksums
# For manual verification:
govman --verbose install 1.21.1  # Shows checksum verification

# Check installed version integrity
cd ~/.govman/versions/go1.21.1
find . -type f -name "*.so" -o -name "go" | xargs shasum -a 256
```

### **Team Security**

```bash
# Lock versions in repository
echo "1.21.1" > .govman-version
git add .govman-version

# Document security requirements
echo "# Security" >> README.md
echo "This project requires Go $(cat .govman-version) for security compliance" >> README.md
```

## 🎯 Advanced Use Cases

### **Cross-Compilation Setup**

```bash
# Install Go versions for different targets
govman install 1.21.1
govman use 1.21.1

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o app-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o app-windows-amd64.exe .
GOOS=darwin GOARCH=arm64 go build -o app-darwin-arm64 .
```

### **Performance Testing Across Versions**

```bash
#!/bin/bash
# benchmark-versions.sh
versions=("1.19.5" "1.20.5" "1.21.1")

for version in "${versions[@]}"; do
    echo "Testing Go $version..."
    govman use $version

    echo "Building..."
    go build -o "bench-$version" .

    echo "Benchmarking..."
    go test -bench=. -benchmem > "results-$version.txt"

    echo "Performance for Go $version:"
    grep "Benchmark" "results-$version.txt"
    echo "---"
done
```

### **Version Migration Strategy**

```bash
#!/bin/bash
# migrate-go-version.sh
OLD_VERSION="1.20.5"
NEW_VERSION="1.21.1"

echo "Migrating from Go $OLD_VERSION to $NEW_VERSION"

# 1. Install new version
govman install $NEW_VERSION

# 2. Test with new version
govman use $NEW_VERSION
go mod tidy
go test ./...

if [ $? -eq 0 ]; then
    echo "✅ Tests pass with Go $NEW_VERSION"

    # 3. Update project version
    govman use $NEW_VERSION --local

    # 4. Commit changes
    git add .govman-version
    git commit -m "Upgrade Go version from $OLD_VERSION to $NEW_VERSION"

    echo "✅ Migration complete"
else
    echo "❌ Tests failed with Go $NEW_VERSION"
    govman use $OLD_VERSION
    echo "Reverted to Go $OLD_VERSION"
fi
```

### **Development Environment Sync**

```bash
#!/bin/bash
# sync-dev-env.sh
# Synchronize development environment across machines

# Export current setup
cat > dev-env.txt << EOF
# GOVMAN Development Environment
go_versions=$(govman list --quiet)
current_version=$(govman current --quiet)
projects_with_versions=$(find ~/projects -name ".govman-version" -exec echo {} \; -exec cat {} \;)
EOF

# On new machine, restore environment
restore_env() {
    # Install all versions from old machine
    grep "Go" dev-env.txt | while read -r version; do
        govman install "$version"
    done

    # Set default version
    default_version=$(grep "current_version" dev-env.txt | cut -d'=' -f2)
    govman use "$default_version" --default
}
```

---

**These examples cover the most common and advanced use cases for GOVMAN. Mix and match these patterns to create workflows that fit your specific needs!**