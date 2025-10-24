# Developer Onboarding

Welcome to the `govman` team! This guide will help you get set up and make your first contribution.

## 1. Philosophy and Goals

Before you start, it's important to understand the project's core philosophy:

-   **Simplicity Over Features**: `govman` should be simple, intuitive, and require zero configuration to use.
-   **Performance Matters**: Operations like downloading and switching versions should be as fast as possible.
-   **Reliability is Key**: The tool must be robust and handle edge cases gracefully (e.g., network errors, corrupted files).
-   **Minimal Dependencies**: We rely on the Go standard library as much as possible to keep the project lean and secure.

## 2. Your First Week: A Checklist

### Day 1: Setup and Exploration

-   [ ] **Clone and Build**: Follow the [Getting Started Guide](getting-started.md) to clone the repository and build the binary.
-   [ ] **Explore the Commands**: Use your local build to try out every command (`install`, `list`, `use`, `clean`, etc.).
-   [ ] **Read the Architecture Docs**: Read the [Architecture Overview](architecture.md) and study the [Architecture Diagrams](architecture-diagrams.md) to understand the high-level structure.

### Day 2: Dive into the Code

-   [ ] **Trace a Command**: Pick a simple command like `govman current` and trace its execution path from `internal/cli/current.go` through `internal/manager/manager.go`.
-   [ ] **Understand Configuration**: Look at `internal/config/config.go` and the `config.yaml` file to see how settings are loaded and used.
-   [ ] **Review the Manager**: The `manager` is the heart of the application. Spend time understanding how it orchestrates different services.

### Day 3: Run the Tests

-   [ ] **Run All Tests**: Run `make test` to execute the unit tests.
-   [ ] **Generate Coverage**: Run `make test-coverage` and open the `coverage/coverage.html` report in your browser. Identify an area with lower test coverage.
-   [ ] **Write a Simple Test**: Find a simple function that is easy to understand and try to add a new test case for it. This is a great way to get comfortable with the testing workflow.

### Day 4-5: Your First Contribution

-   [ ] **Find an Issue**: Look for issues on GitHub tagged with `good first issue` or `documentation`.
-   [ ] **Fix a Small Bug or Add a Docstring**: Even a small contribution is valuable. It could be correcting a typo, improving a comment, or adding a missing test case.
-   [ ] **Submit a Pull Request**: Follow the steps in the [Getting Started Guide](getting-started.md) to submit your PR. Don't worry if it's not perfect; the review process is a collaborative effort.

## 3. Key Areas of the Codebase

-   **Adding a new command?** Start in `internal/cli/`. You'll need to create a new `cobra` command and add it to `internal/cli/command.go`.
-   **Changing core logic?** Most business logic lives in `internal/manager/manager.go`.
-   **Modifying shell integration?** Look at `internal/shell/shell.go`. Each shell has its own struct (`ZshShell`, `PowerShell`, etc.) that implements the `Shell` interface.
-   **Improving download performance?** The `internal/downloader/downloader.go` file contains the download and extraction engine.

## 4. Asking for Help

If you get stuck, don't hesitate to ask for help! You can:
-   Open an issue on GitHub.
-   Leave a comment on an existing issue or pull request.
-   Reach out to one of the maintainers directly.

We are excited to have you on board!