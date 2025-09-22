# Makefile for govman - Enhanced cross-platform build system
.PHONY: build build-all build-binaries build-archives test test-coverage test-integration test-benchmark clean lint fmt vet install dev-setup help deps check release docker version generate validate security

# ==================================================================================
# BUILD CONFIGURATION
# ==================================================================================

# Version and build info
# Priority: 1. Manual VERSION env var, 2. Git tag (exact), 3. Git describe, 4. Fallback
VERSION ?= $(shell \
	if [ -n "$$VERSION" ]; then \
		echo "$$VERSION"; \
	elif git describe --tags --exact-match HEAD >/dev/null 2>&1; then \
		git describe --tags --exact-match HEAD; \
	elif git describe --tags >/dev/null 2>&1; then \
		git describe --tags --always --dirty; \
	else \
		echo "v0.0.0-dev"; \
	fi)

COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_BY ?= $(shell whoami)@$(shell hostname)
GO_VERSION ?= $(shell go version | cut -d' ' -f3)

# Go configuration
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
MAIN_PACKAGE = ./cmd/govman
BINARY_NAME = govman
MODULE_NAME = github.com/sijunda/govman

# Directories
BUILD_DIR = build
DIST_DIR = dist
COVERAGE_DIR = coverage
DOCS_DIR = docs
SCRIPTS_DIR = scripts
TOOLS_DIR = tools

