# WatchDocs CLI (Go)

WatchDocs is a lightweight CLI tool designed to bridge the gap between project dependencies and their official documentation. It automates the process of finding the right docs by scanning project manifests and using Gemini AI to resolve validated URLs.

## Getting Started

### 1. Prerequisites
- [Go](https://go.dev/doc/install) 1.21 or higher.
- A [Google Gemini API Key](https://aistudio.google.com/app/apikey).

### 2. Installation
Build the executable:
```powershell
go build -o watchdocs-cli.exe
```

### 3. Setup
Set your Gemini API key in your environment:
```powershell
# Windows (PowerShell)
$env:GEMINI_API_KEY = "your_api_key_here"

# Linux/macOS
export GEMINI_API_KEY="your_api_key_here"
```

### 4. Usage
Run the scan in your project directory:
```powershell
./watchdocs-cli.exe scan
```

---

### Command: `watchdocs scan`

* **Discovery**: Scans the project root for manifest files.
* **Parsing**: Extracts dependency names and versions.
* **Resolution**: Uses Gemini to identify official documentation URLs.
* **Validation**: Performs HEAD requests to ensure URLs are active.
* **Caching**: Stores results locally to prevent redundant API calls.
* **Output**: Renders a clean table or JSON output.

### Command: `watchdocs open <name>`

* **Browser Integration**: Opens the resolved documentation URL in the default system browser.

---

## Architecture

The tool is built on a modular flow:

1. **CLI Layer**: Handles user input and command execution.
2. **Scanner/Parser**: Detects files (e.g., `package.json`, `go.mod`) and extracts raw data.
3. **DocResolver**: Interface for converting library names into URLs.
4. **Cache Layer**: Checksums manifest files to skip the AI step for unchanged projects.
5. **Output/Browser**: Presents data and handles interaction.

---

## Technical Stack

* **Language**: Go 1.21+
* **CLI Framework**: Cobra
* **Formats Supported**: JSON, TOML, YAML, XML
* **AI Integration**: Google Gemini API (via `google-generative-ai-go`)
* **Data Persistence**: Local JSON cache in `~/.watchdocs/`

---

## Supported Manifests (v1)

* **Web**: `package.json`
* **Python**: `requirements.txt`, `pyproject.toml`
* **Go**: `go.mod`
* **Rust**: `Cargo.toml`
* **Java**: `pom.xml`
* **Flutter**: `pubspec.yaml`

---

## Gemini Integration Logic

* **Strict JSON**: Prompting ensures the AI returns only structured JSON.
* **Official Sources**: Constraints prioritize official documentation over blogs or third-party tutorials.
* **Version Awareness**: Attempts to resolve documentation specific to the installed library version.
* **Security**: Rejects non-HTTPS URLs and validates endpoints via HTTP HEAD checks.

---

## Core Modules

* `cmd/`: CLI definitions.
* `scanner/`: Manifest detection logic.
* `parsers/`: Manifest-specific extraction logic.
* `resolver/`: Gemini API integration and prompt management.
* `cache/`: SHA-256 manifest hashing and storage.
* `output/`: Terminal rendering and JSON formatting.
