# PromptAI for Generating GOVMAN Landing Page

## Introduction
Generate a comprehensive and SEO-friendly landing page for GOVMAN - Go Version Manager. The landing page should provide a clear overview of the project, its features, installation instructions, and usage examples. Ensure the page includes meta tags, alt text for images, and a well-structured URL.

## SEO Considerations
- **Meta Tags:** Include title, description, and keywords.
- **Alt Text:** Provide alt text for all images.
- **URL Structure:** Use a clean and descriptive URL structure.
- **Header Tags:** Use H1, H2, and H3 tags appropriately.
- **Content Quality:** Write high-quality, informative, and engaging content.
- **Internal Links:** Include internal links to other relevant pages within the site.
- **Mobile Responsiveness:** Ensure the page is mobile-friendly.
- **Loading Speed:** Optimize images and assets for faster loading times.
- **Schema Markup:** Add schema markup for better search engine understanding.

## Key Features
- **Version Management:** Install, uninstall, and switch between different Go versions.
- **Project-Specific Versions:** Support for defining Go versions on a per-project basis using a `.govman-version` file.
- **Cross-Platform Compatibility:** Works across Windows, macOS, and Linux operating systems.
- **Automatic Shell Integration:** Seamless integration with various shell environments.
- **Efficient Downloads:** Utilizes fast parallel downloads with resume capabilities.

## Quick Installation
The easiest way to install `govman` on Unix-like systems (Linux, macOS, FreeBSD) is by using the provided `install.sh` script:

## Getting Started
Provide a brief guide on how to get started with GOVMAN. Include initial setup instructions and basic usage examples.

## Documentation
Include detailed documentation sections covering:
- **Installation:** Step-by-step installation instructions.
- **Configuration:** Configuration options and settings.
- **Usage:** Detailed usage examples and command explanations.
- **Troubleshooting:** Common issues and solutions.

## API Reference
Provide a comprehensive API reference if applicable, detailing:
- **Commands:** Available commands and their usage.
- **Flags:** Command flags and their descriptions.
- **Environment Variables:** Environment variables used by GOVMAN.

## FAQ
Answer frequently asked questions about GOVMAN, including:
- **What is GOVMAN?**
- **How do I install GOVMAN?**
- **Can I use GOVMAN on Windows?**
- **How do I report issues or contribute to the project?**

