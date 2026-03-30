# WatchDocs CLI (Go)

WatchDocs is a lightweight CLI tool that scans project manifest files for dependencies and uses Google Gemini AI to automatically find their official documentation URLs.

## Getting Started

### 1. Prerequisites
- [Go](https://go.dev/doc/install) 1.25.2 or higher.
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

## Commands

### `watchdocs scan`

Scans the project root for manifest files, extracts dependencies, and resolves documentation URLs via Gemini AI.

**Workflow:**
1. **Discovery**: Scans the project root for supported manifest files.
2. **Parsing**: Extracts dependency names and versions.
3. **Resolution**: Uses Gemini AI to identify official documentation URLs.
4. **Output**: Displays results in a table format.

---

## Supported Manifests

Currently supported:
- **Web**: `package.json`
- **Go**: `go.mod`

---

## Technical Stack

- **Language**: Go 1.25.2+
- **CLI Framework**: Cobra
- **AI Integration**: Google Gemini API (`google.golang.org/genai`)
- **Model**: `gemini-3-flash-preview`

---

## Gemini Integration

- **Strict JSON Output**: Prompts ensure the AI returns only structured JSON.
- **Official Sources**: Prioritizes official documentation over third-party resources.
- **Version Awareness**: Attempts to resolve documentation for the specific library version.

---

## Project Structure

```
watchdocs-cli/
├── cmd/              # CLI command definitions
├── internal/
│   ├── scanner/      # Manifest file detection
│   ├── parser/       # Manifest-specific parsing logic
│   └── resolver/     # Gemini API integration
└── main.go           # Entry point
```
