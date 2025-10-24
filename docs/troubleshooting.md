# Troubleshooting

This guide provides solutions to common issues you might encounter while using `govman`.

## Installation Issues

### "Permission Denied" during installation

If the installation script fails with a "Permission Denied" error, it usually means it cannot write to your shell's configuration file (e.g., `~/.zshrc`) or create the `~/.govman` directory.

**Solution**:
1.  Ensure you have write permissions for your home directory.
2.  Try running the manual installation steps.

### `govman: command not found` after installation

This error means the `~/.govman/bin` directory was not successfully added to your shell's `PATH`.

**Solution**:
1.  **Restart your terminal**. This is the most common fix, as the `PATH` is only updated when a new shell session starts.
2.  Verify that `~/.govman/bin` is in your `PATH` by running `echo $PATH`.
3.  If it's missing, run `govman init --force` to re-attempt the shell configuration.
4.  As a last resort, manually add `export PATH="$HOME/.govman/bin:$PATH"` to your shell's profile file (`.bashrc`, `.zshrc`, etc.).

## Download & Installation Issues

### "Failed to get latest version information"

This error indicates that `govman` could not connect to the GitHub API to check for new versions.

**Solution**:
-   Check your internet connection.
-   If you are behind a corporate firewall, ensure that `api.github.com` is accessible.

### "Checksum mismatch"

This is a critical error indicating that the downloaded Go archive is corrupt and does not match the official SHA-256 checksum. `govman` will automatically delete the corrupt file.

**Solution**:
-   Run the `govman install` command again. A temporary network issue may have caused the corruption.
-   If the problem persists, run `govman clean` to clear the cache and try again.

## Version Switching Issues

### Auto-switching is not working

If `govman` doesn't automatically switch versions when you `cd` into a directory with a `.govman-version` file:

**Solution**:
1.  Run `govman init --force` to ensure your shell integration is correctly configured.
2.  Restart your terminal.
3.  Verify that `auto_switch: enabled: true` is set in your `~/.govman/config.yaml`.
4.  Make sure the `.govman-version` file contains a valid and installed Go version string.

### "Go version X is not installed"

This means you are trying to `use` a version that `govman` has not installed yet.

**Solution**:
-   Run `govman install <version>` to install the required version first.
-   Check your spelling or use `govman list --remote` to find the correct version string.

## Debugging

You can enable verbose logging to get more insight into what `govman` is doing behind the scenes.

```bash
govman --verbose <command>
```

For example, to debug a failing installation:
```bash
govman --verbose install 1.25.1
```

This will print detailed step-by-step logs, including API calls, file paths, and timer durations, which can be invaluable for diagnosing complex issues.