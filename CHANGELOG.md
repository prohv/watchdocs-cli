# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [v0.1.5] - 2026-06-23

### Changed
- Pruned binary targets to Windows (x64), Linux (x64), and Apple Silicon macOS (arm64) to optimize footprint (dropped Intel macOS)
- Compressed all binaries using UPX, reducing total package size from ~38 MB to ~7.2 MB

## [v0.1.4] - 2026-06-23

### Fixed
- README.md excluded from npm package — now shows correctly on npm package page

### Changed
- Unified `--help` output showing all commands and flags in a single view
- Updated package description and keywords for better discoverability

## [v0.1.2] - 2026-06-22

### Added
- `watchdocs lookup <package> --ecosystem <eco>` — single package lookup with no manifest needed
- `scan --path <dir>` flag — target a specific directory instead of cwd
- `scan --ecosystem <list>` flag — filter results to specific ecosystems
- `--slim` flag on both `scan` and `lookup` — returns only `name` and `docUrl` to save tokens
- Online resolvers for all ecosystems: npm, go, pip, cargo, pub, maven
- Parsers for `requirements.txt`, `pyproject.toml`, `uv.lock`, `Cargo.toml`, `pubspec.yaml`, `pom.xml`
- Recursive manifest discovery via `filepath.WalkDir` with skip list for heavy dirs (`node_modules`, `.git`, `vendor`, etc.)
- Concurrent resolution with 16-worker semaphore — all deps resolved in parallel
- Structured JSON output with `scanned`, `total`, `results`, and per-dep `status`
- Ecosystem-scoped deduplication — prevents duplicate entries when multiple pip manifests coexist
- Shared `internal/models` package with `Dependency` and `DocResult` types
- npm package scaffold (`package.json`, `bin/watchdocs.js`, `.npmignore`) for `npx watchdocs` usage
- `AGENTS.md` for AI agent discoverability

### Removed
- Gemini AI integration — replaced with direct online registry resolvers
- Table output format — replaced with structured JSON

### Changed
- `Dependency` struct extended with `Ecosystem` and `Type` fields
- Go parser now accepts string content instead of file path
- Scanner upgraded from flat root check to full recursive walk

## [v0.1.0-alpha] - 2026-04-03

### Added
- Initial CLI structure with Cobra (`watchdocs scan`)
- Manifest detection for `package.json` and `go.mod`
- Dependency extraction from NPM and Go modules
- Gemini AI integration for documentation URL resolution
- `GEMINI_API_KEY` validation before scan execution

### Fixed
- Go module parser skipping dependencies starting with `go` (e.g., `go.opencensus.io`)
- Typo in Gemini prompt (`librariesl` → `libraries`)
- Empty response handling from Gemini API
- Improved error messages for client creation and API failures
