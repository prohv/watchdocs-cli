# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

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
