# Internals & Advanced Troubleshooting

This guide is for developers who need to debug the internals of `govman` or diagnose complex issues.

## Enabling Verbose and Debug Logging

`govman` has a multi-level logger. To get more insight into its internal operations, use the `--verbose` flag.

```bash
govman --verbose <command>
```

This will print:
-   Step-by-step execution flow.
-   Timings for key operations (e.g., "version resolution", "download").
-   Internal progress messages.

For even more detail, the source code contains `_logger.Debug()` calls that are not enabled by default. To enable them, you would need to modify the logger level in the source code and recompile.

## Common Internal Failure Points

### 1. Configuration Loading (`internal/config`)

-   **Problem**: `govman` fails to start with a "failed to load config" error.
-   **Location**: `internal/config/config.go`
-   **Diagnosis**:
    -   Check if `~/.govman/config.yaml` is valid YAML.
    -   Ensure the user has read/write permissions for the `~/.govman/` directory.
    -   The `expandPath` function can fail if the home directory cannot be determined or if a configured path attempts path traversal (e.g., `~/../`).

### 2. Go Releases API Parsing (`internal/golang`)

-   **Problem**: `govman list --remote` fails or returns an empty list.
-   **Location**: `internal/golang/releases.go`
-   **Diagnosis**:
    -   The `fetchReleasesWithConfig` function might be failing due to a network issue or a change in the Go releases API (`https://go.dev/dl/?mode=json`).
    -   Check the `json.Unmarshal` step. If the structure of the API response has changed, the `Release` and `File` structs may need to be updated.
    -   The `resolveArch` function contains logic to handle architecture-specific download rules (e.g., for Apple Silicon before Go 1.16). An issue here could cause `govman` to look for a non-existent binary.

### 3. Download & Extraction (`internal/downloader`)

-   **Problem**: Installation fails during download or extraction.
-   **Location**: `internal/downloader/downloader.go`
-   **Diagnosis**:
    -   **Download**: The `downloadFile` function uses `http.NewRequest` with a `Range` header to support resuming downloads. A server that doesn't support this could cause issues. The retry logic is also here.
    -   **Checksum**: `verifyChecksum` reads the entire downloaded file to compute the SHA-256 hash. A mismatch here is a critical error.
    -   **Extraction**: The `extractTarGz` and `extractZip` functions contain security checks to prevent path traversal (`../`). An archive with an unsafe path will cause extraction to fail.

### 4. Shell Integration (`internal/shell`)

-   **Problem**: Auto-switching doesn't work, or the `govman use` command doesn't update the shell's `PATH`.
-   **Location**: `internal/shell/shell.go`
-   **Diagnosis**:
    -   The `Detect` function determines the user's current shell. If this detection is wrong, `govman init` will modify the wrong config file.
    -   The `SetupCommands` method for each shell type (`BashShell`, `ZshShell`, etc.) generates the script that gets injected into the user's profile. An error in this script can break shell integration.
    -   The `ExecutePathCommand` is what allows `govman use` to work for the current session. It prints a shell command to `stdout` that must be `eval`'d by the calling process. If this output is captured incorrectly, the `PATH` will not be updated.

## Simulating a Clean Install for Testing

When testing, it's often useful to simulate a clean environment. You can do this by temporarily setting the `HOME` environment variable.

```bash
# Create a temporary fake home directory
mkdir -p /tmp/govman-test-home

# Run govman with a different home directory
HOME=/tmp/govman-test-home ./build/govman install latest

# This will create all config and data in /tmp/govman-test-home/.govman/
```

This prevents your local development tests from interfering with your actual `govman` setup.