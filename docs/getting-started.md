# Getting Started (Developer Guide)

This guide is for developers who want to contribute to `govman` or understand its inner workings.

## Prerequisites

-   Go 1.25 or later
-   Git
-   Make (optional, but recommended for using the Makefile)

## 1. Fork and Clone the Repository

First, create a fork of the `sijunda/govman` repository on GitHub.

Then, clone your forked repository to your local machine:

```bash
git clone https://github.com/YOUR_USERNAME/govman.git
cd govman
```

## 2. Set Up the Development Environment

The project includes a `Makefile` target to automate the setup process. This will install all required Go dependencies and development tools.

```bash
make dev-setup
```

This command will:
-   Download and verify Go modules (`go mod download`, `go mod verify`).
-   Install essential development tools like `golangci-lint`, `gosec`, and `goimports`.

## 3. Build the Binary

You can build the `govman` binary for your local platform using the `build` target:

```bash
make build
```

The compiled binary will be placed in the `build/` directory. You can run it directly:

```bash
./build/govman --version
```

## 4. Running Tests

`govman` has a comprehensive test suite. To run all unit tests:

```bash
make test
```

To run tests with coverage and generate an HTML report:

```bash
make test-coverage
```
The coverage report will be saved at `coverage/coverage.html`.

## 5. Code Quality and Linting

Before committing code, ensure it meets the project's quality standards by running the validation suite:

```bash
make validate
```

This command is a convenient shortcut that runs formatting (`fmt`), static analysis (`vet`), and linting (`lint`) in one step.

## 6. Making Changes

1.  **Create a new branch**:
    ```bash
    git checkout -b feature/my-new-feature
    ```

2.  **Write your code**:
    -   Follow the existing code style.
    -   Add or update tests for your changes.
    -   Update any relevant documentation in the `docs/` directory.

3.  **Run checks**:
    ```bash
    make check  # Runs 'validate' and 'test'
    ```

4.  **Commit your changes**:
    Use a clear and descriptive commit message.
    ```bash
    git commit -m "feat: Add support for new shell"
    ```

5.  **Push to your fork**:
    ```bash
    git push origin feature/my-new-feature
    ```

6.  **Open a Pull Request**:
    Go to the original `govman` repository on GitHub and open a pull request from your forked branch. Provide a clear description of your changes.