# Colors for output
CYAN = \033[36m
GREEN = \033[32m
YELLOW = \033[33m
RED = \033[31m
RESET = \033[0m
BOLD = \033[1m

# FIXED: Better linker flags with version path
LDFLAGS = -ldflags "\
	-s -w \
	-X '$(MODULE_NAME)/internal/version.Version=$(VERSION)' \
	-X '$(MODULE_NAME)/internal/version.Commit=$(COMMIT)' \
	-X '$(MODULE_NAME)/internal/version.Date=$(DATE)' \
	-X '$(MODULE_NAME)/internal/version.BuildBy=$(BUILD_BY)'"

# Build tags
BUILD_TAGS ?= netgo
CGO_ENABLED ?= 0

# Test flags
TEST_FLAGS ?= -race -timeout=10m
TEST_COVERAGE_FLAGS ?= $(TEST_FLAGS) -coverprofile=coverage.out -covermode=atomic

# Platform definitions with comprehensive architecture support
PLATFORMS = \
	linux/amd64 \
	linux/arm64 \
	linux/arm \
	linux/386 \
	linux/mips64 \
	linux/mips64le \
	linux/mips \
	linux/mipsle \
	linux/ppc64le \
	linux/ppc64 \
	linux/s390x \
	linux/riscv64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64 \
	windows/386 \
	freebsd/amd64 \
	freebsd/arm64 \
	freebsd/arm \
	freebsd/386 \
	netbsd/amd64 \
	netbsd/arm64 \
	netbsd/arm \
	netbsd/386 \
	openbsd/amd64 \
	openbsd/arm64 \
	openbsd/arm \
	openbsd/386 \
	dragonfly/amd64 \
	solaris/amd64 \
	aix/ppc64

# ==================================================================================
# HELP TARGET
# ==================================================================================

help: ## Show this help message
	@echo "$(BOLD)$(CYAN)GOVMAN - Go Version Manager$(RESET)"
	@echo "$(CYAN)=================================$(RESET)"
	@echo ""
	@echo "$(BOLD)Usage:$(RESET) make $(YELLOW)<target>$(RESET)"
	@echo ""
	@echo "$(BOLD)Build Information:$(RESET)"
	@echo "  Version: $(GREEN)$(VERSION)$(RESET)"
	@echo "  Commit:  $(GREEN)$(COMMIT)$(RESET)"
	@echo "  Branch:  $(GREEN)$(BRANCH)$(RESET)"
	@echo "  Go:      $(GREEN)$(GO_VERSION)$(RESET)"
	@echo ""
	@echo "$(BOLD)Available Targets:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z0-9_-]+:.*?## / {printf "  $(YELLOW)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

# ==================================================================================
# DEVELOPMENT SETUP
# ==================================================================================

dev-setup: ## Set up development environment
	@echo "$(CYAN)🔧 Setting up development environment...$(RESET)"
	@echo "$(YELLOW)📦 Installing Go dependencies...$(RESET)"
	go mod download
	go mod verify
	@echo "$(YELLOW)🔧 Installing development tools...$(RESET)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/goreleaser/goreleaser@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/kisielk/errcheck@latest
	@echo "$(YELLOW)📁 Creating directories...$(RESET)"
	mkdir -p $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_DIR) $(DOCS_DIR) $(TOOLS_DIR)
	@echo "$(GREEN)✅ Development environment ready!$(RESET)"

deps: ## Download and verify dependencies
	@echo "$(CYAN)📦 Managing dependencies...$(RESET)"
	go mod download
	go mod verify
	go mod tidy
	@echo "$(GREEN)✅ Dependencies updated!$(RESET)"

# ==================================================================================
# VERSION DEBUGGING TARGETS
# ==================================================================================

debug-version: ## Debug version detection
	@echo "$(BOLD)$(CYAN)Version Detection Debug:$(RESET)"
	@echo "  Manual VERSION env: $(GREEN)$${VERSION}$(RESET)"
	@echo "  Git exact tag: $(GREEN)$(shell git describe --tags --exact-match HEAD 2>/dev/null || echo 'none')$(RESET)"
	@echo "  Git describe: $(GREEN)$(shell git describe --tags --always --dirty 2>/dev/null || echo 'none')$(RESET)"
	@echo "  Git status clean: $(GREEN)$(shell if [ -z "$$(git status --porcelain 2>/dev/null)" ]; then echo 'yes'; else echo 'no'; fi)$(RESET)"
	@echo "  Current VERSION: $(GREEN)$(VERSION)$(RESET)"
	@echo ""
	@echo "$(BOLD)$(CYAN)Git Information:$(RESET)"
	@echo "  Current branch: $(GREEN)$(BRANCH)$(RESET)"
	@echo "  Last commit: $(GREEN)$(COMMIT)$(RESET)"
	@echo "  Available tags:"
	@git tag -l | tail -10 | sed 's/^/    /'

check-git-status: ## Check current git status
	@echo "$(CYAN)🔍 Checking git status...$(RESET)"
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "$(RED)❌ Working directory is dirty:$(RESET)"; \
		git status --short; \
		echo "$(YELLOW)💡 Commit your changes to remove '-dirty' suffix$(RESET)"; \
	else \
		echo "$(GREEN)✅ Working directory is clean$(RESET)"; \
	fi

# ==================================================================================
# CODE QUALITY
# ==================================================================================

fmt: ## Format code with goimports
	@echo "$(CYAN)📝 Formatting code...$(RESET)"
	@if command -v goimports >/dev/null 2>&1; then \
		echo "Using goimports..."; \
		goimports -w -local $(MODULE_NAME) .; \
	else \
		echo "Using go fmt..."; \
		go fmt ./...; \
	fi
	@echo "$(GREEN)✅ Code formatted!$(RESET)"

vet: ## Run go vet
	@echo "$(CYAN)🔍 Running go vet...$(RESET)"
	go vet ./...
	@echo "$(GREEN)✅ Go vet passed!$(RESET)"

lint: ## Run comprehensive linting
	@echo "$(CYAN)🔍 Running linter...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not found, running basic checks...$(RESET)"; \
		go vet ./...; \
		if command -v staticcheck >/dev/null 2>&1; then staticcheck ./...; fi; \
		if command -v errcheck >/dev/null 2>&1; then errcheck ./...; fi; \
	fi
	@echo "$(GREEN)✅ Linting completed!$(RESET)"

security: ## Run security checks
	@echo "$(CYAN)🔒 Running security checks...$(RESET)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec -fmt sarif -out gosec-report.sarif -stdout ./...; \
	else \
		echo "$(YELLOW)⚠️  gosec not found, install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest$(RESET)"; \
	fi

validate: fmt vet lint ## Run all validation checks
	@echo "$(GREEN)✅ All validation checks passed!$(RESET)"

# ==================================================================================
# TESTING
# ==================================================================================

test: ## Run unit tests
	@echo "$(CYAN)🧪 Running unit tests...$(RESET)"
	go test $(TEST_FLAGS) -tags=$(BUILD_TAGS) ./...
	@echo "$(GREEN)✅ Unit tests passed!$(RESET)"

test-coverage: ## Run tests with coverage analysis
	@echo "$(CYAN)🧪 Running tests with coverage...$(RESET)"
	mkdir -p $(COVERAGE_DIR)
	go test $(TEST_COVERAGE_FLAGS) -tags=$(BUILD_TAGS) ./...
	go tool cover -html=coverage.out -o $(COVERAGE_DIR)/coverage.html
	go tool cover -func=coverage.out | tail -1 | awk '{print "Coverage: " $$3}'
	@echo "$(GREEN)📊 Coverage report: $(COVERAGE_DIR)/coverage.html$(RESET)"

test-integration: ## Run integration tests
	@echo "$(CYAN)🧪 Running integration tests...$(RESET)"
	go test $(TEST_FLAGS) -tags=integration ./test/integration/...
	@echo "$(GREEN)✅ Integration tests passed!$(RESET)"

test-benchmark: ## Run benchmark tests
	@echo "$(CYAN)🏃 Running benchmark tests...$(RESET)"
	go test -bench=. -benchmem -tags=$(BUILD_TAGS) ./...
	@echo "$(GREEN)✅ Benchmarks completed!$(RESET)"

test-all: test test-integration test-benchmark ## Run all tests

# ==================================================================================
# BUILDING - WITH BETTER VERSION HANDLING
# ==================================================================================

build: debug-version ## Build for current platform
	@echo "$(CYAN)🏗️  Building govman for $(GOOS)/$(GOARCH) with version $(VERSION)...$(RESET)"
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build $(LDFLAGS) -tags=$(BUILD_TAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@if [ "$(GOOS)" = "windows" ]; then mv $(BUILD_DIR)/$(BINARY_NAME) $(BUILD_DIR)/$(BINARY_NAME).exe; fi
	@echo "$(GREEN)✅ Built: $(BUILD_DIR)/$(BINARY_NAME) ($(VERSION))$(RESET)"
	@echo "$(YELLOW)🔍 Testing version output:$(RESET)"
	@$(BUILD_DIR)/$(BINARY_NAME) --version 2>/dev/null || echo "$(YELLOW)⚠️  Binary doesn't support --version flag$(RESET)"

build-binary: ## Build binary for current platform
	@echo "$(CYAN)🏗️  Building binary for $(GOOS)/$(GOARCH)...$(RESET)"
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build $(LDFLAGS) -tags=$(BUILD_TAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@if [ "$(GOOS)" = "windows" ]; then mv $(BUILD_DIR)/$(BINARY_NAME) $(BUILD_DIR)/$(BINARY_NAME).exe; fi
	@echo "$(GREEN)✅ Binary built: $(BUILD_DIR)/$(BINARY_NAME)$(RESET)"

build-debug: ## Build with debug information
	@echo "$(CYAN)🏗️  Building govman with debug info...$(RESET)"
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build -gcflags="all=-N -l" -tags=$(BUILD_TAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-debug $(MAIN_PACKAGE)
	@echo "$(GREEN)✅ Debug build: $(BUILD_DIR)/$(BINARY_NAME)-debug$(RESET)"

# Build with explicit version
build-release: ## Build release version with explicit version check
	@echo "$(CYAN)🚀 Building release version...$(RESET)"
	@if [ -z "$(VERSION)" ] || echo "$(VERSION)" | grep -q "dirty"; then \
		echo "$(RED)❌ Cannot build release with dirty or missing version$(RESET)"; \
		echo "$(YELLOW)💡 Clean your git status and ensure you're on a tagged commit$(RESET)"; \
		exit 1; \
	fi
	$(MAKE) build
	@echo "$(GREEN)🚀 Release build completed: $(VERSION)$(RESET)"

build-all: build-binaries ## Build binaries for all supported platforms (alias)

build-binaries: ## Build binaries for all supported platforms
	@echo "$(CYAN)🏗️  Building binaries for all platforms with version $(VERSION)...$(RESET)"
	@rm -rf $(DIST_DIR)
	@mkdir -p $(DIST_DIR)
	@total=$$(echo "$(PLATFORMS)" | wc -w); \
	current=0; \
	for platform in $(PLATFORMS); do \
		current=$$((current + 1)); \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		binary_name=$(BINARY_NAME)-$$os-$$arch; \
		if [ "$$os" = "windows" ]; then binary_name=$$binary_name.exe; fi; \
		echo "$(YELLOW)[$$current/$$total] Building for $$os/$$arch...$(RESET)"; \
		if CGO_ENABLED=$(CGO_ENABLED) GOOS=$$os GOARCH=$$arch \
			go build $(LDFLAGS) -tags=$(BUILD_TAGS) -o $(DIST_DIR)/$$binary_name $(MAIN_PACKAGE) 2>/dev/null; then \
			echo "$(GREEN)✅ $$binary_name$(RESET)"; \
		else \
			echo "$(RED)❌ Failed to build for $$os/$$arch$(RESET)"; \
		fi; \
	done
	@echo "$(GREEN)🎉 All binaries built with version $(VERSION)! Check $(DIST_DIR)/$(RESET)"
	@ls -la $(DIST_DIR)/

build-archives: build-all ## Build archives for distribution
	@echo "$(CYAN)📦 Creating distribution archives...$(RESET)"
	@cd $(DIST_DIR) && \
	for binary in govman-*; do \
		if [ -f "$$binary" ]; then \
			platform=$$(echo $$binary | sed 's/govman-//; s/\.exe$$//'); \
			echo "$(YELLOW)📦 Creating archive for $$platform...$(RESET)"; \
			if echo "$$binary" | grep -q "windows"; then \
				if [ -f "../README.md" ] && [ -f "../LICENSE.md" ]; then \
					zip "$$platform.zip" "$$binary" ../README.md ../LICENSE.md; \
				elif [ -f "../README.md" ] && [ -f "../LICENSE" ]; then \
					zip "$$platform.zip" "$$binary" ../README.md ../LICENSE; \
				else \
					zip "$$platform.zip" "$$binary"; \
				fi; \
			else \
				if [ -f "../README.md" ] && [ -f "../LICENSE.md" ]; then \
					tar -czf "$$platform.tar.gz" "$$binary" ../README.md ../LICENSE.md; \
				elif [ -f "../README.md" ] && [ -f "../LICENSE" ]; then \
					tar -czf "$$platform.tar.gz" "$$binary" ../README.md ../LICENSE; \
				else \
					tar -czf "$$platform.tar.gz" "$$binary"; \
				fi; \
			fi; \
		fi; \
	done
	@echo "$(GREEN)📦 Archives created in $(DIST_DIR)/$(RESET)"

# ==================================================================================
# INSTALLATION
# ==================================================================================

install: build ## Install to GOPATH/bin
	@echo "$(CYAN)📦 Installing govman to GOPATH/bin...$(RESET)"
	go install $(LDFLAGS) -tags=$(BUILD_TAGS) $(MAIN_PACKAGE)
	@echo "$(GREEN)✅ govman installed to $$(go env GOPATH)/bin/govman$(RESET)"

install-local: build ## Install to /usr/local/bin (requires sudo)
	@echo "$(CYAN)📦 Installing govman to /usr/local/bin...$(RESET)"
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✅ govman installed to /usr/local/bin/govman$(RESET)"

uninstall: ## Uninstall from system
	@echo "$(CYAN)🗑️  Uninstalling govman...$(RESET)"
	@if [ -f "/usr/local/bin/govman" ]; then sudo rm /usr/local/bin/govman; echo "$(GREEN)✅ Removed from /usr/local/bin$(RESET)"; fi
	@if [ -f "$$(go env GOPATH)/bin/govman" ]; then rm "$$(go env GOPATH)/bin/govman"; echo "$(GREEN)✅ Removed from GOPATH/bin$(RESET)"; fi

# ==================================================================================
# RELEASE MANAGEMENT
# ==================================================================================

check-git-clean: ## Check if git working directory is clean
	@echo "$(CYAN)🔍 Checking git status...$(RESET)"
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "$(RED)❌ Working directory is not clean. Please commit or stash changes.$(RESET)"; \
		git status --short; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Working directory is clean$(RESET)"

check-git-tag: ## Check if current commit is tagged
	@echo "$(CYAN)🏷️  Checking git tag...$(RESET)"
	@if ! git describe --tags --exact-match HEAD >/dev/null 2>&1; then \
		echo "$(RED)❌ Current commit is not tagged. Please create and push a tag first.$(RESET)"; \
		echo "$(YELLOW)Example: make tag TAG=v1.0.0$(RESET)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Current commit is tagged: $(shell git describe --tags --exact-match HEAD)$(RESET)"

pre-release-checks: check-git-clean check-git-tag validate test ## Run all pre-release checks
	@echo "$(GREEN)✅ All pre-release checks passed!$(RESET)"

release: pre-release-checks ## Build release with goreleaser
	@echo "$(CYAN)🚀 Building release...$(RESET)"
	@if ! command -v goreleaser >/dev/null 2>&1; then \
		echo "$(RED)❌ goreleaser not found. Install with: go install github.com/goreleaser/goreleaser@latest$(RESET)"; \
		exit 1; \
	fi
	goreleaser release --clean
	@echo "$(GREEN)🚀 Release completed!$(RESET)"

release-snapshot: ## Build snapshot release
	@echo "$(CYAN)📸 Building snapshot release...$(RESET)"
	goreleaser release --snapshot --clean
	@echo "$(GREEN)📸 Snapshot release completed!$(RESET)"

release-dry-run: ## Test release process
	@echo "$(CYAN)🧪 Testing release process...$(RESET)"
	goreleaser release --skip-publish --clean
	@echo "$(GREEN)🧪 Dry run completed!$(RESET)"

tag: ## Create and push a new tag
	@if [ -z "$(TAG)" ]; then \
		echo "$(RED)❌ Please specify TAG: make tag TAG=v1.0.0$(RESET)"; \
		exit 1; \
	fi
	@echo "$(CYAN)🏷️  Creating and pushing tag $(TAG)...$(RESET)"
	git tag -a $(TAG) -m "Release $(TAG)"
	git push origin $(TAG)
	@echo "$(GREEN)🏷️  Tag $(TAG) created and pushed!$(RESET)"

# ==================================================================================
# UTILITIES
# ==================================================================================

version: debug-version ## Show version information
	@echo ""
	@echo "$(BOLD)$(CYAN)Compiled Version Information:$(RESET)"
	@if [ -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		echo "  From Binary: $(GREEN)$$($(BUILD_DIR)/$(BINARY_NAME) --version 2>/dev/null || echo 'N/A')$(RESET)"; \
	else \
		echo "  From Binary: $(YELLOW)No binary found - run 'make build' first$(RESET)"; \
	fi

info: version ## Show detailed build information
	@echo ""
	@echo "$(BOLD)$(CYAN)Environment:$(RESET)"
	@echo "  GOPATH:    $(GREEN)$$(go env GOPATH)$(RESET)"
	@echo "  GOROOT:    $(GREEN)$$(go env GOROOT)$(RESET)"
	@echo "  GOPROXY:   $(GREEN)$$(go env GOPROXY)$(RESET)"
	@echo "  CGO:       $(GREEN)$(CGO_ENABLED)$(RESET)"
	@echo ""
	@echo "$(BOLD)$(CYAN)Build Settings:$(RESET)"
	@echo "  Tags:      $(GREEN)$(BUILD_TAGS)$(RESET)"
	@echo "  LDFLAGS:   $(GREEN)$(LDFLAGS)$(RESET)"

clean: ## Clean build artifacts
	@echo "$(CYAN)🧹 Cleaning up...$(RESET)"
	rm -rf $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	rm -f coverage.out coverage.html gosec-report.sarif
	go clean -cache -testcache -modcache
	@echo "$(GREEN)✅ Cleanup completed!$(RESET)"

check: validate test ## Run all quality checks and tests
	@echo "$(GREEN)🎉 All checks passed successfully!$(RESET)"

ci: deps generate validate test build ## Run full CI pipeline
	@echo "$(GREEN)🎉 CI pipeline completed successfully!$(RESET)"

# ==================================================================================
# PLATFORM-SPECIFIC TARGETS
# ==================================================================================

build-linux: ## Build for Linux (amd64, arm64, arm)
	@echo "$(CYAN)🐧 Building for Linux...$(RESET)"
	@for arch in amd64 arm64 arm 386; do \
		echo "$(YELLOW)Building for linux/$$arch...$(RESET)"; \
		mkdir -p $(DIST_DIR); \
		CGO_ENABLED=0 GOOS=linux GOARCH=$$arch \
			go build $(LDFLAGS) -tags=$(BUILD_TAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-$$arch $(MAIN_PACKAGE); \
	done
	@echo "$(GREEN)✅ Linux builds completed!$(RESET)"

build-darwin: ## Build for macOS (amd64, arm64)
	@echo "$(CYAN)🍎 Building for macOS...$(RESET)"
	@for arch in amd64 arm64; do \
		echo "$(YELLOW)Building for darwin/$$arch...$(RESET)"; \
		mkdir -p $(DIST_DIR); \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$$arch \
			go build $(LDFLAGS) -tags=$(BUILD_TAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-$$arch $(MAIN_PACKAGE); \
	done
	@echo "$(GREEN)✅ macOS builds completed!$(RESET)"

build-windows: ## Build for Windows (amd64, 386, arm64)
	@echo "$(CYAN)🪟 Building for Windows...$(RESET)"
	@for arch in amd64 386 arm64; do \
		echo "$(YELLOW)Building for windows/$$arch...$(RESET)"; \
		mkdir -p $(DIST_DIR); \
		CGO_ENABLED=0 GOOS=windows GOARCH=$$arch \
			go build $(LDFLAGS) -tags=$(BUILD_TAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-$$arch.exe $(MAIN_PACKAGE); \
	done
	@echo "$(GREEN)✅ Windows builds completed!$(RESET)"

build-freebsd: ## Build for FreeBSD (amd64, arm64, arm, 386)
	@echo "$(CYAN)😈 Building for FreeBSD...$(RESET)"
	@for arch in amd64 arm64 arm 386; do \
		echo "$(YELLOW)Building for freebsd/$$arch...$(RESET)"; \
		mkdir -p $(DIST_DIR); \
		CGO_ENABLED=0 GOOS=freebsd GOARCH=$$arch \
			go build $(LDFLAGS) -tags=$(BUILD_TAGS) -o $(DIST_DIR)/$(BINARY_NAME)-freebsd-$$arch $(MAIN_PACKAGE); \
	done
	@echo "$(GREEN)✅ FreeBSD builds completed!$(RESET)"

# ==================================================================================
# ADVANCED FEATURES
# ==================================================================================

profile: ## Build with profiling support
	@echo "$(CYAN)📊 Building with profiling support...$(RESET)"
	mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -tags="$(BUILD_TAGS) pprof" -o $(BUILD_DIR)/$(BINARY_NAME)-profile $(MAIN_PACKAGE)
	@echo "$(GREEN)✅ Profiling build: $(BUILD_DIR)/$(BINARY_NAME)-profile$(RESET)"

size-analysis: build ## Analyze binary size
	@echo "$(CYAN)📏 Analyzing binary size...$(RESET)"
	@if [ -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		ls -lah $(BUILD_DIR)/$(BINARY_NAME); \
		if command -v du >/dev/null 2>&1; then du -h $(BUILD_DIR)/$(BINARY_NAME); fi; \
		if command -v file >/dev/null 2>&1; then file $(BUILD_DIR)/$(BINARY_NAME); fi; \
	fi

watch: ## Watch for changes and rebuild
	@echo "$(CYAN)👀 Watching for changes...$(RESET)"
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o . | xargs -n1 -I{} make build; \
	elif command -v inotifywait >/dev/null 2>&1; then \
		while inotifywait -r -e modify .; do make build; done; \
	else \
		echo "$(RED)❌ No file watcher found. Install fswatch or inotify-tools$(RESET)"; \
	fi

# ==================================================================================
# MAINTENANCE
# ==================================================================================

update-deps: ## Update all dependencies
	@echo "$(CYAN)📦 Updating dependencies...$(RESET)"
	go get -u ./...
	go mod tidy
	@echo "$(GREEN)✅ Dependencies updated!$(RESET)"

outdated: ## Show outdated dependencies
	@echo "$(CYAN)📦 Checking for outdated dependencies...$(RESET)"
	@go list -u -m all | grep '\['

# Show current configuration
show-config:
	@echo "$(BOLD)$(CYAN)Current Configuration:$(RESET)"
	@echo "  GOOS:          $(GREEN)$(GOOS)$(RESET)"
	@echo "  GOARCH:        $(GREEN)$(GOARCH)$(RESET)"
	@echo "  CGO_ENABLED:   $(GREEN)$(CGO_ENABLED)$(RESET)"
	@echo "  BUILD_TAGS:    $(GREEN)$(BUILD_TAGS)$(RESET)"
	@echo "  BUILD_DIR:     $(GREEN)$(BUILD_DIR)$(RESET)"
	@echo "  DIST_DIR:      $(GREEN)$(DIST_DIR)$(RESET)"