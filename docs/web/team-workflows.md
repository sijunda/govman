# 👥 Team Workflows

Comprehensive guide for using GOVMAN in team environments, from small startups to large enterprises.

## 📋 Table of Contents

- [Team Setup Strategies](#-team-setup-strategies)
- [Version Standardization](#-version-standardization)
- [Development Workflows](#-development-workflows)
- [Code Review Process](#-code-review-process)
- [Branch Management](#-branch-management)
- [Release Management](#-release-management)
- [Enterprise Integration](#-enterprise-integration)
- [Troubleshooting Team Issues](#-troubleshooting-team-issues)

## 🚀 Team Setup Strategies

### **Small Team (2-5 Developers)**

**Quick Setup:**
```bash
# Team lead initializes project
cd team-project
govman use 1.21.1 --local
git add .govman-version
git commit -m "Set team Go version to 1.21.1"

# Team members clone and setup
git clone https://github.com/team/project.git
cd project
govman install $(cat .govman-version)
# Auto-switches when entering directory
```

**Benefits:**
- ✅ Instant version consistency
- ✅ No configuration meetings needed
- ✅ New developers get correct version automatically

### **Medium Team (5-20 Developers)**

**Structured Setup:**
```bash
# 1. Define team standards
mkdir -p .govman
echo "1.21.1" > .govman-version
echo "# Team Go Standards" > .govman/README.md
echo "- Production: Go $(cat .govman-version)" >> .govman/README.md
echo "- Development: Same as production" >> .govman/README.md
echo "- Testing: 1.20.5, 1.21.1" >> .govman/README.md

# 2. Create onboarding script
cat > scripts/setup-dev-env.sh << 'EOF'
#!/bin/bash
echo "Setting up development environment..."

# Install GOVMAN if not present
if ! command -v govman &> /dev/null; then
    echo "Installing GOVMAN..."
    curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash
    source ~/.bashrc
fi

# Install project Go version
if [ -f .govman-version ]; then
    PROJECT_VERSION=$(cat .govman-version)
    echo "Installing Go $PROJECT_VERSION..."
    govman install "$PROJECT_VERSION"
    echo "✅ Development environment ready!"
else
    echo "❌ No .govman-version found"
    exit 1
fi
EOF
chmod +x scripts/setup-dev-env.sh

# 3. Update README with setup instructions
cat >> README.md << 'EOF'

## Development Setup

1. Clone the repository
2. Run: `./scripts/setup-dev-env.sh`
3. Start coding!

The project will automatically use Go $(cat .govman-version).
EOF

git add .govman-version .govman/ scripts/ README.md
git commit -m "Add team development environment setup"
```

### **Large Team/Enterprise (20+ Developers)**

**Enterprise Setup:**
```bash
# 1. Create team configuration
mkdir -p .govman/team
cat > .govman/team/config.yaml << 'EOF'
team:
  name: "Backend Team"
  go_versions:
    production: "1.21.1"
    staging: "1.21.1"
    development: "1.21.1"
    testing: ["1.20.5", "1.21.1"]

  policies:
    require_version_file: true
    auto_install: true
    testing_required: true

  contacts:
    lead: "senior-dev@company.com"
    infrastructure: "devops@company.com"
EOF

# 2. Create validation script
cat > .govman/team/validate.sh << 'EOF'
#!/bin/bash
# Team environment validation

PROJECT_VERSION=$(cat .govman-version 2>/dev/null)
CURRENT_VERSION=$(govman current --quiet 2>/dev/null | grep -o '[0-9.]*')

if [ -z "$PROJECT_VERSION" ]; then
    echo "❌ Missing .govman-version file"
    echo "Run: govman use <version> --local"
    exit 1
fi

if [ "$CURRENT_VERSION" != "$PROJECT_VERSION" ]; then
    echo "❌ Version mismatch!"
    echo "Project requires: $PROJECT_VERSION"
    echo "Currently using: $CURRENT_VERSION"
    echo "Run: govman use $PROJECT_VERSION"
    exit 1
fi

echo "✅ Go version is correct: $PROJECT_VERSION"
EOF
chmod +x .govman/team/validate.sh

# 3. Git hooks for validation
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# Validate Go version before commit
if [ -f .govman/team/validate.sh ]; then
    ./.govman/team/validate.sh
fi
EOF
chmod +x .git/hooks/pre-commit
```

## 📌 Version Standardization

### **Single Version Strategy**

**When to use:**
- Small, focused team
- Single product/service
- Rapid development cycle

**Implementation:**
```bash
# Set and lock team version
govman use 1.21.1 --local

# Document the decision
cat > docs/go-version-policy.md << 'EOF'
# Go Version Policy

## Current Version: 1.21.1

### Why this version?
- Latest stable release
- Required for generics support
- Performance improvements
- Security updates

### Update Policy
- Review quarterly
- Update after 2-week testing period
- Coordinate with infrastructure team
EOF

git add .govman-version docs/go-version-policy.md
git commit -m "Standardize on Go 1.21.1"
```

### **Multi-Version Strategy**

**When to use:**
- Large organization
- Multiple products
- Legacy system maintenance

**Implementation:**
```bash
# Different services, different versions
services/
├── auth-service/         # .govman-version: 1.21.1
├── payment-service/      # .govman-version: 1.20.5
├── legacy-api/          # .govman-version: 1.19.5
└── experimental-service/ # .govman-version: 1.22rc1

# Create version matrix
cat > docs/go-version-matrix.md << 'EOF'
# Go Version Matrix

| Service | Go Version | Status | Migration Plan |
|---------|------------|--------|----------------|
| auth-service | 1.21.1 | ✅ Current | N/A |
| payment-service | 1.20.5 | ⚠️ Upgrade planned | Q2 2024 |
| legacy-api | 1.19.5 | 🚨 Legacy | Q3 2024 |
| experimental | 1.22rc1 | 🧪 Testing | TBD |
EOF
```

### **Version Migration Planning**

```bash
#!/bin/bash
# migration-planner.sh
# Plan team-wide Go version migrations

CURRENT_VERSION="1.20.5"
TARGET_VERSION="1.21.1"
MIGRATION_DATE="2024-03-01"

cat > migration-plan.md << EOF
# Go Migration: $CURRENT_VERSION → $TARGET_VERSION

## Timeline: $MIGRATION_DATE

### Phase 1: Preparation (Week 1-2)
- [ ] Install $TARGET_VERSION on all dev machines
- [ ] Run compatibility tests
- [ ] Update CI/CD pipelines
- [ ] Prepare rollback plan

### Phase 2: Testing (Week 3)
- [ ] Create feature branch with $TARGET_VERSION
- [ ] Run full test suite
- [ ] Performance benchmarks
- [ ] Security scan

### Phase 3: Migration (Week 4)
- [ ] Update .govman-version
- [ ] Update documentation
- [ ] Deploy to staging
- [ ] Production deployment

### Rollback Plan
- Revert .govman-version to $CURRENT_VERSION
- Redeploy previous version
- Notify team via Slack #dev-team
EOF

echo "Migration plan created: migration-plan.md"
```

## 🔄 Development Workflows

### **Feature Branch Workflow**

```bash
# Developer starts new feature
git checkout -b feature/user-authentication
cd project-root

# Check current project version
cat .govman-version  # Shows: 1.21.1

# Work with project version
govman current       # Should show: Go 1.21.1
go mod tidy
go test ./...

# If feature needs newer Go version
echo "1.22beta1" > .govman-version
govman install 1.22beta1
govman refresh

# Commit version change
git add .govman-version
git commit -m "Upgrade to Go 1.22beta1 for new feature"
```

### **Code Review with Version Checks**

**Reviewer Checklist:**
```bash
#!/bin/bash
# review-checklist.sh
echo "🔍 Code Review Checklist"

# 1. Check Go version changes
if git diff main -- .govman-version | grep -q "^+"; then
    echo "⚠️  Go version changed in this PR"
    echo "Old: $(git show main:.govman-version)"
    echo "New: $(cat .govman-version)"
    echo "✅ Reason documented in PR description?"
fi

# 2. Test with specified version
PROJECT_VERSION=$(cat .govman-version)
govman use "$PROJECT_VERSION"
echo "✅ Using Go $PROJECT_VERSION for review"

# 3. Run compatibility tests
echo "🧪 Running tests..."
go mod tidy
go test ./...
go vet ./...

# 4. Check backward compatibility
if govman list | grep -q "1.20.5"; then
    echo "🔙 Testing backward compatibility..."
    govman use 1.20.5
    go test ./...
    govman use "$PROJECT_VERSION"  # Switch back
fi

echo "✅ Review complete"
```

### **Pair Programming Setup**

```bash
# Both developers sync to same version
echo "Setting up pair programming session..."

# Developer 1 shares their version
DEV1_VERSION=$(govman current --quiet | grep -o '[0-9.]*')
echo "Dev 1 using Go $DEV1_VERSION"

# Developer 2 matches the version
echo "Dev 2 switching to Go $DEV1_VERSION"
govman use "$DEV1_VERSION"

# Verify both are in sync
echo "Both developers now using:"
govman current
```

## 🔍 Code Review Process

### **Automated PR Checks**

**GitHub Actions Workflow:**
```yaml
# .github/workflows/pr-checks.yml
name: PR Checks

on:
  pull_request:
    branches: [main]

jobs:
  go-version-check:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
      with:
        fetch-depth: 0

    - name: Check Go version changes
      run: |
        if git diff origin/main -- .govman-version | grep -q "^+"; then
          echo "Go version changed in this PR"
          echo "::warning::Go version updated - ensure team approval"

          OLD_VERSION=$(git show origin/main:.govman-version)
          NEW_VERSION=$(cat .govman-version)
          echo "Old version: $OLD_VERSION"
          echo "New version: $NEW_VERSION"
        fi

    - name: Install GOVMAN
      run: curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

    - name: Setup Go version
      run: |
        PROJECT_VERSION=$(cat .govman-version)
        govman install $PROJECT_VERSION
        govman use $PROJECT_VERSION

    - name: Validate setup
      run: |
        govman current
        go version

    - name: Run tests
      run: |
        go mod tidy
        go test ./...
        go vet ./...
```

### **Manual Review Guidelines**

**Reviewer Checklist:**
```markdown
## Go Version Review Checklist

### Version Changes
- [ ] Is Go version change documented in PR description?
- [ ] Is the version change necessary for the feature?
- [ ] Are backward compatibility implications addressed?
- [ ] Is team lead approval obtained for version changes?

### Testing
- [ ] Tests pass with new Go version
- [ ] Performance benchmarks included (if relevant)
- [ ] Security implications reviewed
- [ ] Deployment impact assessed

### Documentation
- [ ] README updated if version requirement changed
- [ ] Migration guide provided (if breaking change)
- [ ] Team notified of version change
```

### **Version Change Approval Process**

```bash
# Script for major version changes
#!/bin/bash
# version-change-approval.sh

OLD_VERSION=$(git show main:.govman-version)
NEW_VERSION=$(cat .govman-version)

if [ "$OLD_VERSION" != "$NEW_VERSION" ]; then
    echo "🚨 Go version change detected"
    echo "From: $OLD_VERSION"
    echo "To:   $NEW_VERSION"

    # Check if major version change
    OLD_MAJOR=$(echo $OLD_VERSION | cut -d. -f1-2)
    NEW_MAJOR=$(echo $NEW_VERSION | cut -d. -f1-2)

    if [ "$OLD_MAJOR" != "$NEW_MAJOR" ]; then
        echo "⚠️  MAJOR VERSION CHANGE"
        echo "Requires team lead approval"
        echo "Create RFC document: docs/rfcs/go-version-upgrade.md"
        exit 1
    else
        echo "✅ Minor version change - acceptable"
    fi
fi
```

## 🌿 Branch Management

### **Version per Branch Strategy**

```bash
# Main development branch
git checkout main
echo "1.21.1" > .govman-version  # Stable version

# Experimental branch
git checkout -b experimental
echo "1.22beta1" > .govman-version  # Cutting edge

# Long-term support branch
git checkout -b lts
echo "1.20.5" > .govman-version  # Conservative version

# Feature branch inherits from parent
git checkout main
git checkout -b feature/new-api
# Automatically uses 1.21.1 from main
```

### **Branch Protection Rules**

**GitHub Branch Protection:**
```yaml
# .github/branch-protection.yml
branch_protection_rules:
  main:
    required_status_checks:
      - "go-version-check"
      - "test-matrix"
    required_reviews: 2
    restrictions:
      - "team-leads"

  experimental:
    required_status_checks:
      - "go-version-check"
    required_reviews: 1
```

### **Merge Strategy for Version Changes**

```bash
#!/bin/bash
# smart-merge.sh
# Handle Go version conflicts during merges

merge_branches() {
    local target_branch=$1
    local feature_branch=$2

    # Check for version conflicts
    TARGET_VERSION=$(git show "$target_branch":.govman-version)
    FEATURE_VERSION=$(git show "$feature_branch":.govman-version)

    if [ "$TARGET_VERSION" != "$FEATURE_VERSION" ]; then
        echo "Version conflict detected:"
        echo "Target ($target_branch): $TARGET_VERSION"
        echo "Feature ($feature_branch): $FEATURE_VERSION"

        # Use newer version
        if [ "$(printf '%s\n' "$TARGET_VERSION" "$FEATURE_VERSION" | sort -V | tail -1)" = "$FEATURE_VERSION" ]; then
            echo "Using feature version: $FEATURE_VERSION"
            echo "$FEATURE_VERSION" > .govman-version
        else
            echo "Using target version: $TARGET_VERSION"
            echo "$TARGET_VERSION" > .govman-version
        fi

        # Test with chosen version
        govman use "$(cat .govman-version)"
        go test ./...
    fi
}
```

## 🚀 Release Management

### **Version-Controlled Releases**

```bash
# Release preparation script
#!/bin/bash
# prepare-release.sh

RELEASE_VERSION="v1.2.0"
GO_VERSION="1.21.1"

echo "Preparing release $RELEASE_VERSION"

# 1. Lock Go version for release
echo "$GO_VERSION" > .govman-version
git add .govman-version

# 2. Update changelog
cat > CHANGELOG.md << EOF
# Changelog

## [$RELEASE_VERSION] - $(date +%Y-%m-%d)

### Requirements
- Go $GO_VERSION

### Added
- New features...

### Changed
- Breaking changes...

### Fixed
- Bug fixes...
EOF

# 3. Create release tag
git add CHANGELOG.md
git commit -m "Release $RELEASE_VERSION with Go $GO_VERSION"
git tag -a "$RELEASE_VERSION" -m "Release $RELEASE_VERSION"

echo "✅ Release $RELEASE_VERSION prepared"
```

### **Multi-Environment Deployment**

```bash
# Deployment matrix
environments=(
    "development:1.21.1"
    "staging:1.21.1"
    "production:1.20.5"  # Conservative for production
)

deploy_to_env() {
    local env_name=$1
    local go_version=$2

    echo "Deploying to $env_name with Go $go_version"

    # Set environment-specific version
    govman use "$go_version"

    # Build with specific version
    go build -ldflags "-X main.Version=$RELEASE_VERSION -X main.GoVersion=$go_version" -o "app-$env_name" .

    # Deploy
    case $env_name in
        "development")
            cp "app-$env_name" /opt/dev/app
            ;;
        "staging")
            scp "app-$env_name" staging-server:/opt/app
            ;;
        "production")
            # Blue-green deployment
            scp "app-$env_name" prod-server:/opt/app-new
            ;;
    esac
}

# Deploy to all environments
for env_spec in "${environments[@]}"; do
    IFS=':' read -r env_name go_version <<< "$env_spec"
    deploy_to_env "$env_name" "$go_version"
done
```

## 🏢 Enterprise Integration

### **Corporate Standards**

```yaml
# .govman/corporate-policy.yaml
corporate_policy:
  allowed_versions:
    - "1.19.5"  # LTS support
    - "1.20.5"  # Current stable
    - "1.21.1"  # Latest approved

  restricted_versions:
    - "1.22*"   # Not yet approved
    - "*rc*"    # No release candidates
    - "*beta*"  # No beta versions

  approval_required:
    - major_version_changes
    - security_updates
    - performance_critical_updates

  compliance:
    scanning_required: true
    vulnerability_checks: true
    license_validation: true
```

### **LDAP/SSO Integration**

```bash
#!/bin/bash
# enterprise-auth.sh
# Example: Integrate with corporate systems

check_permissions() {
    local user=$(whoami)
    local requested_version=$1

    # Check against corporate LDAP
    if ldap_check_group "$user" "go-version-admin"; then
        echo "✅ Admin privileges - all versions allowed"
        return 0
    fi

    # Check approved versions
    if echo "$requested_version" | grep -qE "^1\.(19|20|21)\."; then
        echo "✅ Standard version approved"
        return 0
    fi

    echo "❌ Version $requested_version requires approval"
    echo "Contact: platform-team@company.com"
    return 1
}

# Wrapper around govman install
enterprise_install() {
    local version=$1

    if check_permissions "$version"; then
        govman install "$version"
    else
        exit 1
    fi
}
```

### **Audit and Compliance**

```bash
#!/bin/bash
# audit-report.sh
# Generate compliance reports

generate_audit_report() {
    cat > audit-report.txt << EOF
# Go Version Audit Report
Generated: $(date)

## Current Environment
User: $(whoami)
Host: $(hostname)
GOVMAN Version: $(govman --version)

## Installed Versions
$(govman list)

## Current Project
Project: $(basename $(pwd))
Required Version: $(cat .govman-version 2>/dev/null || echo "Not specified")
Active Version: $(govman current --quiet 2>/dev/null || echo "None")

## Security Status
$(govman --verbose list | grep -E "(checksum|security)")

## Compliance
$(check_compliance)
EOF

    echo "Audit report generated: audit-report.txt"
}

check_compliance() {
    # Check against corporate policy
    if [ -f .govman/corporate-policy.yaml ]; then
        echo "✅ Corporate policy file present"
    else
        echo "⚠️ No corporate policy found"
    fi

    # Check for unapproved versions
    govman list --quiet | while read -r version; do
        if echo "$version" | grep -qE "(rc|beta|alpha)"; then
            echo "⚠️ Unapproved version detected: $version"
        fi
    done
}
```

## 🔧 Troubleshooting Team Issues

### **Common Team Problems**

**Problem: Team members have different Go versions**

**Solution:**
```bash
#!/bin/bash
# sync-team-versions.sh
echo "Synchronizing team Go versions..."

# Get project version
PROJECT_VERSION=$(cat .govman-version)
echo "Project requires Go $PROJECT_VERSION"

# Check current version
CURRENT_VERSION=$(govman current --quiet | grep -o '[0-9.]*')

if [ "$CURRENT_VERSION" != "$PROJECT_VERSION" ]; then
    echo "Installing Go $PROJECT_VERSION..."
    govman install "$PROJECT_VERSION"
    govman use "$PROJECT_VERSION"
    echo "✅ Synchronized to Go $PROJECT_VERSION"
else
    echo "✅ Already using correct version"
fi
```

**Problem: CI/CD fails with version mismatches**

**Solution:**
```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Install GOVMAN
      run: curl -fsSL https://raw.githubusercontent.com/sijunda/govman/main/scripts/install.sh | bash

    - name: Use project Go version
      run: |
        if [ -f .govman-version ]; then
          PROJECT_VERSION=$(cat .govman-version)
          govman install $PROJECT_VERSION
          govman use $PROJECT_VERSION
        else
          echo "No .govman-version found, using latest"
          govman install latest
          govman use latest
        fi

    - name: Verify version consistency
      run: |
        echo "GOVMAN says: $(govman current)"
        echo "Go says: $(go version)"
```

**Problem: Shell integration not working for some team members**

**Solution:**
```bash
#!/bin/bash
# fix-shell-integration.sh
echo "Fixing shell integration for team member..."

# Detect shell
SHELL_NAME=$(basename "$SHELL")
echo "Detected shell: $SHELL_NAME"

# Backup current config
cp ~/.${SHELL_NAME}rc ~/.${SHELL_NAME}rc.backup.$(date +%s)

# Remove old GOVMAN integration
sed -i.bak '/# GOVMAN/,/# END GOVMAN/d' ~/.${SHELL_NAME}rc

# Reinstall integration
govman init --force

# Reload shell
case $SHELL_NAME in
    bash)
        source ~/.bashrc
        ;;
    zsh)
        source ~/.zshrc
        ;;
    fish)
        source ~/.config/fish/config.fish
        ;;
esac

echo "✅ Shell integration fixed"
echo "Test with: cd $(pwd) && govman current"
```

### **Team Communication Templates**

**Slack Notification for Version Changes:**
```bash
#!/bin/bash
# notify-team.sh
NEW_VERSION=$(cat .govman-version)
OLD_VERSION=$(git show HEAD~1:.govman-version)

if [ "$NEW_VERSION" != "$OLD_VERSION" ]; then
    curl -X POST -H 'Content-type: application/json' \
    --data "{\"text\":\"🐹 Go version updated in $(basename $(pwd))\n• From: $OLD_VERSION\n• To: $NEW_VERSION\n• Run: \`govman use $NEW_VERSION\`\"}" \
    $SLACK_WEBHOOK_URL
fi
```

**Email Template for Major Updates:**
```text
Subject: Go Version Update Required - Project: [PROJECT_NAME]

Team,

The Go version for [PROJECT_NAME] has been updated:

Previous Version: [OLD_VERSION]
New Version: [NEW_VERSION]

Action Required:
1. Pull latest changes: git pull
2. Update Go version: govman use [NEW_VERSION]
3. Test your current work: go test ./...

Questions? Contact the platform team.

Best regards,
Platform Team
```

---

**These workflows ensure smooth team collaboration while maintaining version consistency and reducing configuration conflicts.**