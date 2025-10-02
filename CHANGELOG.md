# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2025-10-02

This release focuses on major improvements to cross-platform compatibility, build automation, code generation templates, and overall robustness.

### Added

- **GitHub Actions CI/CD Workflow**: A complete CI/CD pipeline (`.github/workflows/ci.yml`) has been added. It automatically runs tests, lints, and builds on every push and pull request. It also automates the creation of GitHub releases with cross-platform binaries when a new version tag is pushed.
- **Project `.gitignore`**: A comprehensive `.gitignore` file has been added to the project root to exclude build artifacts, test projects, and OS/IDE-specific files.

### Changed

- **Build-Time Version Injection**: The application version is no longer hardcoded. It defaults to `dev` and can be injected at build time using ldflags (`-X 'github.com/night-slayer18/goforge/cmd.version=...'`), which is a standard practice for Go applications.
- **`add` Command Reliability**: The `goforge add` command now automatically runs `go mod tidy` after getting a new dependency. This ensures the `go.sum` file is always consistent and prevents build failures.
- **Code Generation Templates**: The default templates for handlers, services, and the main application have been updated to ensure that a newly generated project or component compiles successfully out of the box without requiring manual code changes.

### Fixed

- **Cross-Platform `build` Command**: Replaced the Unix-specific `cp -r` command with a native Go implementation for copying assets, making the `build` command fully compatible with Windows, macOS, and Linux.
- **Cross-Platform `watch` Command**: The `watch` command is now fully cross-platform.
  - It now uses the correct shell (`cmd /C`) for executing scripts on Windows.
  - The automatic port cleanup feature now works on Windows (using `netstat`/`taskkill`) in addition to macOS/Linux (using `lsof`).
- **`watch` Command Configuration**: The `watch` command now correctly loads file-watching patterns from `goforge.yml` instead of using hardcoded values.
- **`add` Command Parsing**: The logic for parsing module versions in the `add` command has been improved to correctly handle module paths that may contain an `@` symbol.

### Removed

- **Obsolete Code**: Removed a redundant and unused `WatchMode` implementation from the `internal/runner` package, simplifying the codebase.
