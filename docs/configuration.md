# Configuration

`govman` is designed to work out-of-the-box with zero configuration. However, you can customize its behavior by editing the `config.yaml` file located at `~/.govman/config.yaml`.

If the file does not exist, `govman` will create it with default values the first time it runs.

## Default Configuration

Here is an example of the default configuration file:

```yaml
# Installation directory for Go versions
install_dir: "~/.govman/versions"

# Cache directory for downloads
cache_dir: "~/.govman/cache"

# Default Go version (empty = none)
default_version: ""

# Download configuration
download:
  parallel: true          # Enable parallel downloads
  max_connections: 4      # Maximum concurrent connections
  timeout: 300s           # Download timeout
  retry_count: 3          # Number of retry attempts
  retry_delay: 5s         # Delay between retries

# Mirror configuration (for users in China or with network restrictions)
mirror:
  enabled: false
  url: "https://golang.google.cn/dl/"

# Auto-switch configuration
auto_switch:
  enabled: true
  project_file: ".govman-version"

# Shell integration
shell:
  auto_detect: true
  completion: true

# Go releases API
go_releases:
  api_url: "https://go.dev/dl/?mode=json&include=all"
  download_url: "https://go.dev/dl/%s"
  cache_expiry: 10m

# Self-update configuration
self_update:
  github_api_url: "https://api.github.com/repos/sijunda/govman/releases/latest"
  github_releases_url: "https://api.github.com/repos/sijunda/govman/releases?per_page=1"

# Logging
quiet: false    # Suppress normal output
verbose: false  # Enable verbose logging
```

## Configuration Sections

### `install_dir` and `cache_dir`

-   `install_dir`: The directory where different Go versions will be installed.
-   `cache_dir`: The directory where downloaded Go archives are stored. `govman` will use these cached files to avoid re-downloading.

### `default_version`

-   The Go version to be used by default in new shell sessions.
-   This value is set automatically when you run `govman use <version> --default`.

### `download`

-   Customize the behavior of the download engine. You can disable parallel downloads or adjust connection and timeout settings if you are on an unstable network.

### `mirror`

-   `enabled`: Set to `true` to use the official Go mirror.
-   `url`: The base URL for the mirror. The default is the official mirror for users in China.

### `auto_switch`

-   `enabled`: Set to `false` to disable automatic version switching when changing directories.
-   `project_file`: The name of the file `govman` looks for to determine the project-specific version. Defaults to `.govman-version`.

### `logging`

-   `quiet`: Suppresses all output except for errors. Can be overridden by the `--quiet` flag.
-   `verbose`: Enables detailed logging for debugging. Can be overridden by the `--verbose` flag.