## Contact
Provide contact information for users to reach out:
- **Email:** [support@govman.io](mailto:support@govman.io)
- **GitHub Issues:** [https://github.com/sijunda/govman/issues](https://github.com/sijunda/govman/issues)
- **Discord:** [Join our Discord channel](https://discord.gg/govman)

## Meta Tags
- **Title:** GOVMAN - Go Version Manager | Simplify Go Version Management
- **Description:** GOVMAN is a cross-platform CLI tool for managing multiple Go versions. Install, switch, and manage Go versions easily.
- **Keywords:** GOVMAN, Go Version Manager, Go, CLI, version management, cross-platform, shell integration, efficient downloads

## Alt Text
Ensure all images have descriptive alt text. For example:
- ![GOVMAN Logo](logo.png "GOVMAN Logo")

## URL Structure
Use a clean and descriptive URL structure. For example:
- `https://govman.io/`
- `https://govman.io/features`
- `https://govman.io/installation`

## Header Tags
Use H1, H2, and H3 tags appropriately. For example:
- `# GOVMAN - Go Version Manager`
- `## Quick Installation`
- `### Prerequisites`

## Content Quality
Write high-quality, informative, and engaging content. Ensure the page is well-written and easy to understand.

## Internal Links
Include internal links to other relevant pages within the site. For example:
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)

## Mobile Responsiveness
Ensure the page is mobile-friendly. Use responsive design techniques to ensure the page looks good on all devices.

## Loading Speed
Optimize images and assets for faster loading times. Use techniques such as image compression and lazy loading.

## Schema Markup
Add schema markup for better search engine understanding. For example:
- Use JSON-LD to mark up the page structure and content.
- Include schema for software application, organization, and product.
```bash
curl -sSL https://raw.githubusercontent.com/sijunda/govman/main/install.sh | bash
```
This script will download the latest stable release of `govman` for your system, install it to `$HOME/.govman/bin`, and add it to your shell's `PATH`.

## Building and Running
### Prerequisites
- Go (version 1.25.1 or higher, as specified in `go.mod`)
- `git` (for version information in builds)
- `curl` (for the `install.sh` script)

### Development Setup
To set up the development environment and install necessary tools:
```bash
make dev-setup
```
This command will download Go module dependencies and install development tools like `golangci-lint`, `goreleaser`, `goimports`, etc.

### Dependency Management
To download and verify Go module dependencies:
```bash
make deps
```

### Building
- **Build for the current platform:**
  ```bash
  make build
  ```
  The executable will be placed in the `build/` directory.
- **Build for all supported platforms:**
  ```bash
  make build-all
  ```
  Cross-compiled binaries will be placed in the `dist/` directory.

### Installation (from source)
- **Install to your GOPATH/bin:**
  ```bash
  make install
  ```
- **Install to `/usr/local/bin` (requires `sudo`):**
  ```bash
  make install-local
  ```

### Running
After installation, you can run `govman` commands from your terminal:
```bash
govman --help
```

## Usage
Here are some common `govman` commands:
- **Install a specific Go version:**
  ```bash
  govman install 1.20.1
  govman install latest
  ```
- **Switch to a specific Go version:**
  ```bash
  govman use 1.20.1
  ```
- **Set a Go version as default:**
  ```bash
  govman use 1.20.1 --default
  ```
- **Set a project-specific Go version:**
  ```bash
  govman use 1.20.1 --local
  ```
- **List installed Go versions:**
  ```bash
  govman list
  ```
- **List available Go versions for download:**
  ```bash
  govman list --remote
  ```
- **Uninstall a Go version:**
  ```bash
  govman uninstall 1.20.1
  ```
- **Show current active Go version:**
  ```bash
  govman current
  ```
- **Clean cached Go archives:**
  ```bash
  govman clean
  ```

## Configuration
`govman` uses `viper` for configuration. The default configuration file is located at `$HOME/.govman/config.yaml`. You can also specify a custom configuration file using the `--config` flag.

## Testing
The project includes various test targets:
- **Run unit tests:**
  ```bash
  make test
  ```
- **Run tests with coverage analysis:**
  ```bash
  make test-coverage
  ```
- **Run integration tests:**
  ```bash
  make test-integration
  ```
- **Run all tests (unit, integration, benchmark):**
  ```bash
  make test-all
  ```

## Code Quality and Linting
- **Format code using `goimports`:**
  ```bash
  make fmt
  ```
- **Run `go vet` for static analysis:**
  ```bash
  make vet
  ```
- **Run comprehensive linting using `golangci-lint`:**
  ```bash
  make lint
  ```
- **Run all validation checks (fmt, vet, lint):**
  ```bash
  make validate
  ```

## Development Conventions
- **Build System:** Uses a `Makefile` for consistent task automation.
- **Code Formatting:** Enforces code formatting using `goimports`.
- **Linting:** Utilizes `golangci-lint` for comprehensive code quality checks.
- **Static Analysis:** Employs `go vet` for identifying suspicious constructs.
- **Testing:** Follows standard Go testing practices with `_test.go` files. Integration tests are specifically tagged.
- **CLI Framework:** Built with the Cobra library for structured command-line interfaces.
- **Configuration:** Uses Viper for flexible application configuration.
- **Version Information:** Build-time version, commit, and branch information are injected into the binary using `ldflags`.

## Contributing
Contributions are welcome! Please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project. (Note: `CONTRIBUTING.md` is a placeholder and needs to be created if not present).

## License
This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.